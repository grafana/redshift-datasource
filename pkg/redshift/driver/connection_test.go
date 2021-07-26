package driver

import (
	"context"
	"fmt"
	"testing"
	"time"

	redshiftservicemock "github.com/grafana/redshift-datasource/pkg/redshift/driver/mock"
	"github.com/jpillora/backoff"
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
		// for tests we override backoff instance to always take 1 millisecond so the tests run quickly
		c := &conn{backoffInstance: backoff.Backoff{
			Min:    1 * time.Millisecond,
			Max:   1 * time.Millisecond,
		},}
		redshiftServiceMock := redshiftservicemock.NewMockRedshiftService()
		redshiftServiceMock.CalledTimesCountDown = tc.calledTimesCountDown
		err := c.waitOnQuery(context.Background(), redshiftServiceMock, tc.statementStatus)
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.calledTimesCountDown, redshiftServiceMock.CalledTimesCounter)
	}
}