package driver

import (
	"context"
	"fmt"
	"testing"
	"time"

	redshiftservicemock "github.com/grafana/redshift-datasource/pkg/redshift/driver/mock"
	"github.com/stretchr/testify/assert"
)

var waitOnQueryTestCases = []struct {
	calledTimesCountDown int
	statementStatus      string
	err                  error
}{
	{1, redshiftservicemock.DESCRIBE_STATEMENT_SUCCEEDED, nil},
	{10, redshiftservicemock.DESCRIBE_STATEMENT_SUCCEEDED, nil},
	{1, redshiftservicemock.DESCRIBE_STATEMENT_FAILED, fmt.Errorf(redshiftservicemock.DESCRIBE_STATEMENT_FAILED)},
	{10, redshiftservicemock.DESCRIBE_STATEMENT_FAILED, fmt.Errorf(redshiftservicemock.DESCRIBE_STATEMENT_FAILED)},
}

func TestConnection_waitOnQuery(t *testing.T) {
	t.Parallel()
	for _, tc := range waitOnQueryTestCases {
		c := &conn{pollingInterval: time.Millisecond}
		redshiftServiceMock := redshiftservicemock.NewMockRedshiftService()
		redshiftServiceMock.CalledTimesCountDown = tc.calledTimesCountDown
		err := c.waitOnQuery(context.Background(), redshiftServiceMock, tc.statementStatus)
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.calledTimesCountDown, redshiftServiceMock.CalledTimesCounter)
	}
}

func TestConnection_waitOnQuery_success2(t *testing.T) {
	t.Parallel()
	c := &conn{pollingInterval: time.Millisecond}
	redshiftServiceMock := redshiftservicemock.NewMockRedshiftService()
	redshiftServiceMock.CalledTimesCountDown = 10
	err := c.waitOnQuery(context.Background(), redshiftServiceMock, redshiftservicemock.DESCRIBE_STATEMENT_SUCCEEDED)
	assert.Nil(t, err)
	assert.Equal(t, 10, redshiftServiceMock.CalledTimesCounter)
}
