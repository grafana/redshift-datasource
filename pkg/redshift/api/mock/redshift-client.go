package mock

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshift/redshiftiface"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type MockRedshiftClient struct {
	ExecutionResult         *redshiftdataapiservice.ExecuteStatementOutput
	DescribeStatementOutput *redshiftdataapiservice.DescribeStatementOutput
	ListStatementsOutput    *redshiftdataapiservice.ListStatementsOutput
	// Schemas > Tables > Columns
	Resources map[string]map[string][]string
	Secrets   []string
	Secret    string
	Clusters  []string

	secretsmanageriface.SecretsManagerAPI
	redshiftdataapiservice.RedshiftDataAPIService
	redshiftiface.RedshiftAPI
}

type MockRedshiftClientError struct {
	redshiftiface.RedshiftAPI
}

type MockRedshiftClientNil struct {
	redshiftiface.RedshiftAPI
}

func (m *MockRedshiftClient) ExecuteStatementWithContext(ctx aws.Context, input *redshiftdataapiservice.ExecuteStatementInput, opts ...request.Option) (*redshiftdataapiservice.ExecuteStatementOutput, error) {
	return m.ExecutionResult, nil
}

func (m *MockRedshiftClient) DescribeStatementWithContext(_ aws.Context, input *redshiftdataapiservice.DescribeStatementInput, _ ...request.Option) (*redshiftdataapiservice.DescribeStatementOutput, error) {
	return m.DescribeStatementOutput, nil
}

func (m *MockRedshiftClient) ListStatementsWithContext(_ aws.Context, input *redshiftdataapiservice.ListStatementsInput, _ ...request.Option) (*redshiftdataapiservice.ListStatementsOutput, error) {
	return m.ListStatementsOutput, nil
}

func (m *MockRedshiftClient) ListSchemasWithContext(ctx aws.Context, input *redshiftdataapiservice.ListSchemasInput, opts ...request.Option) (*redshiftdataapiservice.ListSchemasOutput, error) {
	res := &redshiftdataapiservice.ListSchemasOutput{}
	for sc := range m.Resources {
		res.Schemas = append(res.Schemas, aws.String(sc))
	}
	return res, nil
}

func (m *MockRedshiftClient) ListTablesWithContext(ctx aws.Context, input *redshiftdataapiservice.ListTablesInput, opts ...request.Option) (*redshiftdataapiservice.ListTablesOutput, error) {
	res := &redshiftdataapiservice.ListTablesOutput{}
	for t := range m.Resources[*input.SchemaPattern] {
		res.Tables = append(res.Tables, &redshiftdataapiservice.TableMember{Name: aws.String(t)})
	}
	return res, nil
}

func (m *MockRedshiftClient) DescribeTableWithContext(ctx aws.Context, input *redshiftdataapiservice.DescribeTableInput, opts ...request.Option) (*redshiftdataapiservice.DescribeTableOutput, error) {
	res := &redshiftdataapiservice.DescribeTableOutput{}
	tables := m.Resources[*input.Schema]
	for _, c := range tables[*input.Table] {
		res.ColumnList = append(res.ColumnList, &redshiftdataapiservice.ColumnMetadata{Name: aws.String(c)})
	}
	return res, nil
}

func (m *MockRedshiftClient) ListSecretsWithContext(ctx aws.Context, input *secretsmanager.ListSecretsInput, opts ...request.Option) (*secretsmanager.ListSecretsOutput, error) {
	r := &secretsmanager.ListSecretsOutput{}
	for _, c := range m.Secrets {
		r.SecretList = append(r.SecretList, &secretsmanager.SecretListEntry{ARN: aws.String(fmt.Sprintf("arn:%s", c)), Name: aws.String(c)})
	}
	return r, nil
}

func (m *MockRedshiftClient) GetSecretValueWithContext(ctx aws.Context, input *secretsmanager.GetSecretValueInput, opts ...request.Option) (*secretsmanager.GetSecretValueOutput, error) {
	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(m.Secret),
	}, nil
}

func (m *MockRedshiftClient) DescribeClusters(input *redshift.DescribeClustersInput) (*redshift.DescribeClustersOutput, error) {
	r := []*redshift.Cluster{}
	for _, c := range m.Clusters {
		r = append(r, &redshift.Cluster{
			ClusterIdentifier: aws.String(c),
			Endpoint: &redshift.Endpoint{
				Address: aws.String(c),
				Port:    aws.Int64(123),
			},
			DBName: aws.String(c),
		})
	}
	res := redshift.DescribeClustersOutput{
		Clusters: r,
	}
	return &res, nil
}

func (m *MockRedshiftClientError) DescribeClusters(input *redshift.DescribeClustersInput) (*redshift.DescribeClustersOutput, error) {
	return nil, fmt.Errorf("Boom")
}
func (m *MockRedshiftClientNil) DescribeClusters(input *redshift.DescribeClustersInput) (*redshift.DescribeClustersOutput, error) {
	return nil, nil
}
