package redshift

import (
	"context"
	"database/sql"
	"encoding/json"

	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/driver"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v2"
)

type RedshiftDatasourceIface interface {
	sqlds.Driver
	sqlds.Completable
	sqlAPI.Resources
	sqlds.AsyncDBGetter
	Schemas(ctx context.Context, options sqlds.Options) ([]string, error)
	Tables(ctx context.Context, options sqlds.Options) ([]string, error)
	Columns(ctx context.Context, options sqlds.Options) ([]string, error)
	Secrets(ctx context.Context, options sqlds.Options) ([]models.ManagedSecret, error)
	Secret(ctx context.Context, options sqlds.Options) (*models.RedshiftSecret, error)
	Cluster(ctx context.Context, options sqlds.Options) (*models.RedshiftCluster, error)
}

type RedshiftDatasource struct {
	awsDS *datasource.AWSDatasource
}

func New() *RedshiftDatasource {
	return &RedshiftDatasource{awsDS: datasource.New()}
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
func (s *RedshiftDatasource) Connect(config backend.DataSourceInstanceSettings, queryArgs json.RawMessage) (*sql.DB, error) {
	s.awsDS.Init(config)
	args, err := sqlds.ParseOptions(queryArgs)
	if err != nil {
		return nil, err
	}

	return s.awsDS.GetDB(config.ID, args, models.New, api.New, driver.NewSync)
}

func (s *RedshiftDatasource) GetAsyncDB(config backend.DataSourceInstanceSettings, queryArgs json.RawMessage) (sqlds.AsyncDB, error) {
	s.awsDS.Init(config)
	args, err := sqlds.ParseOptions(queryArgs)
	if err != nil {
		return nil, err
	}

	return s.awsDS.GetAsyncDB(config.ID, args, models.New, api.New, driver.New)
}

func (s *RedshiftDatasource) getApi(ctx context.Context, options sqlds.Options) (*api.API, error) {
	id := datasource.GetDatasourceID(ctx)
	res, err := s.awsDS.GetAPI(id, options, models.New, api.New)
	return res.(*api.API), err
}

func (s *RedshiftDatasource) Regions(ctx context.Context) ([]string, error) {
	api, err := s.getApi(ctx, sqlds.Options{})
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
	api, err := s.getApi(ctx, options)
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
	api, err := s.getApi(ctx, options)
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
	api, err := s.getApi(ctx, options)
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
	api, err := s.getApi(ctx, options)
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
	api, err := s.getApi(ctx, options)
	if err != nil {
		return nil, err
	}
	return api.Secrets(ctx)
}

func (s *RedshiftDatasource) Secret(ctx context.Context, options sqlds.Options) (*models.RedshiftSecret, error) {
	api, err := s.getApi(ctx, options)
	if err != nil {
		return nil, err
	}
	return api.Secret(ctx, options)
}

func (s *RedshiftDatasource) Cluster(ctx context.Context, options sqlds.Options) (*models.RedshiftCluster, error) {
	api, err := s.getApi(ctx, options)
	if err != nil {
		return nil, err
	}
	return api.Cluster(options)
}
