module github.com/grafana/redshift-datasource

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.62
	github.com/google/go-cmp v0.5.6
	// TODO: Replace with final version when ready
	github.com/grafana/grafana-aws-sdk v0.7.1-0.20211213152905-b5c6292f2a5a
	github.com/grafana/grafana-plugin-sdk-go v0.114.0
	github.com/grafana/sqlds/v2 v2.3.3
	github.com/jpillora/backoff v1.0.0
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
)
