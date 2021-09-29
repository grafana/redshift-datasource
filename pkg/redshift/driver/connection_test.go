package driver

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/google/go-cmp/cmp"
	redshiftservicemock "github.com/grafana/redshift-datasource/pkg/redshift/driver/mock"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
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
	backoffMin = 1 * time.Millisecond
	backoffMax = 1 * time.Millisecond

	for _, tc := range waitOnQueryTestCases {
		// for tests we override backoff instance to always take 1 millisecond so the tests run quickly
		c := &conn{}
		redshiftServiceMock := redshiftservicemock.NewMockRedshiftService()
		redshiftServiceMock.CalledTimesCountDown = tc.calledTimesCountDown
		err := c.waitOnQuery(context.Background(), redshiftServiceMock, tc.statementStatus)
		assert.Equal(t, tc.err, err)
		assert.Equal(t, tc.calledTimesCountDown, redshiftServiceMock.CalledTimesCounter)
	}
}

func Test_parseStatementInput(t *testing.T) {
	tests := []struct {
		description string
		query       string
		settings    *models.RedshiftDataSourceSettings
		expected    *redshiftdataapiservice.ExecuteStatementInput
	}{
		{
			"using temporary creds",
			"select * from table",
			&models.RedshiftDataSourceSettings{
				ClusterIdentifier: "cluster",
				Database:          "db",
				DBUser:            "user",
			},
			&redshiftdataapiservice.ExecuteStatementInput{
				ClusterIdentifier: aws.String("cluster"),
				Database:          aws.String("db"),
				Sql:               aws.String("select * from table"),
				DbUser:            aws.String("user"),
			},
		},
		{
			"using managed secret",
			"select * from table",
			&models.RedshiftDataSourceSettings{
				ClusterIdentifier: "cluster",
				Database:          "db",
				ManagedSecret:     "arn:...",
				// ignored
				DBUser: "user",
			},
			&redshiftdataapiservice.ExecuteStatementInput{
				ClusterIdentifier: aws.String("cluster"),
				Database:          aws.String("db"),
				Sql:               aws.String("select * from table"),
				SecretArn:         aws.String("arn:..."),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			res := parseStatementInput(tt.query, tt.settings)
			if !cmp.Equal(res, tt.expected) {
				t.Errorf("unexpected result: %v", cmp.Diff(res, tt.expected))
			}
		})
	}
}
