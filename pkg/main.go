package main

import (
	"os"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/redshift-datasource/pkg/redshift/routes"
)

func main() {
	// Start listening to requests sent from Grafana.
	s := redshift.New()
	ds := awsds.NewAsyncAWSDatasource(s)
	ds.SqlDatasource.Completable = s
	ds.SqlDatasource.CustomRoutes = routes.New(s).Routes()

	if err := datasource.Manage(
		"grafana-redshift-datasource",
		ds.NewDatasource,
		datasource.ManageOpts{},
	); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
