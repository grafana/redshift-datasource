package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/redshift-datasource/pkg/redshift/routes"
	"github.com/grafana/sqlds/v2"
)

func main() {
	// Start listening to requests sent from Grafana.
	s := redshift.New()
	ds := sqlds.NewAsyncDatasource(s)
	ds.Completable = s
	ds.CustomRoutes = routes.New(s).Routes()

	if err := datasource.Manage(
		"grafana-redshift-datasource",
		ds.NewDatasource,
		datasource.ManageOpts{},
	); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
