package redshift

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
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
	Schemas(ctx context.Context) ([]string, error)
	Tables(ctx context.Context, schema string) ([]string, error)
	Columns(ctx context.Context, table string) ([]string, error)
	Secrets(ctx context.Context) ([]models.ManagedSecret, error)
	Secret(ctx context.Context, arn string) (*models.RedshiftSecret, error)
}

type RedshiftDatasource struct {
	sessionCache *awsds.SessionCache
	db           *sql.DB
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
	s.db = db

	return db, nil
}

func (s *RedshiftDatasource) Converters() (sc []sqlutil.Converter) {
	return sc
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

	return api.New(s.sessionCache, &settings)
}

func (s *RedshiftDatasource) Schemas(ctx context.Context) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	schemas, err := api.ListSchemas(ctx)
	if err != nil {
		return nil, err
	}
	return schemas, nil
}

func (s *RedshiftDatasource) Tables(ctx context.Context, schema string) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	tables, err := api.ListTables(ctx, schema)
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func (s *RedshiftDatasource) Columns(ctx context.Context, table string) ([]string, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: Add support for other schemas
	cols, err := api.ListColumns(ctx, "public", table)
	if err != nil {
		return nil, err
	}
	return cols, nil
}

func (s *RedshiftDatasource) Secrets(ctx context.Context) ([]models.ManagedSecret, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	return api.ListSecrets(ctx)
}

func (s *RedshiftDatasource) Secret(ctx context.Context, arn string) (*models.RedshiftSecret, error) {
	api, err := s.getApi(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetSecret(ctx, arn)
}
