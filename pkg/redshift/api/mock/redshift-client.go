package mock

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	redshifttypes "github.com/aws/aws-sdk-go-v2/service/redshift/types"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	redshiftdatatypes "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
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
	ExecutionResult         *redshiftdata.ExecuteStatementOutput
	DescribeStatementOutput *redshiftdata.DescribeStatementOutput
	ListStatementsOutput    *redshiftdata.ListStatementsOutput
	// Schemas > Tables > Columns
	Resources map[string]map[string][]string
	Clusters  []string

	redshiftdata.ListDatabasesAPIClient
	redshiftdata.GetStatementResultAPIClient
	types.CancelStatementAPIClient
	redshiftdata.DescribeTableAPIClient
}

type MockRedshiftClientError struct {
}

type MockRedshiftClientNil struct {
}

func (mc *MockRedshiftClient) ExecuteStatement(_ context.Context, _ *redshiftdata.ExecuteStatementInput, _ ...func(*redshiftdata.Options)) (*redshiftdata.ExecuteStatementOutput, error) {
	return mc.ExecutionResult, nil
}

func (mc *MockRedshiftClient) DescribeStatement(_ context.Context, _ *redshiftdata.DescribeStatementInput, _ ...func(*redshiftdata.Options)) (*redshiftdata.DescribeStatementOutput, error) {
	return mc.DescribeStatementOutput, nil
}

func (mc *MockRedshiftClient) ListStatements(_ context.Context, _ *redshiftdata.ListStatementsInput, _ ...func(*redshiftdata.Options)) (*redshiftdata.ListStatementsOutput, error) {
	return mc.ListStatementsOutput, nil
}

func (mc *MockRedshiftClient) ListSchemas(_ context.Context, _ *redshiftdata.ListSchemasInput, _ ...func(*redshiftdata.Options)) (*redshiftdata.ListSchemasOutput, error) {
	res := &redshiftdata.ListSchemasOutput{}
	for sc := range mc.Resources {
		res.Schemas = append(res.Schemas, sc)
	}
	return res, nil
}

func (mc *MockRedshiftClient) ListTables(_ context.Context, input *redshiftdata.ListTablesInput, _ ...func(*redshiftdata.Options)) (*redshiftdata.ListTablesOutput, error) {
	res := &redshiftdata.ListTablesOutput{}
	for t := range mc.Resources[*input.SchemaPattern] {
		res.Tables = append(res.Tables, redshiftdatatypes.TableMember{Name: aws.String(t)})
	}
	return res, nil
}

func (mc *MockRedshiftClient) DescribeTable(_ context.Context, input *redshiftdata.DescribeTableInput, _ ...func(*redshiftdata.Options)) (*redshiftdata.DescribeTableOutput, error) {
	res := &redshiftdata.DescribeTableOutput{}
	tables := mc.Resources[*input.Schema]
	for _, c := range tables[*input.Table] {
		res.ColumnList = append(res.ColumnList, redshiftdatatypes.ColumnMetadata{Name: aws.String(c)})
	}
	return res, nil
}

func (mc *MockRedshiftClient) DescribeClusters(_ context.Context, _ *redshift.DescribeClustersInput, _ ...func(*redshift.Options)) (*redshift.DescribeClustersOutput, error) {
	r := []redshifttypes.Cluster{}
	for _, c := range mc.Clusters {
		r = append(r, redshifttypes.Cluster{
			ClusterIdentifier: aws.String(c),
			Endpoint: &redshifttypes.Endpoint{
				Address: aws.String(c),
				Port:    aws.Int32(123),
			},
			DBName: aws.String(c),
		})
	}
	res := redshift.DescribeClustersOutput{
		Clusters: r,
	}
	return &res, nil
}

func (m *MockRedshiftClientError) DescribeClusters(_ context.Context, _ *redshift.DescribeClustersInput, _ ...func(*redshift.Options)) (*redshift.DescribeClustersOutput, error) {
	return nil, fmt.Errorf("Boom")
}
func (m *MockRedshiftClientNil) DescribeClusters(_ context.Context, _ *redshift.DescribeClustersInput, _ ...func(*redshift.Options)) (*redshift.DescribeClustersOutput, error) {
	return nil, nil
}
