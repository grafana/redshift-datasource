package redshift

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/driver"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v2"
	"github.com/pkg/errors"
)

type RedshiftDatasourceIface interface {
	sqlds.Driver
	sqlds.Completable
	sqlAPI.Resources
	Schemas(ctx context.Context, options sqlds.Options) ([]string, error)
	Tables(ctx context.Context, options sqlds.Options) ([]string, error)
	Columns(ctx context.Context, options sqlds.Options) ([]string, error)
	Secrets(ctx context.Context, options sqlds.Options) ([]models.ManagedSecret, error)
	Secret(ctx context.Context, options sqlds.Options) (*models.RedshiftSecret, error)
}

type RedshiftDatasource struct {
	sessionCache *awsds.SessionCache
}

func New() *RedshiftDatasource {
	return &RedshiftDatasource{
		sessionCache: awsds.NewSessionCache(),
	}
}

func (s *RedshiftDatasource) Settings(_ backend.DataSourceInstanceSettings) sqlds.DriverSettings {
	return sqlds.DriverSettings{
		FillMode: &data.FillMissing{
			Mode: data.FillModeNull,
		},
	}
}

func (s *RedshiftDatasource) Converters() (sc []sqlutil.Converter) {
	return sc
}

// Connect opens a sql.DB connection using datasource settings
func (s *RedshiftDatasource) Connect(config backend.DataSourceInstanceSettings, _ json.RawMessage) (*sql.DB, error) {
	settings := models.RedshiftDataSourceSettings{}
	err := settings.Load(config)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	api, err := api.New(s.sessionCache, &settings)
	if err != nil {
		return nil, err
	}

	db, err := driver.Open(api)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to connect to database. Is the hostname and port correct?")
	}

	return db, nil
}

func (s *RedshiftDatasource) getApi(ctx context.Context) (*api.API, error) {
	plugin := httpadapter.PluginConfigFromContext(ctx)
	if plugin.DataSourceInstanceSettings == nil {
		return nil, fmt.Errorf("unable to get settings from request")
	}

	settings := models.RedshiftDataSourceSettings{}
	err := settings.Load(*plugin.DataSourceInstanceSettings)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	res, err := api.New(s.sessionCache, &settings)
	if err != nil {
		return nil, err
	}
	return res.(*api.API), nil
}

func (s *RedshiftDatasource) Regions(ctx context.Context) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	regions, err := api.Regions(ctx)
	if err != nil {
		return nil, err
	}
	return regions, nil
}

func (s *RedshiftDatasource) Databases(ctx context.Context, options sqlds.Options) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	dbs, err := api.Databases(ctx, options)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (s *RedshiftDatasource) Schemas(ctx context.Context, options sqlds.Options) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	schemas, err := api.Schemas(ctx, options)
	if err != nil {
		return nil, err
	}
	return schemas, nil
}

func (s *RedshiftDatasource) Tables(ctx context.Context, options sqlds.Options) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	tables, err := api.Tables(ctx, options)
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func (s *RedshiftDatasource) Columns(ctx context.Context, options sqlds.Options) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	cols, err := api.Columns(ctx, options)
	if err != nil {
		return nil, err
	}
	return cols, nil
}

func (s *RedshiftDatasource) Secrets(ctx context.Context, options sqlds.Options) ([]models.ManagedSecret, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	return api.Secrets(ctx)
}

func (s *RedshiftDatasource) Secret(ctx context.Context, options sqlds.Options) (*models.RedshiftSecret, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	return api.Secret(ctx, options)
}
