package driver

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice/redshiftdataapiserviceiface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/jpillora/backoff"
)

var (
	backoffMin = 200 * time.Millisecond
	backoffMax = 10 * time.Minute
)

type conn struct {
	sessionCache *awsds.SessionCache
	settings     *models.RedshiftDataSourceSettings
}

func newConnection(sessionCache *awsds.SessionCache, settings *models.RedshiftDataSourceSettings) *conn {
	return &conn{
		sessionCache: sessionCache,
		settings:     settings,
	}
}

func parseStatementInput(query string, settings *models.RedshiftDataSourceSettings) *redshiftdataapiservice.ExecuteStatementInput {
	statementInput := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(settings.ClusterIdentifier),
		Database:          aws.String(settings.Database),
		Sql:               aws.String(query),
	}
	if settings.ManagedSecret.ARN != "" {
		statementInput.SecretArn = aws.String(settings.ManagedSecret.ARN)
	} else {
		statementInput.DbUser = aws.String(settings.DBUser)
	}
	return statementInput
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	session, err := c.sessionCache.GetSession(c.settings.DefaultRegion, c.settings.AWSDatasourceSettings)
	if err != nil {
		return nil, err
	}
	client := redshiftdataapiservice.New(session)

	statementInput := parseStatementInput(query, c.settings)
	executeStatementResult, err := client.ExecuteStatement(statementInput)
	if err != nil {
		return nil, err
	}

	if err := c.waitOnQuery(ctx, client, *executeStatementResult.Id); err != nil {
		return nil, err
	}

	return newRows(client, *executeStatementResult.Id)
}

// waitOnQuery polls the redshift api until the query finishes, returning an error if it failed.
func (c *conn) waitOnQuery(ctx context.Context, service redshiftdataapiserviceiface.RedshiftDataAPIServiceAPI, queryID string) error {
	backoffInstance := backoff.Backoff{
		Min:    backoffMin,
		Max:    backoffMax,
		Factor: 2,
	}
	for {
		statusResp, err := service.DescribeStatementWithContext(ctx, &redshiftdataapiservice.DescribeStatementInput{
			Id: aws.String(queryID),
		})
		if err != nil {
			return err
		}

		switch *statusResp.Status {
		case redshiftdataapiservice.StatusStringFailed,
			redshiftdataapiservice.StatusStringAborted:
			reason := *statusResp.Error
			return errors.New(reason)
		case redshiftdataapiservice.StatusStringFinished:
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoffInstance.Duration()):
			continue
		}
	}
}

func (c *conn) Ping(ctx context.Context) error {
	rows, err := c.QueryContext(ctx, "SELECT 1", nil)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("redshift driver doesn't support begin statements")
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("redshift driver doesn't support prepared statements")
}

func (c *conn) Close() error {
	return nil
}
