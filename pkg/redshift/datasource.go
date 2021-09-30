package redshift

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"time"

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

var (
	// TODO: This supports basic table names (which is incomplete)
	tableNameRegex = regexp.MustCompile("^[0-9A-Za-z_]+$")
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
		Timeout: 5 * time.Minute,
	}
}

// Connect opens a sql.DB connection using datasource settings
func (s *RedshiftDatasource) Connect(config backend.DataSourceInstanceSettings, _ json.RawMessage) (*sql.DB, error) {
	settings := models.RedshiftDataSourceSettings{}
	err := settings.Load(config)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	db, err := driver.Open(settings, s.sessionCache)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to connect to database. Is the hostname and port correct?")
	}
	s.db = db

	return db, nil
}

func (s *RedshiftDatasource) Converters() (sc []sqlutil.Converter) {
	return []sqlutil.Converter{{ // This converter can be removed as soon as it's a part of SQLUtil. See https://github.com/grafana/grafana-plugin-sdk-go/pull/369
		Name:          "nullable bool converter",
		InputScanType: reflect.TypeOf(sql.NullBool{}),
		InputTypeName: "BOOLEAN",
		FrameConverter: sqlutil.FrameConverter{
			FieldType: data.FieldTypeNullableBool,
			ConverterFunc: func(n interface{}) (interface{}, error) {
				v := n.(*sql.NullBool)

				if !v.Valid {
					return (*bool)(nil), nil
				}

				return &v.Bool, nil
			},
		},
	}}
}

func getStringArr(rows *sql.Rows) ([]string, error) {
	res := []string{}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		res = append(res, name)
	}
	return res, nil
}

func (s *RedshiftDatasource) Schemas(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT nspname FROM pg_namespace")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getStringArr(rows)
}

func (s *RedshiftDatasource) Tables(ctx context.Context, schema string) ([]string, error) {
	// We use the "public" schema by default if not specified
	if schema == "" {
		schema = "public"
	}
	if !tableNameRegex.Match([]byte(schema)) {
		return nil, fmt.Errorf("unsupported schema name %s", schema)
	}
	// Rather than injecting the table_schema, we should use arguments but the Redshift API only allow
	// arguments for prepared statements (which has no support in the golang sdk)
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema='%s'", schema))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getStringArr(rows)
}

func (s *RedshiftDatasource) Columns(ctx context.Context, table string) ([]string, error) {
	if !tableNameRegex.Match([]byte(table)) {
		return nil, fmt.Errorf("unsupported table name %s", table)
	}
	// Rather than injecting the table_name, we should use arguments but the Redshift API only allow
	// arguments for prepared statements (which has no support in the golang sdk)
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%s'", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getStringArr(rows)
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
