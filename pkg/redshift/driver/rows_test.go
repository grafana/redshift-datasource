package driver

import (
	"database/sql/driver"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	redshiftservicemock "github.com/grafana/redshift-datasource/pkg/redshift/driver/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnePageSuccess(t *testing.T) {
	redshiftServiceMock := &redshiftservicemock.RedshiftService{}
	redshiftServiceMock.CalledTimesCountDown = 1
	rows, rowErr := newRows(redshiftServiceMock, redshiftservicemock.SinglePageResponseQueryId)
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
	redshiftServiceMock := &redshiftservicemock.RedshiftService{}
	redshiftServiceMock.CalledTimesCountDown = 5
	rows, rowErr := newRows(redshiftServiceMock, redshiftservicemock.MultiPageResponseQueryId)
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
		metadata      *redshiftdataapiservice.ColumnMetadata
		data          *redshiftdataapiservice.Field
		expectedType  string
		expectedValue string
	}{
		{
			name: "numeric type int",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				Name:     aws.String("num"),
				TypeName: aws.String(REDSHIFT_INT),
			},
			data:          &redshiftdataapiservice.Field{LongValue: aws.Int64(1)},
			expectedType:  "int32",
			expectedValue: "1",
		},
		{
			name: "numeric type int2",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_INT2),
			},
			data: &redshiftdataapiservice.Field{
				LongValue: aws.Int64(2),
			},
			expectedType:  "int16",
			expectedValue: "2",
		},
		{
			name: "numeric type int4",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				Name:     aws.String("num"),
				TypeName: aws.String(REDSHIFT_INT4),
			},
			data: &redshiftdataapiservice.Field{
				LongValue: aws.Int64(3),
			},
			expectedType:  "int32",
			expectedValue: "3",
		},
		{
			name: "time as int4",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				Name:     aws.String("time"),
				TypeName: aws.String(REDSHIFT_INT4),
			},
			data: &redshiftdataapiservice.Field{
				LongValue: aws.Int64(1624741200),
			},
			expectedType:  "time.Time",
			expectedValue: "2021-06-26 21:00:00 +0000 UTC",
		},
		{
			name: "numeric type int8",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_INT8),
			},
			data: &redshiftdataapiservice.Field{
				LongValue: aws.Int64(4),
			},
			expectedType:  "int64",
			expectedValue: "4",
		},
		{
			name: "numeric type float4",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				Name:     aws.String("other"),
				TypeName: aws.String(REDSHIFT_FLOAT4),
			},
			data: &redshiftdataapiservice.Field{
				DoubleValue: aws.Float64(1.1),
			},
			expectedType:  "float64",
			expectedValue: "1.1",
		},
		{
			name: "numeric type numeric",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_NUMERIC),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("1.2"),
			},
			expectedType:  "float64",
			expectedValue: "1.2",
		},
		{
			name: "numeric type float",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				Name:     aws.String("other"),
				TypeName: aws.String(REDSHIFT_FLOAT),
			},
			data: &redshiftdataapiservice.Field{
				DoubleValue: aws.Float64(1.3),
			},
			expectedType:  "float64",
			expectedValue: "1.3",
		},
		{
			name: "numeric float8",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				Name:     aws.String("other"),
				TypeName: aws.String(REDSHIFT_FLOAT8),
			},
			data: &redshiftdataapiservice.Field{
				DoubleValue: aws.Float64(1.4),
			},
			expectedType:  "float64",
			expectedValue: "1.4",
		},
		{
			name: "bool type",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_BOOL),
			},
			data: &redshiftdataapiservice.Field{
				BooleanValue: aws.Bool(false),
			},
			expectedType:  "bool",
			expectedValue: "false",
		},
		{
			name: "character",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_CHARACTER),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("f"),
			},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "nchar",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_NCHAR),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("f"),
			},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "bpchar",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_BPCHAR),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("f"),
			},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "character varying",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_CHARACTER_VARYING),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("f"),
			},
			expectedType:  "string",
			expectedValue: "f",
		},
		{
			name: "text",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TEXT),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("foo"),
			},
			expectedType:  "string",
			expectedValue: "foo",
		},
		{
			name: "varchar",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_VARCHAR),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("foo"),
			},
			expectedType:  "string",
			expectedValue: "foo",
		},
		{
			name: "date",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_DATE),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("2008-01-01"),
			},
			expectedType:  "time.Time",
			expectedValue: "2008-01-01 00:00:00 +0000 UTC",
		},
		{
			name: "timestamp",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIMESTAMP),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("2008-01-01 20:00:00.00"),
			},
			expectedType:  "time.Time",
			expectedValue: "2008-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "timestamp with tz",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIMESTAMP_WITH_TIME_ZONE),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("2008-01-01 20:00:00+00"),
			},
			expectedType:  "time.Time",
			expectedValue: "2008-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "time without tz",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIME_WITHOUT_TIME_ZONE),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("20:00:00.00"),
			},
			expectedType:  "time.Time",
			expectedValue: "0000-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "time with tz",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_TIME_WITH_TIME_ZONE),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("20:00:00.00"),
			},
			expectedType:  "time.Time",
			expectedValue: "0000-01-01 20:00:00 +0000 UTC",
		},
		{
			name: "geometry",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_GEOMETRY),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String("[B@f69ae81"),
			},
			expectedType:  "string",
			expectedValue: "[B@f69ae81",
		},
		{
			name: "hllsketch",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_HLLSKETCH),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String(`{"version":1,"logm":15,"sparse":{"indices":[40242751],"values":[2]}}`),
			},
			expectedType:  "string",
			expectedValue: `{"version":1,"logm":15,"sparse":{"indices":[40242751],"values":[2]}}`,
		},
		{
			name: "super",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_SUPER),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String(`{"foo":"bar"}`),
			},
			expectedType:  "string",
			expectedValue: `{"foo":"bar"}`,
		},
		{
			name: "name",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				TypeName: aws.String(REDSHIFT_NAME),
			},
			data: &redshiftdataapiservice.Field{
				StringValue: aws.String(`table`),
			},
			expectedType:  "string",
			expectedValue: `table`,
		},
		{
			name: "unix time",
			metadata: &redshiftdataapiservice.ColumnMetadata{
				Name:     aws.String("time"),
				TypeName: aws.String(REDSHIFT_FLOAT8),
			},
			data: &redshiftdataapiservice.Field{
				DoubleValue: aws.Float64(1626357600),
			},
			expectedType:  "time.Time",
			expectedValue: `2021-07-15 14:00:00 +0000 UTC`,
		},
		{
			name:     "null",
			metadata: &redshiftdataapiservice.ColumnMetadata{},
			data: &redshiftdataapiservice.Field{
				IsNull: aws.Bool(true),
			},
			expectedType:  "<nil>",
			expectedValue: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := make([]driver.Value, 1)
			err := convertRow(
				[]*redshiftdataapiservice.ColumnMetadata{tt.metadata},
				[]*redshiftdataapiservice.Field{tt.data},
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

		metadata := []*redshiftdataapiservice.ColumnMetadata{
			{Name: aws.String("num"), TypeName: aws.String(REDSHIFT_INT)},
			{},
		}
		data := []*redshiftdataapiservice.Field{
			{LongValue: aws.Int64(3)},
			{IsNull: aws.Bool(true)},
		}

		err := convertRow(metadata, data, res)
		require.NoError(t, err)

		expectedValue := []driver.Value{int32(3), nil}
		assert.Equal(t, expectedValue, res)
	})

	t.Run("error returned for missing column type", func(t *testing.T) {
		assert.EqualError(t, convertRow(
			[]*redshiftdataapiservice.ColumnMetadata{{}},
			[]*redshiftdataapiservice.Field{{}},
			[]driver.Value{},
		), "error in convertRow: col.TypeName is nil")
	})
}
