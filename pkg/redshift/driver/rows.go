package driver

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice/redshiftdataapiserviceiface"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Rows struct {
	service redshiftdataapiserviceiface.RedshiftDataAPIServiceAPI
	queryID string

	done   bool
	result *redshiftdataapiservice.GetStatementResultOutput
}

func newRows(service redshiftdataapiserviceiface.RedshiftDataAPIServiceAPI, queryId string) (*Rows, error) {
	r := Rows{
		service: service,
		queryID: queryId,
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
	col := *r.result.ColumnMetadata[index]

	if *col.Nullable == 1 {
		return true, true
	}

	return false, true
}

// ColumnTypeScanType returns the value type that can be used to scan types into.
// For example, the database column type "bigint" this should return "reflect.TypeOf(int64(0))"
func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	col := *r.result.ColumnMetadata[index]

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
		if *col.Name == "time" {
			return reflect.TypeOf(time.Time{})
		}
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

	r.result, err = r.service.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
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
func convertRow(columns []*redshiftdataapiservice.ColumnMetadata, data []*redshiftdataapiservice.Field, ret []driver.Value) error {
	for i, curr := range data {
		if curr.IsNull != nil && *curr.IsNull {
			continue
		}

		col := columns[i]
		typeName := strings.ToUpper(*col.TypeName)
		switch typeName {
		case REDSHIFT_INT2:
			ret[i] = int16(*curr.LongValue)
		case REDSHIFT_INT, REDSHIFT_INT4:
			ret[i] = int32(*curr.LongValue)
		case REDSHIFT_INT8:
			ret[i] = *curr.LongValue
		case REDSHIFT_NUMERIC, REDSHIFT_FLOAT, REDSHIFT_FLOAT4:
			bitSize := 64
			if typeName == REDSHIFT_FLOAT4 {
				bitSize = 32
			}
			v, err := strconv.ParseFloat(*curr.StringValue, bitSize)
			if err != nil {
				return err
			}
			ret[i] = v
		case REDSHIFT_FLOAT8:
			if *col.Name == "time" {
				ret[i] = time.Unix(int64(*curr.DoubleValue), 0)
			} else {
				ret[i] = *curr.DoubleValue
			}
		case REDSHIFT_BOOL:
			// don't know why boolean values are not passed as curr.BooleanValue
			boolValue, err := strconv.ParseBool(*curr.StringValue)
			if err != nil {
				return err
			}
			ret[i] = boolValue

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
			ret[i] = *curr.StringValue
		// Time formats from
		// https://docs.aws.amazon.com/redshift/latest/dg/r_Datetime_types.html
		case REDSHIFT_DATE:
			t, err := time.Parse("2006-01-02", *curr.StringValue)
			if err != nil {
				return err
			}
			ret[i] = t
		case REDSHIFT_TIMESTAMP:
			t, err := time.Parse("2006-01-02 15:04:05", *curr.StringValue)
			if err != nil {
				return err
			}
			ret[i] = t
		case REDSHIFT_TIMESTAMP_WITH_TIME_ZONE:
			t, err := time.Parse("2006-01-02 15:04:05+00", *curr.StringValue)
			if err != nil {
				return err
			}
			ret[i] = t
		case REDSHIFT_TIME_WITHOUT_TIME_ZONE,
			REDSHIFT_TIME_WITH_TIME_ZONE:
			t, err := time.Parse("15:04:05", *curr.StringValue)
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
