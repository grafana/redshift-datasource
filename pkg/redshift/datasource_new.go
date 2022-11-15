package redshift

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice/redshiftdataapiserviceiface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/sunker/async-datasource/pkg/asyncds"
)

type clientGetterFunc func(region string) (srv redshiftdataapiserviceiface.RedshiftDataAPIServiceAPI, err error)

type Datasource struct {
	GetClient clientGetterFunc
	backend.QueryDataHandler
	Settings models.RedshiftDataSourceSettings
}

func NewDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	ds := &Datasource{}
	ds.Settings = models.RedshiftDataSourceSettings{}
	err := ds.Settings.Load(settings)
	if err != nil {
		return nil, err
	}

	sessions := awsds.NewSessionCache()
	ds.GetClient = func(region string) (srv redshiftdataapiserviceiface.RedshiftDataAPIServiceAPI, err error) {
		sess, err := sessions.GetSession(awsds.SessionConfig{
			Settings:      ds.Settings.AWSDatasourceSettings,
			UserAgentName: aws.String("Redshift"),
		})
		if err != nil {
			return nil, err
		}

		return redshiftdataapiservice.New(sess), nil
	}
	ds.QueryDataHandler = asyncds.NewAsyncQueryDataHandler(RedshiftAsyncQueryData{ds: ds})

	return ds, nil
}

// CheckHealth pings the connected SQL database
func (ds *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
