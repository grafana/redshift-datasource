package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/sqlds"
)

func main() {
	// Start listening to requests sent from Grafana.
	s := &redshift.RedshiftDatasource{}
	ds := sqlds.NewDatasource(s)
	ds.Completable = s

	if err := datasource.Manage(
		"grafana-redshift-datasource",
		ds.NewDatasource,
		datasource.ManageOpts{},
	); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
