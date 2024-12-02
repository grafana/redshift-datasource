package mock

import (
	"fmt"
	redshift2 "github.com/aws/aws-sdk-go-v2/service/redshift"
	redshiftV2types "github.com/aws/aws-sdk-go-v2/service/redshift/types"
	redshiftV2data "github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	redshiftdataV2types "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/grafana/redshift-datasource/pkg/redshift/api/types"
	"golang.org/x/net/context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

type MockRedshiftSecretsManager struct {
	Secret  string
	Secrets []string
}

func (msm *MockRedshiftSecretsManager) GetSecretValue(_ context.Context, _ *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(msm.Secret),
	}, nil
}
func (msm *MockRedshiftSecretsManager) ListSecrets(_ context.Context, _ *secretsmanager.ListSecretsInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.ListSecretsOutput, error) {
	r := &secretsmanager.ListSecretsOutput{}
	for _, c := range msm.Secrets {
		r.SecretList = append(r.SecretList, secretsmanagertypes.SecretListEntry{ARN: aws.String(fmt.Sprintf("arn:%s", c)), Name: aws.String(c)})
	}
	return r, nil
}

type MockRedshiftClient struct {
	ExecutionResult         *redshiftV2data.ExecuteStatementOutput
	DescribeStatementOutput *redshiftV2data.DescribeStatementOutput
	ListStatementsOutput    *redshiftV2data.ListStatementsOutput
	// Schemas > Tables > Columns
	Resources map[string]map[string][]string
	Clusters  []string

	redshiftV2data.ListDatabasesAPIClient
	redshiftV2data.GetStatementResultAPIClient
	types.CancelStatementAPIClient
	redshiftV2data.DescribeTableAPIClient
}

type MockRedshiftClientError struct {
}

type MockRedshiftClientNil struct {
}

func (mc *MockRedshiftClient) ExecuteStatement(_ context.Context, _ *redshiftV2data.ExecuteStatementInput, _ ...func(*redshiftV2data.Options)) (*redshiftV2data.ExecuteStatementOutput, error) {
	return mc.ExecutionResult, nil
}

func (mc *MockRedshiftClient) DescribeStatement(_ context.Context, _ *redshiftV2data.DescribeStatementInput, _ ...func(*redshiftV2data.Options)) (*redshiftV2data.DescribeStatementOutput, error) {
	return mc.DescribeStatementOutput, nil
}

func (mc *MockRedshiftClient) ListStatements(_ context.Context, _ *redshiftV2data.ListStatementsInput, _ ...func(*redshiftV2data.Options)) (*redshiftV2data.ListStatementsOutput, error) {
	return mc.ListStatementsOutput, nil
}

func (mc *MockRedshiftClient) ListSchemas(_ context.Context, _ *redshiftV2data.ListSchemasInput, _ ...func(*redshiftV2data.Options)) (*redshiftV2data.ListSchemasOutput, error) {
	res := &redshiftV2data.ListSchemasOutput{}
	for sc := range mc.Resources {
		res.Schemas = append(res.Schemas, sc)
	}
	return res, nil
}

func (mc *MockRedshiftClient) ListTables(_ context.Context, input *redshiftV2data.ListTablesInput, _ ...func(*redshiftV2data.Options)) (*redshiftV2data.ListTablesOutput, error) {
	res := &redshiftV2data.ListTablesOutput{}
	for t := range mc.Resources[*input.SchemaPattern] {
		res.Tables = append(res.Tables, redshiftdataV2types.TableMember{Name: aws.String(t)})
	}
	return res, nil
}

func (mc *MockRedshiftClient) DescribeTable(_ context.Context, input *redshiftV2data.DescribeTableInput, _ ...func(*redshiftV2data.Options)) (*redshiftV2data.DescribeTableOutput, error) {
	res := &redshiftV2data.DescribeTableOutput{}
	tables := mc.Resources[*input.Schema]
	for _, c := range tables[*input.Table] {
		res.ColumnList = append(res.ColumnList, redshiftdataV2types.ColumnMetadata{Name: aws.String(c)})
	}
	return res, nil
}

func (mc *MockRedshiftClient) DescribeClusters(_ context.Context, _ *redshift2.DescribeClustersInput, _ ...func(*redshift2.Options)) (*redshift2.DescribeClustersOutput, error) {
	r := []redshiftV2types.Cluster{}
	for _, c := range mc.Clusters {
		r = append(r, redshiftV2types.Cluster{
			ClusterIdentifier: aws.String(c),
			Endpoint: &redshiftV2types.Endpoint{
				Address: aws.String(c),
				Port:    aws.Int32(123),
			},
			DBName: aws.String(c),
		})
	}
	res := redshift2.DescribeClustersOutput{
		Clusters: r,
	}
	return &res, nil
}

func (m *MockRedshiftClientError) DescribeClusters(_ context.Context, _ *redshift2.DescribeClustersInput, _ ...func(*redshift2.Options)) (*redshift2.DescribeClustersOutput, error) {
	return nil, fmt.Errorf("Boom")
}
func (m *MockRedshiftClientNil) DescribeClusters(_ context.Context, _ *redshift2.DescribeClustersInput, _ ...func(*redshift2.Options)) (*redshift2.DescribeClustersOutput, error) {
	return nil, nil
}
