package driver

import (
	"context"
	"database/sql/driver"
	"fmt"

	redshiftdataV2types "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
)

var _ awsds.AsyncDB = &db{}

// Implements AsyncDB
type db struct {
	api    *api.API
	closed bool
}

func newDB(api *api.API) *db {
	return &db{
		api: api,
	}
}

func (d *db) StartQuery(ctx context.Context, query string, _ ...interface{}) (string, error) {
	output, err := d.api.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: query})
	if err != nil {
		return "", err
	}
	return output.ID, nil
}

func (d *db) GetQueryID(ctx context.Context, query string, args ...interface{}) (bool, string, error) {
	return d.api.GetQueryID(ctx, query, args)
}

func (d *db) QueryStatus(ctx context.Context, queryID string) (awsds.QueryStatus, error) {
	status, err := d.api.Status(ctx, &sqlAPI.ExecuteQueryOutput{ID: queryID})
	if err != nil {
		return awsds.QueryUnknown, err
	}
	var returnStatus awsds.QueryStatus
	switch redshiftdataV2types.StatusString(status.State) {
	case redshiftdataV2types.StatusStringSubmitted,
		redshiftdataV2types.StatusStringPicked:
		returnStatus = awsds.QuerySubmitted
	case redshiftdataV2types.StatusStringStarted:
		returnStatus = awsds.QueryRunning
	case redshiftdataV2types.StatusStringFinished:
		returnStatus = awsds.QueryFinished
	case redshiftdataV2types.StatusStringAborted:
		returnStatus = awsds.QueryCanceled
	case redshiftdataV2types.StatusStringFailed:
		returnStatus = awsds.QueryFailed
	}
	backend.Logger.Debug("QueryStatus", "state", status.State, "queryID", queryID)
	return returnStatus, nil
}

func (d *db) CancelQuery(_ context.Context, queryID string) error {
	return d.api.Stop(&sqlAPI.ExecuteQueryOutput{ID: queryID})
}

func (d *db) GetRows(_ context.Context, queryID string) (driver.Rows, error) {
	return newRows(d.api.DataClient, queryID)
}

func (d *db) Ping(ctx context.Context) error {
	_, err := d.api.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: "SELECT 1"})
	if err != nil {
		return err
	}
	return nil
}

func (d *db) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("redshift driver doesn't support begin statements")
}

func (d *db) Prepare(_ string) (driver.Stmt, error) {
	return nil, fmt.Errorf("redshift driver doesn't support prepared statements")
}

func (d *db) Close() error {
	d.closed = true
	return nil
}
