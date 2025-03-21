package types

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type ExecuteStatementAPIClient interface {
	ExecuteStatement(context.Context, *redshiftdata.ExecuteStatementInput, ...func(*redshiftdata.Options)) (*redshiftdata.ExecuteStatementOutput, error)
}
type DescribeStatementAPIClient interface {
	DescribeStatement(context.Context, *redshiftdata.DescribeStatementInput, ...func(*redshiftdata.Options)) (*redshiftdata.DescribeStatementOutput, error)
}
type CancelStatementAPIClient interface {
	CancelStatement(context.Context, *redshiftdata.CancelStatementInput, ...func(*redshiftdata.Options)) (*redshiftdata.CancelStatementOutput, error)
}

type RedshiftDataClient interface {
	redshiftdata.DescribeTableAPIClient
	redshiftdata.ListDatabasesAPIClient
	redshiftdata.ListSchemasAPIClient
	redshiftdata.ListTablesAPIClient
	redshiftdata.DescribeTableAPIClient
	redshiftdata.GetStatementResultAPIClient

	ExecuteStatementAPIClient
	DescribeStatementAPIClient
	CancelStatementAPIClient
}

type RedshiftManagementClient interface {
	redshift.DescribeClustersAPIClient
}

type RedshiftSecretsClient interface {
	secretsmanager.ListSecretsAPIClient
	GetSecretValue(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type ServerlessAPIClient interface {
	redshiftserverless.ListWorkgroupsAPIClient
}
