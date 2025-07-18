package redshift

import (
	"context"
	"database/sql"
	"encoding/json"

	awsDriver "github.com/grafana/grafana-aws-sdk/pkg/sql/driver"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/driver/async"
	sqlModels "github.com/grafana/grafana-aws-sdk/pkg/sql/models"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/driver"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v4"
)

type RedshiftDatasourceIface interface {
	sqlds.Driver
	sqlds.Completable
	sqlAPI.Resources
	awsds.AsyncDriver
	Schemas(ctx context.Context, options sqlds.Options) ([]string, error)
	Tables(ctx context.Context, options sqlds.Options) ([]string, error)
	Columns(ctx context.Context, options sqlds.Options) ([]string, error)
	Secrets(ctx context.Context, options sqlds.Options) ([]models.ManagedSecret, error)
	Secret(ctx context.Context, options sqlds.Options) (*models.RedshiftSecret, error)
	Clusters(ctx context.Context, options sqlds.Options) ([]models.RedshiftCluster, error)
	Workgroups(ctx context.Context, options sqlds.Options) ([]models.RedshiftWorkgroup, error)
}

type Loader struct{}

func (l Loader) LoadAPI(ctx context.Context, settings sqlModels.Settings) (sqlAPI.AWSAPI, error) {
	return api.New(ctx, settings)
}

func (l Loader) LoadDriver(ctx context.Context, awsapi sqlAPI.AWSAPI) (awsDriver.Driver, error) {
	return driver.New(ctx, awsapi)
}

func (l Loader) LoadAsyncDriver(ctx context.Context, awsapi sqlAPI.AWSAPI) (async.Driver, error) {
	return driver.New(ctx, awsapi)
}

func (l Loader) LoadSettings(ctx context.Context) sqlModels.Settings {
	return models.New(ctx)
}

type RedshiftDatasource struct {
	awsDS datasource.AWSClient
}

func New() *RedshiftDatasource {
	return &RedshiftDatasource{awsDS: datasource.New(Loader{})}
}

func (s *RedshiftDatasource) Settings(ctx context.Context, _ backend.DataSourceInstanceSettings) sqlds.DriverSettings {
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
func (s *RedshiftDatasource) Connect(ctx context.Context, config backend.DataSourceInstanceSettings, queryArgs json.RawMessage) (*sql.DB, error) {
	s.awsDS.Init(config)
	args, err := sqlds.ParseOptions(queryArgs)
	if err != nil {
		return nil, err
	}
	args["updated"] = config.Updated.String()

	return s.awsDS.GetDB(ctx, config.ID, args)
}

func (s *RedshiftDatasource) GetAsyncDB(ctx context.Context, config backend.DataSourceInstanceSettings, queryArgs json.RawMessage) (awsds.AsyncDB, error) {
	s.awsDS.Init(config)
	args, err := sqlds.ParseOptions(queryArgs)
	if err != nil {
		return nil, err
	}
	args["updated"] = config.Updated.String()

	return s.awsDS.GetAsyncDB(ctx, config.ID, args)
}

func (s *RedshiftDatasource) getApi(ctx context.Context, options sqlds.Options) (*api.API, error) {
	id := datasource.GetDatasourceID(ctx)
	args := sqlds.Options{}
	for key, val := range options {
		args[key] = val
	}
	// the updated time makes sure that we don't use a token for a stale version of the datasource
	args["updated"] = datasource.GetDatasourceLastUpdatedTime(ctx)

	res, err := s.awsDS.GetAPI(ctx, id, args)
	if err != nil {
		return nil, err
	}

	return res.(*api.API), err
}

func (s *RedshiftDatasource) Regions(ctx context.Context) ([]string, error) {
	// This is not used. If regions are out of date, update them in the @grafana/aws-sdk-react package
	return []string{}, nil
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

func (s *RedshiftDatasource) CancelQuery(ctx context.Context, options sqlds.Options, queryID string) error {
	api, err := s.getApi(ctx, options)
	if err != nil {
		return err
	}
	return api.CancelQuery(ctx, options, queryID)
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

func (s *RedshiftDatasource) Clusters(ctx context.Context, options sqlds.Options) ([]models.RedshiftCluster, error) {
	api, err := s.getApi(ctx, options)
	if err != nil {
		return nil, err
	}
	return api.Clusters(ctx)
}

func (s *RedshiftDatasource) Workgroups(ctx context.Context, options sqlds.Options) ([]models.RedshiftWorkgroup, error) {
	api, err := s.getApi(ctx, options)
	if err != nil {
		return nil, err
	}
	return api.Workgroups(ctx)
}
