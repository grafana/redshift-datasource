package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/redshift-datasource/pkg/redshift"
)

func main() {
	// Start listening to requests sent from Grafana.
	// s := redshift.New()
	// ds := sqlds.NewDatasource(s)
	// ds.Completable = s
	// ds.CustomRoutes = routes.New(s).Routes()

	if err := datasource.Manage(
		"grafana-redshift-datasource",
		redshift.NewDatasource,
		datasource.ManageOpts{},
	); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
