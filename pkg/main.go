package main

import (
	"context"
	"os"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/mcp"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/mcp/fromschema"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/redshift-datasource/pkg/redshift/routes"
)

const dsID = "grafana-redshift-datasource"

func main() {
	schema, err := loadSchema()
	if err != nil {
		log.DefaultLogger.Error("schema load failed", "err", err)
		os.Exit(1)
	}

	mcpServer := mcp.NewServer(mcp.ServerOpts{
		Name:    dsID,
		Version: "2.5.0",
		Addr:    ":7402",
	})

	factory := MakeDatasourceFactory()
	im := datasource.NewInstanceManager(factory)
	mgr := datasource.NewAutomanagementHandler(im)

	mcpServer.BindQueryDataHandler(mgr)
	mcpServer.BindCallResourceHandler(mgr)
	mcpServer.BindCheckHealthHandler(mgr)

	fromschema.RegisterQueryTools(mcpServer, schema)
	fromschema.RegisterRouteTools(mcpServer, schema)
	fromschema.RegisterQueryExamples(mcpServer, schema)
	fromschema.RegisterHealthCheckTool(mcpServer)

	if err := datasource.Manage(dsID, factory, datasource.ManageOpts{
		MCPServer: mcpServer,
	}); err != nil {
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
