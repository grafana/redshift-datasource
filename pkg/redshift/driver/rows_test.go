package driver

import (
	"database/sql/driver"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOnePageSuccess(t *testing.T) {
	redshiftServiceMock := &mockRedshiftService{}
	redshiftServiceMock.calledTimesCountDown = 1
	rows, rowErr := newRows(redshiftServiceMock, singlePageResponseQueryId)
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
	redshiftServiceMock := &mockRedshiftService{}
	redshiftServiceMock.calledTimesCountDown = 5
	rows, rowErr := newRows(redshiftServiceMock, multiPageResponseQueryId)
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
	require.Equal(t, 5, redshiftServiceMock.calledTimesCounter)
}
