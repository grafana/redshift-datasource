package main

import (
	"context"
	"os"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/redshift-datasource/pkg/redshift/routes"
)

func main() {
	if err := datasource.Manage(
		"grafana-redshift-datasource",
		MakeDatasourceFactory(),
		datasource.ManageOpts{},
	); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

func MakeDatasourceFactory() datasource.InstanceFactoryFunc {
	return func(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
		log.DefaultLogger.FromContext(ctx).Debug("building new datasource instance")
		s := redshift.New()
		ds := awsds.NewAsyncAWSDatasource(s)
		ds.Completable = s
		ds.CustomRoutes = routes.New(s).Routes()
		ds.EnableRowLimit = true
		return ds.NewDatasource(ctx, settings)
	}
}
