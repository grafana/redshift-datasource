package driver

import (
	"database/sql/driver"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

type Rows struct {
	client  *redshiftdataapiservice.RedshiftDataAPIService
	queryID string

	done      bool
	result    *redshiftdataapiservice.GetStatementResultOutput
}

func newRows(client *redshiftdataapiservice.RedshiftDataAPIService, queryId string) (*Rows, error) {
	r := Rows{
		client:  client,
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
	case "INT2":
		return reflect.TypeOf(int16(0))
	case "INT", "INT4":
		return reflect.TypeOf(int32(0))
	case "INT8":
		return reflect.TypeOf(int64(0))
	case "FLOAT4":
		return reflect.TypeOf(float32(0))
	case "NUMERIC", "FLOAT", "FLOAT8":
		return reflect.TypeOf(float64(0))
	case "BOOL":
		return reflect.TypeOf(false)
	case "CHARACTER", "NCHAR", "BPCHAR", "VARYING", "NVARCHAR", "TEXT":
		return reflect.TypeOf("")
	case "TIMESTAMP":
		return reflect.TypeOf(time.Time{})
	default:
		return reflect.TypeOf("")
	}
}

// ColumnTypeDatabaseTypeName converts a redshift data type to a corresponding go sql type
func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
	columnTypeMapper := map[string]string{
		"INT2":                        "SMALLINT",
		"INT":                         "INTEGER",
		"INT4":                        "INTEGER",
		"INT8":                        "BIGINT",
		"NUMERIC":                     "DECIMAL",
		"FLOAT4":                      "REAL",
		"FLOAT8":                      "DOUBLE",
		"FLOAT":                       "DOUBLE",
		"BOOL":                        "BOOLEAN",
		"CHARACTER":                   "CHAR",
		"NCHAR":                       "CHAR",
		"BPCHAR":                      "CHAR",
		"CHARACTER VARYING":           "VARCHAR",
		"NVARCHAR":                    "VARCHAR",
		"TEXT":                        "VARCHAR",
		"DATE":                        "DATE",
		"TIMESTAMP":                   "TIMESTAMP",
		"TIMESTAMP WITHOUT TIME ZONE": "TIMESTAMP",
		"TIMESTAMP WITH TIME ZONE":    "TIMESTAMPTZ",
		"TIME WITHOUT TIME ZONE":      "TIME",
		"TIME WITH TIME ZONE":         "TIMETZ",
		"VARCHAR":                     "VARCHAR",
	}

	typeName := *r.result.ColumnMetadata[index].TypeName
	if val, ok := columnTypeMapper[strings.ToUpper(typeName)]; ok {
		return val
	}

	// TODO: Replace this with return "VARCHAR" once this ds is no longer in development
	panic("could not map redshift type to go sql type")
}

// Close closes the rows iterator.
func (r *Rows) Close() error {
	r.done = true
	return nil
}


// fetchNextPage fetches the next statement result page and adds the result to the row
func (r *Rows) fetchNextPage(token *string) error {
	var err error

	r.result, err = r.client.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
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
		switch strings.ToUpper(*col.TypeName) {
		case "INT2":
			ret[i] = int16(*curr.LongValue)
		case "INT", "INT4":
			ret[i] = int32(*curr.LongValue)
		case "INT8":
			ret[i] = *curr.LongValue
		case "FLOAT4":
			ret[i] = float32(*curr.LongValue)
		case "NUMERIC", "FLOAT", "FLOAT8":
			v, err := strconv.ParseFloat(*curr.StringValue, 64)
			if err != nil {
				return err
			}
			ret[i] = v
		case "BOOL":
			// don't know why boolean values are not passed as curr.BooleanValue
			boolValue, err := strconv.ParseBool(*curr.StringValue)
			if err != nil {
				return err
			}
			ret[i] = boolValue
		case "VARCHAR", "CHARACTER", "NCHAR", "BPCHAR", "VARYING", "NVARCHAR", "TEXT":
			ret[i] = *curr.StringValue
		case "TIMESTAMP", "TIMESTAMPTZ":
			// TODO: Replace this with something more robust
			t, err := dateparse.ParseAny(*curr.StringValue)
			if err != nil {
				return err
			}
			ret[i] = t
		default:
			// Unhandled type should probably be handled like this: ret[i] = *curr.StringValue
			// But while this driver is still in development, let's panic so we get a chance to add them.
			panic("unhandled type name")
		}
	}

	return nil
}
