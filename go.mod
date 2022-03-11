module github.com/grafana/redshift-datasource

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.62
	github.com/google/go-cmp v0.5.7
	github.com/grafana/grafana-aws-sdk v0.10.1
	github.com/grafana/grafana-plugin-sdk-go v0.125.0
	github.com/grafana/sqlds/v2 v2.3.3
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/grafana/grafana-plugin-sdk-go => ../grafana-plugin-sdk-go

replace github.com/grafana/sqlds/v2 => ../sqlds

replace github.com/grafana/grafana-aws-sdk => ../grafana-aws-sdk
