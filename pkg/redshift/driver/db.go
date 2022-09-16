package driver

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
)

var _ awsds.AsyncDB = &DB{}

// Implements AsyncDB
type DB struct {
	api    *api.API
	closed bool
}

func (d *DB) StartQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
	output, err := d.api.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: query})
	if err != nil {
		return "", err
	}
	return output.ID, nil
}

func (d *DB) GetQueryID(ctx context.Context, query string, args ...interface{}) (bool, string, error) {
	return d.api.GetQueryID(ctx, query, args)
}

func (d *DB) QueryStatus(ctx context.Context, queryID string) (awsds.QueryStatus, error) {
	status, err := d.api.Status(ctx, &sqlAPI.ExecuteQueryOutput{ID: queryID})
	if err != nil {
		return awsds.QueryUnknown, err
	}
	var returnStatus awsds.QueryStatus
	switch status.State {
	case "SUBMITTED", "PICKED":
		returnStatus = awsds.QuerySubmitted
	case "STARTED":
		returnStatus = awsds.QueryRunning
	case "FINISHED":
		returnStatus = awsds.QueryFinished
	case "ABORTED":
		returnStatus = awsds.QueryCanceled
	case "FAILED":
		returnStatus = awsds.QueryFailed
	}
	return returnStatus, nil

}

func (d *DB) CancelQuery(ctx context.Context, queryID string) error {
	return d.api.Stop(&sqlAPI.ExecuteQueryOutput{ID: queryID})

}

func (d *DB) GetRows(ctx context.Context, queryID string) (driver.Rows, error) {
	return newRows(d.api.DataClient, queryID)

}

func (d *DB) Ping(ctx context.Context) error {
	_, err := d.api.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: "SELECT 1"})
	if err != nil {
		return err
	}
	return nil

}

func (d *DB) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("redshift driver doesn't support begin statements")

}

func (d *DB) Prepare(query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("redshift driver doesn't support prepared statements")

}

func (d *DB) Close() error {
	d.closed = true
	return nil

}
