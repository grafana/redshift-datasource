package driver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection_waitOnQuery_failed(t *testing.T) {
	t.Parallel()
	c := &conn{pollingInterval: time.Millisecond}
	redshiftServiceMock := newMockRedshiftService()
	err := c.waitOnQuery(context.Background(), redshiftServiceMock, DESCRIBE_STATEMENT_FAILED)
	assert.NotNil(t, err)
	assert.Equal(t, DESCRIBE_STATEMENT_FAILED, err.Error())
	assert.Equal(t, 1, redshiftServiceMock.calledTimesCounter)
}

func TestConnection_waitOnQuery_(t *testing.T) {
	t.Parallel()
	c := &conn{pollingInterval: time.Millisecond}
	redshiftServiceMock := newMockRedshiftService()
	redshiftServiceMock.calledTimesCountDown = 10
	err := c.waitOnQuery(context.Background(), redshiftServiceMock, DESCRIBE_STATEMENT_SUCCEEDED)
	assert.Nil(t, err)
	assert.Equal(t, 10, redshiftServiceMock.calledTimesCounter)
}
