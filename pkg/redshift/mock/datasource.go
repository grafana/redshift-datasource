package mock

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v2"
)

type RedshiftMockDatasource struct {
	SecretList  []models.ManagedSecret
	RSecret     models.RedshiftSecret
	RClusters   []models.RedshiftCluster
	RWorkgroups []models.RedshiftWorkgroup
}

func (s *RedshiftMockDatasource) Settings(_ backend.DataSourceInstanceSettings) sqlds.DriverSettings {
	return sqlds.DriverSettings{}
}

func (s *RedshiftMockDatasource) Converters() (sc []sqlutil.Converter) {
	return sc
}

func (s *RedshiftMockDatasource) Connect(config backend.DataSourceInstanceSettings, queryArgs json.RawMessage) (*sql.DB, error) {
	return &sql.DB{}, nil
}

func (s *RedshiftMockDatasource) GetAsyncDB(config backend.DataSourceInstanceSettings, queryArgs json.RawMessage) (awsds.AsyncDB, error) {
	return nil, nil
}

func (s *RedshiftMockDatasource) Macros() sqlds.Macros {
	return sqlds.Macros{}
}

func (s *RedshiftMockDatasource) Regions(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftMockDatasource) Databases(ctx context.Context, options sqlds.Options) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftMockDatasource) CancelQuery(ctx context.Context, options sqlds.Options, queryID string) error {
	return nil
}

func (s *RedshiftMockDatasource) Schemas(ctx context.Context, options sqlds.Options) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftMockDatasource) Tables(ctx context.Context, options sqlds.Options) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftMockDatasource) Columns(ctx context.Context, options sqlds.Options) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftMockDatasource) Secrets(ctx context.Context, options sqlds.Options) ([]models.ManagedSecret, error) {
	return s.SecretList, nil
}

func (s *RedshiftMockDatasource) Secret(ctx context.Context, options sqlds.Options) (*models.RedshiftSecret, error) {
	return &s.RSecret, nil
}
func (s *RedshiftMockDatasource) Clusters(ctx context.Context, options sqlds.Options) ([]models.RedshiftCluster, error) {
	return s.RClusters, nil
}
func (s *RedshiftMockDatasource) Workgroups(ctx context.Context, options sqlds.Options) ([]models.RedshiftWorkgroup, error) {
	return s.RWorkgroups, nil
}
