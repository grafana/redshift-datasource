package driver

import (
	"database/sql/driver"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

func convertRow(columns []*redshiftdataapiservice.ColumnMetadata, data []*redshiftdataapiservice.Field, ret []driver.Value) error {
	for i, curr := range data {
		col := columns[i]

		var value interface{}
		if curr.BlobValue != nil {
			value = curr.BlobValue
		} else if curr.BooleanValue != nil {
			value = *curr.BooleanValue
		} else if curr.DoubleValue != nil {
			value = *curr.DoubleValue
		} else if curr.LongValue != nil {
			value = *curr.LongValue
		} else if curr.StringValue != nil {
			if *col.TypeName == "timestamp" {
				t, err := dateparse.ParseAny(*curr.StringValue)
				if err != nil {
					return err
				}
				value = t
			} else {
				value = *curr.StringValue
			}
		} else if *curr.IsNull {
			value = nil
		}
		
		// switch strings.ToUpper(*col.TypeName) {
		// 	case "INT2": 
		// 		v, err := strconv.ParseInt(value, 10, 64)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		return &v, nil
		// 		return reflect.TypeOf(int16(0))
		// 	case "INT", "INT4": 
		// 		return reflect.TypeOf(int32(0))
		// 	case "INT8": 
		// 		return reflect.TypeOf(int64(0))
		// 	case "FLOAT4":
		// 		return reflect.TypeOf(float32(0))
		// 	case "NUMERIC", "FLOAT", "FLOAT8":
		// 		return reflect.TypeOf(float64(0))
		// 	case "BOOL":
		// 		return reflect.TypeOf(false)
		// 	case "CHARACTER", "NCHAR", "BPCHAR", "VARYING", "NVARCHAR", "TEXT":
		// 		return reflect.TypeOf("")
		// 	case "TIMESTAMP, ":
		// 		return reflect.TypeOf(time.Time{})
		// 	default:
		// 		return reflect.TypeOf("")
		// 	}

			ret[i] = value

		
	}


	return nil
}
