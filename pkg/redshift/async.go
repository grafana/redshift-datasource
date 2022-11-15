package redshift

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/pkg/errors"
	"github.com/sunker/async-datasource/pkg/asyncds"
)

type RedshiftAsyncQueryData struct {
	ds *Datasource
}

func (async RedshiftAsyncQueryData) StartQuery(ctx context.Context, query backend.DataQuery) (string, error) {
	q, err := models.GetQuery(query)
	if err != nil {
		return "", err
	}
	service, err := async.ds.GetClient(async.ds.Settings.Region)
	if err != nil {
		return "", err
	}
	redshiftInput := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(async.ds.Settings.ClusterIdentifier),
		Database:          aws.String(async.ds.Settings.Database),
		DbUser:            aws.String(async.ds.Settings.DBUser),
		Sql:               aws.String(q.RawSQL),
		WithEvent:         aws.Bool(async.ds.Settings.WithEvent),
	}

	output, err := service.ExecuteStatementWithContext(ctx, redshiftInput)
	if err != nil {
		return "", fmt.Errorf("%w: %v", api.ExecuteError, err)
	}
	return *output.Id, nil
}
func (async RedshiftAsyncQueryData) GetQueryID(ctx context.Context, query backend.DataQuery) (string, error) {
	return "", nil
}
func (async RedshiftAsyncQueryData) GetQueryStatus(ctx context.Context, queryId string) (asyncds.QueryStatus, error) {
	service, err := async.ds.GetClient(async.ds.Settings.Region)
	if err != nil {
		return asyncds.QueryUnknown, err
	}
	statusResp, err := service.DescribeStatementWithContext(ctx, &redshiftdataapiservice.DescribeStatementInput{
		Id: aws.String(queryId),
	})
	if err != nil {
		return asyncds.QueryUnknown, fmt.Errorf("%w: %v", api.StatusError, err)
	}

	var status asyncds.QueryStatus
	switch *statusResp.Status {
	case redshiftdataapiservice.StatusStringFailed,
		redshiftdataapiservice.StatusStringAborted:
		status = asyncds.QueryFailed
		err = errors.New(*statusResp.Error)
	case redshiftdataapiservice.StatusStringFinished:
		status = asyncds.QueryFinished
	default:
		status = asyncds.QueryRunning
	}

	return status, err

}
func (async RedshiftAsyncQueryData) CancelQuery(ctx context.Context, queryId string) error {
	service, err := async.ds.GetClient(async.ds.Settings.Region)
	if err != nil {
		return fmt.Errorf("%w: %v", api.ExecuteError, err)
	}
	_, err = service.CancelStatement(&redshiftdataapiservice.CancelStatementInput{
		Id: aws.String(queryId),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", err, api.StopError)
	}
	return nil
}
func (async RedshiftAsyncQueryData) GetResult(ctx context.Context, queryId string) (data.Frames, error) {
	service, err := async.ds.GetClient(async.ds.Settings.Region)
	if err != nil {
		return data.Frames{}, fmt.Errorf("%w: %v", api.ExecuteError, err)
	}
	//TODO: add pagination supports
	result, err := service.GetStatementResult(&redshiftdataapiservice.GetStatementResultInput{
		Id: aws.String(queryId),
	})
	if err != nil {
		return data.Frames{}, fmt.Errorf("%w: %v", api.ExecuteError, err)
	}

	frames, err := QueryResultToDataFrame(queryId, result)
	if err != nil {
		return data.Frames{}, fmt.Errorf("%w: %v", api.ExecuteError, err)
	}
	return frames, nil
}
