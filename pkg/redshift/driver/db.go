package driver

import (
	"context"
	"database/sql/driver"

	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
	"github.com/grafana/sqlds/v2"
)

type DB struct {
	sqlds.AsyncDB
	API *api.API
}

func (d *DB) StartQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
	output, err := d.API.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: query})
	if err != nil {
		return "", err
	}
	return output.ID, nil
}

func (d *DB) GetQueryID(ctx context.Context, query string, args ...interface{}) (bool, string, error) {
	res, err := d.API.ListStatements(ctx, query)

	if err != nil {
		return false, "", err
	}

	return false, res, nil
}

func (d *DB) QueryStatus(ctx context.Context, queryID string) (bool, string, error) {
	status, err := d.API.Status(ctx, &sqlAPI.ExecuteQueryOutput{ID: queryID})
	if err != nil {
		return false, "", err
	}
	return status.Finished, status.State, nil
}

func (d *DB) CancelQuery(ctx context.Context, queryID string) error {
	return d.API.Stop(&sqlAPI.ExecuteQueryOutput{ID: queryID})
}

func (d *DB) GetRows(ctx context.Context, queryID string) (driver.Rows, error) {
	return newRows(d.API.DataClient, queryID)
}
