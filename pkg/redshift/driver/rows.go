package driver

import (
	"database/sql/driver"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

type Rows struct {
	client  *redshiftdataapiservice.RedshiftDataAPIService
	queryID string

	done          bool
	result           *redshiftdataapiservice.GetStatementResultOutput
	pageCount       int64
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


func (r *Rows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}

	// If nothing left to iterate...
	if len(r.result.Records) == 0 {
		// And if nothing more to paginate...
		if r.result.NextToken == nil || *r.result.NextToken == "" {
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


func (r *Rows) Columns() []string {
	columnNames := []string{}
	for _, column := range r.result.ColumnMetadata {
		columnNames = append(columnNames, *column.Name)
	}
	return columnNames
}

func (r *Rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	col := *r.result.ColumnMetadata[index]

	if *col.Nullable == 1 {
		return true, true
	}

	return false, true
}

func (r *Rows)  ColumnTypeScanType(index int) reflect.Type {
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

	// switch fd.DataTypeOID {
	// case pgtype.Float8OID:
	// 	return reflect.TypeOf(float64(0))
	// case pgtype.Float4OID:
	// 	return reflect.TypeOf(float32(0))
	// case pgtype.Int8OID:
	// 	return reflect.TypeOf(int64(0))
	// case pgtype.Int4OID:
	// 	return reflect.TypeOf(int32(0))
	// case pgtype.Int2OID:
	// 	return reflect.TypeOf(int16(0))
	// case pgtype.BoolOID:
	// 	return reflect.TypeOf(false)
	// case pgtype.NumericOID:
	// 	return reflect.TypeOf(float64(0))
	// case pgtype.DateOID, pgtype.TimestampOID, pgtype.TimestamptzOID:
	// 	return reflect.TypeOf(time.Time{})
	// case pgtype.ByteaOID:
	// 	return reflect.TypeOf([]byte(nil))
	// default:
	// 	return reflect.TypeOf("")
	// }
}

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
		"TIMESTAMP WITHOUT TIME ZONE": "TIMESTAMP",
		"TIMESTAMP WITH TIME ZONE":    "TIMESTAMPTZ",
		"TIME WITHOUT TIME ZONE":      "TIME",
		"TIME WITH TIME ZONE":         "TIMETZ",
	}

	typeName := *r.result.ColumnMetadata[index].TypeName
	// TODO: Create mapping between redshift and go.SQL
	if val, ok := columnTypeMapper[strings.ToUpper(typeName)]; ok {
		return val
	}
	return "VARCHAR"
}

func (r *Rows) Close() error {
	r.done = true
	return nil
}

func (r *Rows) fetchNextPage(token *string) error {
	var err error

	r.result, err = r.client.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
		Id:        aws.String(r.queryID),
		NextToken: token,
	})

	if err != nil {
		return err
	}

	r.pageCount++

	// r.done = r.result.NextToken == nil
	return nil
}