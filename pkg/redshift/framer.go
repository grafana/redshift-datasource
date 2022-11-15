package redshift

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const (
	REDSHIFT_INT                      = "INT"
	REDSHIFT_INT2                     = "INT2"
	REDSHIFT_INT4                     = "INT4"
	REDSHIFT_INT8                     = "INT8"
	REDSHIFT_FLOAT4                   = "FLOAT4"
	REDSHIFT_NUMERIC                  = "NUMERIC"
	REDSHIFT_FLOAT                    = "FLOAT"
	REDSHIFT_FLOAT8                   = "FLOAT8"
	REDSHIFT_BOOL                     = "BOOL"
	REDSHIFT_CHARACTER                = "CHARACTER"
	REDSHIFT_NCHAR                    = "NCHAR"
	REDSHIFT_BPCHAR                   = "BPCHAR"
	REDSHIFT_CHARACTER_VARYING        = "CHARACTER VARYING"
	REDSHIFT_NVARCHAR                 = "NVARCHAR"
	REDSHIFT_TEXT                     = "TEXT"
	REDSHIFT_VARCHAR                  = "VARCHAR"
	REDSHIFT_DATE                     = "DATE"
	REDSHIFT_TIMESTAMP                = "TIMESTAMP"
	REDSHIFT_TIMESTAMP_WITH_TIME_ZONE = "TIMESTAMPTZ"
	REDSHIFT_TIME_WITHOUT_TIME_ZONE   = "TIME"
	REDSHIFT_TIME_WITH_TIME_ZONE      = "TIMETZ"
	REDSHIFT_GEOMETRY                 = "GEOMETRY"
	REDSHIFT_HLLSKETCH                = "HLLSKETCH"
	REDSHIFT_SUPER                    = "SUPER"
	REDSHIFT_NAME                     = "NAME"
)

// QueryResultToDataFrame creates a DataFrame from query results
func QueryResultToDataFrame(refID string, res *redshiftdataapiservice.GetStatementResultOutput) (data.Frames, error) {
	fields := map[int]*data.Field{}
	fieldArr := []*data.Field{}

	// Inspect the column structure
	for index, column := range res.ColumnMetadata {
		typeName := strings.ToUpper(*column.TypeName)
		var fieldType data.FieldType
		switch typeName {
		case REDSHIFT_INT2:
			fieldType = data.FieldTypeNullableInt16
		case REDSHIFT_INT, REDSHIFT_INT4:
			if *column.Name == "time" {
				fieldType = data.FieldTypeTime
			} else {
				fieldType = data.FieldTypeNullableInt32
			}
		case REDSHIFT_INT8:
			fieldType = data.FieldTypeNullableInt64
		case REDSHIFT_NUMERIC, REDSHIFT_FLOAT, REDSHIFT_FLOAT4:
			fieldType = data.FieldTypeNullableFloat64
			if typeName == REDSHIFT_FLOAT4 {
				fieldType = data.FieldTypeNullableFloat32
			}
		case REDSHIFT_FLOAT8:
			if *column.Name == "time" {
				fieldType = data.FieldTypeTime
			} else {
				fieldType = data.FieldTypeNullableFloat64
			}
		case REDSHIFT_BOOL:
			// don't know why boolean values are not passed as curr.BooleanValue
			fieldType = data.FieldTypeNullableBool

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
			fieldType = data.FieldTypeNullableString
		// Time formats from
		// https://docs.aws.amazon.com/redshift/latest/dg/r_Datetime_types.html
		case REDSHIFT_DATE:
			fieldType = data.FieldTypeNullableTime
		case REDSHIFT_TIMESTAMP:
			fieldType = data.FieldTypeNullableTime
		case REDSHIFT_TIMESTAMP_WITH_TIME_ZONE:
			fieldType = data.FieldTypeNullableTime
		case REDSHIFT_TIME_WITHOUT_TIME_ZONE,
			REDSHIFT_TIME_WITH_TIME_ZONE:
			fieldType = data.FieldTypeNullableTime
		default:
			// return fmt.Errorf("unsupported type %s", typeName)
		}
		field := data.NewFieldFromFieldType(fieldType, int(*res.TotalNumRows))
		field.Name = *column.Name
		fields[index] = field
		fieldArr = append(fieldArr, field)
	}

	frame := data.NewFrame(refID, fieldArr...)
	frame.RefID = refID
	frame.Name = refID

	for columnIndex, col := range res.ColumnMetadata {
		for rowIndex, r := range res.Records {
			curr := r[columnIndex]
			var value interface{}

			typeName := strings.ToUpper(*col.TypeName)
			switch typeName {
			case REDSHIFT_INT2:
				value = pointer(int16(*curr.LongValue))
			case REDSHIFT_INT, REDSHIFT_INT4:
				if *col.Name == "time" {
					value = pointer(time.Unix(*curr.LongValue, 0).UTC())
				} else {
					value = pointer(int32(*curr.LongValue))
				}
			case REDSHIFT_INT8:
				value = curr.LongValue
			case REDSHIFT_NUMERIC, REDSHIFT_FLOAT, REDSHIFT_FLOAT4:
				bitSize := 64
				if typeName == REDSHIFT_FLOAT4 {
					bitSize = 32
				}
				v, err := strconv.ParseFloat(*curr.StringValue, bitSize)
				if err != nil {
					return nil, err
				}
				value = pointer(v)
			case REDSHIFT_FLOAT8:
				if *col.Name == "time" {
					value = pointer(time.Unix(int64(*curr.DoubleValue), 0).UTC())
				} else {
					value = curr.DoubleValue
				}
			case REDSHIFT_BOOL:
				// don't know why boolean values are not passed as curr.BooleanValue
				boolValue, err := strconv.ParseBool(*curr.StringValue)
				if err != nil {
					return nil, err
				}
				value = pointer(boolValue)

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
				value = curr.StringValue
			// Time formats from
			// https://docs.aws.amazon.com/redshift/latest/dg/r_Datetime_types.html
			case REDSHIFT_DATE:
				t, err := time.Parse("2006-01-02", *curr.StringValue)
				if err != nil {
					return nil, err
				}
				value = pointer(t)
			case REDSHIFT_TIMESTAMP:
				t, err := time.Parse("2006-01-02 15:04:05", *curr.StringValue)
				if err != nil {
					return nil, err
				}
				value = pointer(t)
			case REDSHIFT_TIMESTAMP_WITH_TIME_ZONE:
				t, err := time.Parse("2006-01-02 15:04:05+00", *curr.StringValue)
				if err != nil {
					return nil, err
				}
				value = pointer(t)
			case REDSHIFT_TIME_WITHOUT_TIME_ZONE,
				REDSHIFT_TIME_WITH_TIME_ZONE:
				t, err := time.Parse("15:04:05", *curr.StringValue)
				if err != nil {
					return nil, err
				}
				value = pointer(t)
			default:
				// return fmt.Errorf("unsupported type %s", typeName)
			}

			fields[columnIndex].Set(rowIndex, value)
		}
	}

	return data.Frames{frame}, nil
}

func pointer[T any](arg T) *T { return &arg }

// typeName := strings.ToUpper(*col.TypeName)
// fieldBuilder := &fieldBuilder{columnIdx: index, name: *columnMeta.Name}
// switch typeName {
// case REDSHIFT_INT2:
// 	fieldType = data.FieldTypeNullableInt16
// case REDSHIFT_INT, REDSHIFT_INT4:
// 	fieldType = data.FieldTypeTime
// case REDSHIFT_INT8:
// 	fieldType = data.FieldTypeNullableInt64
// case REDSHIFT_NUMERIC, REDSHIFT_FLOAT, REDSHIFT_FLOAT4:
// 	fieldType = data.FieldTypeNullableFloat32
// case REDSHIFT_FLOAT8:
// 	fieldType = data.FieldTypeNullableFloat64
// case REDSHIFT_BOOL:
// 	// don't know why boolean values are not passed as curr.BooleanValue
// 	fieldType = data.FieldTypeNullableBool

// case REDSHIFT_CHARACTER,
// 	REDSHIFT_VARCHAR,
// 	REDSHIFT_NCHAR,
// 	REDSHIFT_BPCHAR,
// 	REDSHIFT_CHARACTER_VARYING,
// 	REDSHIFT_NVARCHAR,
// 	REDSHIFT_TEXT,
// 	// Complex types are returned as a string
// 	REDSHIFT_GEOMETRY,
// 	REDSHIFT_HLLSKETCH,
// 	REDSHIFT_SUPER,
// 	REDSHIFT_NAME:
// 	fieldType = data.FieldTypeNullableString
// // Time formats from
// // https://docs.aws.amazon.com/redshift/latest/dg/r_Datetime_types.html
// case REDSHIFT_DATE:
// 	fieldType = data.FieldTypeNullableTime
// case REDSHIFT_TIMESTAMP:
// 	fieldType = data.FieldTypeNullableTime
// case REDSHIFT_TIMESTAMP_WITH_TIME_ZONE:
// 	fieldType = data.FieldTypeNullableTime
// case REDSHIFT_TIME_WITHOUT_TIME_ZONE,
// 	REDSHIFT_TIME_WITH_TIME_ZONE:
// 	fieldType = data.FieldTypeNullableTime
// default:
// 	// return fmt.Errorf("unsupported type %s", typeName)
// }
// builders = append(builders, fieldBuilder)
