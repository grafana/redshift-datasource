module github.com/grafana/redshift-datasource

go 1.16

// pointing to a local package until this PR is merged: https://github.com/grafana/sqlds/pull/17
replace github.com/grafana/sqlds => ../../sqlds

require (
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/aws/aws-sdk-go v1.38.54
	github.com/grafana/grafana-aws-sdk v0.6.0
	github.com/grafana/grafana-plugin-sdk-go v0.104.0
	github.com/grafana/sqlds v1.0.11 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
