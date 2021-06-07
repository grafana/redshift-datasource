package driver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

type conn struct {
	sessionCache   *awsds.SessionCache
	settings *models.RedshiftDataSourceSettings
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	panic("not implemented")
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	panic("not implemented")
}

func (c *conn) Ping(ctx context.Context) error {
	const testQuery = "SELECT 1"

	session, err := c.sessionCache.GetSession(c.settings.DefaultRegion, c.settings.AWSDatasourceSettings)
	if err != nil {
		return err
	}

	client := redshiftdataapiservice.New(session)
	statementInput := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(c.settings.ClusterIdentifier),
		Database: aws.String(c.settings.Database),
		DbUser: aws.String(c.settings.DBUser),
		Sql	: aws.String(testQuery),
	}
	executeStatementResult, err := client.ExecuteStatement(statementInput)
	if err != nil {
		return err
	}

	// wait for a second so that the statement gets a chance to finish before querying the statement result. 
	// this will be replace by something non-blocking eventually
	time.Sleep(1 * time.Second)

	statementResult, err := client.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
		Id: executeStatementResult.Id,
	})

	log.DefaultLogger.Info("healthcheck", "statementResult", statementResult.TotalNumRows)

	if err != nil {
		describeStatementResult, _ := client.DescribeStatement(&redshiftdataapiservice.DescribeStatementInput{
			Id: executeStatementResult.Id,
		})
		return fmt.Errorf(*describeStatementResult.Error)
	}

	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	panic("not implemented")
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	panic("Athena doesn't support prepared statements")
}

func (c *conn) Close() error {
	return nil
}
