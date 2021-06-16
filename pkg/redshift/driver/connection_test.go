package driver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var waitOnQueryTestCases = []struct {
	calledTimesCountDown int
	statementStatus      string
	err                  error
}{
	{1, DESCRIBE_STATEMENT_SUCCEEDED, nil},
	{10, DESCRIBE_STATEMENT_SUCCEEDED, nil},
	{1, DESCRIBE_STATEMENT_FAILED, fmt.Errorf(DESCRIBE_STATEMENT_FAILED)},
	{10, DESCRIBE_STATEMENT_FAILED, fmt.Errorf(DESCRIBE_STATEMENT_FAILED)},
}

func TestConnection_waitOnQuery(t *testing.T) {
	t.Parallel()
	for _, tc := range waitOnQueryTestCases {
		c := &conn{pollingInterval: time.Millisecond}
		redshiftServiceMock := newMockRedshiftService()
		redshiftServiceMock.calledTimesCountDown = tc.calledTimesCountDown
		err := c.waitOnQuery(context.Background(), redshiftServiceMock, tc.statementStatus)
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.calledTimesCountDown, redshiftServiceMock.calledTimesCounter)
	}
}

func TestConnection_waitOnQuery_success2(t *testing.T) {
	t.Parallel()
	c := &conn{pollingInterval: time.Millisecond}
	redshiftServiceMock := newMockRedshiftService()
	redshiftServiceMock.calledTimesCountDown = 10
	err := c.waitOnQuery(context.Background(), redshiftServiceMock, DESCRIBE_STATEMENT_SUCCEEDED)
	assert.Nil(t, err)
	assert.Equal(t, 10, redshiftServiceMock.calledTimesCounter)
}
