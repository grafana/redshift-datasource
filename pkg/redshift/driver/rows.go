package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	redshiftdatatypes "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Rows struct {
	service redshiftdata.GetStatementResultAPIClient
	queryID string
	context context.Context

	done   bool
	result *redshiftdata.GetStatementResultOutput
}

func newRows(ctx context.Context, service redshiftdata.GetStatementResultAPIClient, queryId string) (*Rows, error) {
	r := Rows{
		service: service,
		queryID: queryId,
		context: ctx,
	}

	if err := r.fetchNextPage(nil); err != nil {
		return nil, err
	}

	return &r, nil
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide. io.EOF should be returned when there are no more rows.
func (r *Rows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}

	// If nothing left to iterate...
	if len(r.result.Records) == 0 {
		// And if nothing more to paginate...
		if r.result.NextToken == nil || *r.result.NextToken == "" {
			r.done = true
			return io.EOF
		}

		err := r.fetchNextPage(r.result.NextToken)
		if err != nil {
			return err
		}
	}

	// Shift to next row
	current := r.result.Records[0]
	if err := convertRow(r.result.ColumnMetadata, current, dest); err != nil {
		return err
	}

	r.result.Records = r.result.Records[1:]
	return nil
}

// Columns returns the names of the columns.
func (r *Rows) Columns() []string {
	columnNames := []string{}
	for _, column := range r.result.ColumnMetadata {
		columnNames = append(columnNames, *column.Name)
	}
	return columnNames
}

// ColumnTypeNullable returns true if it is known the column may be null,
// or false if the column is known to be not nullable. If the column nullability is unknown, ok should be false.
func (r *Rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	return r.result.ColumnMetadata[index].Nullable == 1, true
}

// ColumnTypeScanType returns the value type that can be used to scan types into.
// For example, the database column type "bigint" this should return "reflect.TypeOf(int64(0))"
func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	col := r.result.ColumnMetadata[index]

	switch strings.ToUpper(*col.TypeName) {
	case REDSHIFT_INT, REDSHIFT_INT4,
		REDSHIFT_NUMERIC, REDSHIFT_FLOAT, REDSHIFT_FLOAT8:
		// If the value is numeric and the name is "time", assume a Unix timestamp
		if *col.Name == "time" {
			return reflect.TypeOf(time.Time{})
		}
	}

	switch strings.ToUpper(*col.TypeName) {
	case REDSHIFT_INT2:
		return reflect.TypeOf(int16(0))
	case REDSHIFT_INT, REDSHIFT_INT4:
		return reflect.TypeOf(int32(0))
	case REDSHIFT_INT8:
		return reflect.TypeOf(int64(0))
	case REDSHIFT_FLOAT4:
		return reflect.TypeOf(float32(0))
	case REDSHIFT_NUMERIC, REDSHIFT_FLOAT, REDSHIFT_FLOAT8:
		return reflect.TypeOf(float64(0))
	case REDSHIFT_BOOL:
		return reflect.TypeOf(false)
	case REDSHIFT_CHARACTER,
		REDSHIFT_VARCHAR,
		REDSHIFT_NCHAR,
		REDSHIFT_BPCHAR,
		REDSHIFT_CHARACTER_VARYING,
		REDSHIFT_NVARCHAR,
		REDSHIFT_TEXT:
		return reflect.TypeOf("")
	case REDSHIFT_TIMESTAMP,
		REDSHIFT_TIMESTAMP_WITH_TIME_ZONE,
		REDSHIFT_TIME_WITHOUT_TIME_ZONE,
		REDSHIFT_TIME_WITH_TIME_ZONE:
		return reflect.TypeOf(time.Time{})
	default:
		return reflect.TypeOf("")
	}
}

// ColumnTypeDatabaseTypeName converts a redshift data type to a corresponding go sql type
func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	columnTypeMapper := map[string]string{
		REDSHIFT_INT2:                     "SMALLINT",
		REDSHIFT_INT:                      "INTEGER",
		REDSHIFT_INT4:                     "INTEGER",
		REDSHIFT_INT8:                     "BIGINT",
		REDSHIFT_NUMERIC:                  "DECIMAL",
		REDSHIFT_FLOAT4:                   "REAL",
		REDSHIFT_FLOAT8:                   "DOUBLE",
		REDSHIFT_FLOAT:                    "DOUBLE",
		REDSHIFT_BOOL:                     "BOOLEAN",
		REDSHIFT_CHARACTER:                "CHAR",
		REDSHIFT_NCHAR:                    "CHAR",
		REDSHIFT_BPCHAR:                   "CHAR",
		REDSHIFT_CHARACTER_VARYING:        "VARCHAR",
		REDSHIFT_NVARCHAR:                 "VARCHAR",
		REDSHIFT_TEXT:                     "VARCHAR",
		REDSHIFT_VARCHAR:                  "VARCHAR",
		REDSHIFT_DATE:                     "DATE",
		REDSHIFT_TIMESTAMP:                "TIMESTAMP",
		REDSHIFT_TIMESTAMP_WITH_TIME_ZONE: "TIMESTAMPTZ",
		REDSHIFT_TIME_WITHOUT_TIME_ZONE:   "TIME",
		REDSHIFT_TIME_WITH_TIME_ZONE:      "TIMETZ",
		REDSHIFT_GEOMETRY:                 "GEOMETRY",
		// HLLSKETCH and SUPER are redshift specific types
		REDSHIFT_HLLSKETCH: "VARCHAR",
		REDSHIFT_SUPER:     "VARCHAR",
	}

	typeName := strings.ToUpper(*r.result.ColumnMetadata[index].TypeName)
	if val, ok := columnTypeMapper[typeName]; ok {
		return val
	}

	backend.Logger.Warn("unexpected type, using VARCHAR instead", "type name", typeName)
	return "VARCHAR"
}

// Close closes the rows iterator.
func (r *Rows) Close() error {
	r.done = true
	return nil
}

// fetchNextPage fetches the next statement result page and adds the result to the row
func (r *Rows) fetchNextPage(token *string) error {
	var err error

	r.result, err = r.service.GetStatementResult(r.context, &redshiftdata.GetStatementResultInput{
		Id:        aws.String(r.queryID),
		NextToken: token,
	})

	if err != nil {
		return err
	}

	return nil
}

// convertRow converts values in a redshift data api row into its corresponding type in Go. Mapping is based on:
// https://docs.aws.amazon.com/redshift/latest/dg/c_Supported_data_types.html
// https://docs.aws.amazon.com/redshift/latest/mgmt/jdbc20-data-type-mapping.html
func convertRow(columns []redshiftdatatypes.ColumnMetadata, data []redshiftdatatypes.Field, ret []driver.Value) error {
	for i, curr := range data {
		// FIXME: I think this is the correct translation of the previous aws-sdk-v1 behavior
		// but I'm not sure it's actually the correct behavior
		if isNull, ok := curr.(*redshiftdatatypes.FieldMemberIsNull); ok && isNull.Value {
			if isNull.Value {
				ret[i] = nil
			}
			continue
		}

		col := columns[i]
		if col.TypeName == nil {
			return fmt.Errorf("error in convertRow: col.TypeName is nil")
		}
		typeName := strings.ToUpper(*col.TypeName)
		switch typeName {
		case REDSHIFT_INT2:
			if long, ok := AsInt(curr); ok {
				ret[i] = int16(long)
			} else {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
		case REDSHIFT_INT, REDSHIFT_INT4:
			if long, ok := AsInt(curr); ok {
				if *col.Name == "time" {
					ret[i] = time.Unix(long, 0).UTC()
				} else {
					ret[i] = int32(long)
				}
			} else {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
		case REDSHIFT_INT8:
			if long, ok := AsInt(curr); ok {
				ret[i] = long
			} else {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
		case REDSHIFT_NUMERIC:
			s, ok := AsString(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			value, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			ret[i] = value
		case REDSHIFT_FLOAT, REDSHIFT_FLOAT4, REDSHIFT_FLOAT8:
			value, ok := AsFloat(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			if *col.Name == "time" {
				ret[i] = time.Unix(int64(value), 0).UTC()
			} else {
				ret[i] = value
			}
		case REDSHIFT_BOOL:
			value, ok := AsBool(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			ret[i] = value

		case REDSHIFT_CHARACTER,
			REDSHIFT_VARCHAR,
			REDSHIFT_NCHAR,
			REDSHIFT_BPCHAR,
			REDSHIFT_CHARACTER_VARYING,
			REDSHIFT_NVARCHAR,
			REDSHIFT_TEXT,
			// Complex types are returned as a string
			REDSHIFT_GEOMETRY,
			REDSHIFT_HLLSKETCH,
			REDSHIFT_SUPER,
			REDSHIFT_NAME:
			value, ok := AsString(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			ret[i] = value
		// Time formats from
		// https://docs.aws.amazon.com/redshift/latest/dg/r_Datetime_types.html
		case REDSHIFT_DATE:
			value, ok := AsString(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			t, err := time.Parse("2006-01-02", value)
			if err != nil {
				return err
			}
			ret[i] = t
		case REDSHIFT_TIMESTAMP:
			value, ok := AsString(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			t, err := time.Parse("2006-01-02 15:04:05", value)
			if err != nil {
				return err
			}
			ret[i] = t
		case REDSHIFT_TIMESTAMP_WITH_TIME_ZONE:
			value, ok := AsString(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			t, err := time.Parse("2006-01-02 15:04:05+00", value)
			if err != nil {
				return err
			}
			ret[i] = t
		case REDSHIFT_TIME_WITHOUT_TIME_ZONE,
			REDSHIFT_TIME_WITH_TIME_ZONE:
			value, ok := AsString(curr)
			if !ok {
				return fmt.Errorf("column %s with typeName %s could not be converted", *col.Name, *col.TypeName)
			}
			t, err := time.Parse("15:04:05", value)
			if err != nil {
				return err
			}
			ret[i] = t
		default:
			return fmt.Errorf("unsupported type %s", typeName)
		}
	}
	return nil
}

func AsInt(field redshiftdatatypes.Field) (int64, bool) {
	var value int64
	v, ok := field.(*redshiftdatatypes.FieldMemberLongValue)
	if ok {
		value = v.Value
	}
	return value, ok
}

func AsFloat(field redshiftdatatypes.Field) (float64, bool) {
	var value float64
	v, ok := field.(*redshiftdatatypes.FieldMemberDoubleValue)
	if ok {
		value = v.Value
	}
	return value, ok
}

func AsString(field redshiftdatatypes.Field) (string, bool) {
	var value string
	v, ok := field.(*redshiftdatatypes.FieldMemberStringValue)
	if ok {
		value = v.Value
	}
	return value, ok

}
func AsBool(field redshiftdatatypes.Field) (bool, bool) {
	var value bool
	v, ok := field.(*redshiftdatatypes.FieldMemberBooleanValue)
	if ok {
		value = v.Value
	}
	return value, ok

}
