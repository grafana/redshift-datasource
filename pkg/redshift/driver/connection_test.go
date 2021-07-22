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
	minDuration          int
}{
	{1, redshiftservicemock.DESCRIBE_STATEMENT_SUCCEEDED, nil, 0},
	{2, redshiftservicemock.DESCRIBE_STATEMENT_SUCCEEDED, nil, 1},
	{3, redshiftservicemock.DESCRIBE_STATEMENT_SUCCEEDED, nil, 1 + 1},
	{4, redshiftservicemock.DESCRIBE_STATEMENT_FAILED, fmt.Errorf(redshiftservicemock.DESCRIBE_STATEMENT_FAILED), 1 + 1 + 2},
	{5, redshiftservicemock.DESCRIBE_STATEMENT_FAILED, fmt.Errorf(redshiftservicemock.DESCRIBE_STATEMENT_FAILED), 1 + 1 + 2 + 3},
	{6, redshiftservicemock.DESCRIBE_STATEMENT_FAILED, fmt.Errorf(redshiftservicemock.DESCRIBE_STATEMENT_FAILED), 1 + 1 + 2 + 3 + 5},
}

func TestConnection_waitOnQuery(t *testing.T) {
	t.Parallel()
	for _, tc := range waitOnQueryTestCases {
		c := newConnection(nil, nil)
		redshiftServiceMock := redshiftservicemock.NewMockRedshiftService()
		redshiftServiceMock.CalledTimesCountDown = tc.calledTimesCountDown
		beforeCall := time.Now()
		err := c.waitOnQuery(context.Background(), redshiftServiceMock, tc.statementStatus)
		durationOfCall := time.Since(beforeCall)
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.calledTimesCountDown, redshiftServiceMock.CalledTimesCounter)
		assert.Greater(t, durationOfCall, time.Second * time.Duration(tc.minDuration))
	}
}