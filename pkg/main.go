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
	// Start listening to requests sent from Grafana.
	s := redshift.New()
	ds := awsds.NewAsyncAWSDatasource(s)
	ds.Completable = s
	ds.CustomRoutes = routes.New(s).Routes()

	// newDatasourceForUpgradedPluginSdk adds context to the NewDatasource function, which is the new signature expected
	// for an InstanceFactoryFunc in grafana-plugin-sdk-go
	// TODO: Remove this function and create a NewDatasource method with Context in grafana-aws-sdk, see https://github.com/grafana/oss-plugin-partnerships/issues/648
	newDatasourceForUpgradedPluginSdk := func(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
		return ds.NewDatasource(settings)
	}

	if err := datasource.Manage(
		"grafana-redshift-datasource",
		newDatasourceForUpgradedPluginSdk,
		datasource.ManageOpts{},
	); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
