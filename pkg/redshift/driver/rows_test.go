package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	redshiftdatatypes "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"

	mock "github.com/grafana/redshift-datasource/pkg/redshift/driver/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnePageSuccess(t *testing.T) {
	redshiftServiceMock := &mock.RedshiftService{}
	redshiftServiceMock.CalledTimesCountDown = 1
	rows, rowErr := newRows(context.Background(), redshiftServiceMock, mock.SinglePageResponseQueryId)
	require.NoError(t, rowErr)
	cnt := 0
	for {
		var col1, col2 string
		err := rows.Next([]driver.Value{
			&col1,
			&col2,
		})
		if err != nil {
			require.ErrorIs(t, io.EOF, err)
			break
		}
		require.NoError(t, err)
		cnt++
	}
	require.Equal(t, 2, cnt)
}

func TestMultiPageSuccess(t *testing.T) {
	redshiftServiceMock := &mock.RedshiftService{}
	redshiftServiceMock.CalledTimesCountDown = 5
	rows, rowErr := newRows(context.Background(), redshiftServiceMock, mock.MultiPageResponseQueryId)
	require.NoError(t, rowErr)
	cnt := 0
	for {
		var col1, col2 string
		err := rows.Next([]driver.Value{
			&col1,
			&col2,
		})
		if err != nil {
			require.ErrorIs(t, io.EOF, err)
			break
		}
		require.NoError(t, err)
		cnt++
	}
	require.Equal(t, 10, cnt)
	require.Equal(t, 5, redshiftServiceMock.CalledTimesCounter)
}

func Test_convertRow(t *testing.T) {

	tests := []struct {
		name          string
		metadata      redshiftdatatypes.ColumnMetadata
		data          redshiftdatatypes.Field
		expectedType  string
		expectedValue string
	}{
		{
			name: "numeric type int",
			metadata: redshiftdatatypes.ColumnMetadata{
				Name:     aws.String("num"),
				TypeName: aws.String(REDSHIFT_INT),
			},
			data:          &redshiftdatatypes.FieldMemberLongValue{Value: 1},
			expectedType:  "int32",
			expectedValue: "1",
		},
		{
			name: "numeric type int2",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_INT2),
			},
			data:          &redshiftdatatypes.FieldMemberLongValue{Value: 2},
			expectedType:  "int16",
			expectedValue: "2",
		},
		{
			name: "numeric type int4",
			metadata: redshiftdatatypes.ColumnMetadata{
				Name:     aws.String("num"),
				TypeName: aws.String(REDSHIFT_INT4),
			},
			data:          &redshiftdatatypes.FieldMemberLongValue{Value: 3},
			expectedType:  "int32",
			expectedValue: "3",
		},
		{
			name: "time as int4",
			metadata: redshiftdatatypes.ColumnMetadata{
				Name:     aws.String("time"),
				TypeName: aws.String(REDSHIFT_INT4),
			},
			data: &redshiftdatatypes.FieldMemberLongValue{
				Value: 1624741200,
			},
			expectedType:  "time.Time",
			expectedValue: "2021-06-26 21:00:00 +0000 UTC",
		},
		{
			name: "numeric type int8",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_INT8),
			},
			data:          &redshiftdatatypes.FieldMemberLongValue{Value: 4},
			expectedType:  "int64",
			expectedValue: "4",
		},
		{
			name: "numeric type float4",
			metadata: redshiftdatatypes.ColumnMetadata{
				Name:     aws.String("other"),
				TypeName: aws.String(REDSHIFT_FLOAT4),
			},
			data:          &redshiftdatatypes.FieldMemberDoubleValue{Value: 1.1},
			expectedType:  "float64",
			expectedValue: "1.1",
		},
		{
			name: "numeric type numeric",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_NUMERIC),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "1.2"},
			expectedType:  "float64",
			expectedValue: "1.2",
		},
		{
			name: "numeric type float",
			metadata: redshiftdatatypes.ColumnMetadata{
				Name:     aws.String("other"),
				TypeName: aws.String(REDSHIFT_FLOAT),
			},
			data:          &redshiftdatatypes.FieldMemberDoubleValue{Value: 1.3},
			expectedType:  "float64",
			expectedValue: "1.3",
		},
		{
			name: "numeric float8",
			metadata: redshiftdatatypes.ColumnMetadata{
				Name:     aws.String("other"),
				TypeName: aws.String(REDSHIFT_FLOAT8),
			},
			data:          &redshiftdatatypes.FieldMemberDoubleValue{Value: 1.4},
			expectedType:  "float64",
			expectedValue: "1.4",
		},
		{
			name: "bool type",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_BOOL),
			},
			data:          &redshiftdatatypes.FieldMemberBooleanValue{Value: false},
			expectedType:  "bool",
			expectedValue: "false",
		},
		{
			name: "character",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_CHARACTER),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "f"},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "nchar",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_NCHAR),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "f"},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "bpchar",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_BPCHAR),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "f"},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "character varying",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_CHARACTER_VARYING),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "f"},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "text",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TEXT),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "foo"},
			expectedType:  "string",
			expectedValue: "foo",
		},
		{
			name: "varchar",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_VARCHAR),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "foo"},
			expectedType:  "string",
			expectedValue: "foo",
		},
		{
			name: "date",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_DATE),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "2008-01-01"},
			expectedType:  "time.Time",
			expectedValue: "2008-01-01 00:00:00 +0000 UTC",
		},
		{
			name: "timestamp",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIMESTAMP),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "2008-01-01 20:00:00.00"},
			expectedType:  "time.Time",
			expectedValue: "2008-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "timestamp with tz",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIMESTAMP_WITH_TIME_ZONE),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "2008-01-01 20:00:00+00"},
			expectedType:  "time.Time",
			expectedValue: "2008-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "time without tz",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIME_WITHOUT_TIME_ZONE),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "20:00:00.00"},
			expectedType:  "time.Time",
			expectedValue: "0000-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "time with tz",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIME_WITH_TIME_ZONE),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "20:00:00.00"},
			expectedType:  "time.Time",
			expectedValue: "0000-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "geometry",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_GEOMETRY),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: "[B@f69ae81"},
			expectedType:  "string",
			expectedValue: "[B@f69ae81",
		},
		{
			name: "hllsketch",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_HLLSKETCH),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: `{"version":1,"logm":15,"sparse":{"indices":[40242751],"values":[2]}}`},
			expectedType:  "string",
			expectedValue: `{"version":1,"logm":15,"sparse":{"indices":[40242751],"values":[2]}}`,
		},
		{
			name: "super",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_SUPER),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: `{"foo":"bar"}`},
			expectedType:  "string",
			expectedValue: `{"foo":"bar"}`,
		},
		{
			name: "name",
			metadata: redshiftdatatypes.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_NAME),
			},
			data:          &redshiftdatatypes.FieldMemberStringValue{Value: `table`},
			expectedType:  "string",
			expectedValue: `table`,
		},
		{
			name: "unix time",
			metadata: redshiftdatatypes.ColumnMetadata{
				Name:     aws.String("time"),
				TypeName: aws.String(REDSHIFT_FLOAT8),
			},
			data:          &redshiftdatatypes.FieldMemberDoubleValue{Value: 1626357600},
			expectedType:  "time.Time",
			expectedValue: `2021-07-15 14:00:00 +0000 UTC`,
		},
		{
			name:          "null",
			metadata:      redshiftdatatypes.ColumnMetadata{},
			data:          &redshiftdatatypes.FieldMemberIsNull{Value: true},
			expectedType:  "<nil>",
			expectedValue: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := make([]driver.Value, 1)
			err := convertRow(
				[]redshiftdatatypes.ColumnMetadata{tt.metadata},
				[]redshiftdatatypes.Field{tt.data},
				res,
			)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedType, fmt.Sprintf("%T", res[0]))
			assert.Equal(t, tt.expectedValue, fmt.Sprintf("%v", res[0]))
		})
	}

	t.Run("a value followed by a null value", func(t *testing.T) {
		// simulate previous value
		res := []driver.Value{int32(1), int32(2)}

		metadata := []redshiftdatatypes.ColumnMetadata{
			{Name: aws.String("num"), TypeName: aws.String(REDSHIFT_INT)},
			{},
		}
		data := []redshiftdatatypes.Field{
			&redshiftdatatypes.FieldMemberLongValue{Value: 3},
			&redshiftdatatypes.FieldMemberIsNull{Value: true},
		}

		err := convertRow(metadata, data, res)
		require.NoError(t, err)

		expectedValue := []driver.Value{int32(3), nil}
		assert.Equal(t, expectedValue, res)
	})

	t.Run("error returned for missing column type", func(t *testing.T) {
		empty := []redshiftdatatypes.Field{}
		empty = append(empty, &redshiftdatatypes.FieldMemberStringValue{})
		assert.EqualError(t, convertRow(
			[]redshiftdatatypes.ColumnMetadata{{}},
			empty,
			[]driver.Value{},
		), "error in convertRow: col.TypeName is nil")
	})
}
