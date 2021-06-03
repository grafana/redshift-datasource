package redshift

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

var (
	_ backend.QueryDataHandler      = (*RedshiftDatasource)(nil)
	_ backend.CheckHealthHandler    = (*RedshiftDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*RedshiftDatasource)(nil)
)

// RedshiftDatasource is the redshift data source backend host
type RedshiftDatasource struct{
	settings *models.RedshiftDataSourceSettings
	sessionCache *awsds.SessionCache
}

// NewRedshiftDatasource creates a new datasource instance.
func NewRedshiftDatasource(dataSourceInstanceSettings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings := &models.RedshiftDataSourceSettings{}
	err := settings.Load(dataSourceInstanceSettings)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}
	return &RedshiftDatasource{
		settings: settings,
		sessionCache: awsds.NewSessionCache(),
	}, nil
}

// Dispose cleans up datasource instance resources.
func (d *RedshiftDatasource) Dispose() {
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *RedshiftDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
}

func (d *RedshiftDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	frame := data.NewFrame("response")

	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("values", nil, []int64{10, 20}),
	)

	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *RedshiftDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("CheckHealth called", "request", req)
	const testQuery = "SELECT 1"

	session, err := d.sessionCache.GetSession(d.settings.DefaultRegion, d.settings.AWSDatasourceSettings)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	client := redshiftdataapiservice.New(session)
	statementInput := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: aws.String(d.settings.ClusterIdentifier),
		Database: aws.String(d.settings.Database),
		DbUser: aws.String(d.settings.DBUser),
		Sql	: aws.String(testQuery),
	}
	executeStatementResult, err := client.ExecuteStatement(statementInput)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
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
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: *describeStatementResult.Error,
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
