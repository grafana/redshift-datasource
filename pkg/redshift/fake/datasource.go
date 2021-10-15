package fake

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v2"
)

type RedshiftFakeDatasource struct {
	SecretList []models.ManagedSecret
	RSecret    models.RedshiftSecret
}

func (s *RedshiftFakeDatasource) Settings(_ backend.DataSourceInstanceSettings) sqlds.DriverSettings {
	return sqlds.DriverSettings{}
}

func (s *RedshiftFakeDatasource) Converters() (sc []sqlutil.Converter) {
	return sc
}

func (s *RedshiftFakeDatasource) Connect(config backend.DataSourceInstanceSettings, queryArgs json.RawMessage) (*sql.DB, error) {
	return &sql.DB{}, nil
}

func (s *RedshiftFakeDatasource) Macros() sqlds.Macros {
	return sqlds.Macros{}
}

func (s *RedshiftFakeDatasource) Schemas(ctx context.Context, options sqlds.Options) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftFakeDatasource) Tables(ctx context.Context, options sqlds.Options) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftFakeDatasource) Columns(ctx context.Context, options sqlds.Options) ([]string, error) {
	return []string{}, nil
}

func (s *RedshiftFakeDatasource) Secrets(ctx context.Context) ([]models.ManagedSecret, error) {
	return s.SecretList, nil
}

func (s *RedshiftFakeDatasource) Secret(ctx context.Context, arn string) (*models.RedshiftSecret, error) {
	return &s.RSecret, nil
}
