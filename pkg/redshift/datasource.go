package redshift

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/redshift-datasource/pkg/redshift/driver"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds"
	"github.com/pkg/errors"
)

type RedshiftDatasource struct{}

func (s *RedshiftDatasource) FillMode() *data.FillMissing {
	return &data.FillMissing{
		Mode: data.FillModeNull,
	}
}

// Connect opens a sql.DB connection using datasource settings
func (s *RedshiftDatasource) Connect(config backend.DataSourceInstanceSettings) (*sql.DB, error) {
	settings := models.RedshiftDataSourceSettings{}
	err := settings.Load(config)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	db, err := driver.Open(settings)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to connect to database. Is the hostname and port correct?")
	}

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

func (s *RedshiftDatasource) Macros() sqlds.Macros {
	return nil
}
