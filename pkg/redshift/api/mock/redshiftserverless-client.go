package mock

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	"github.com/aws/aws-sdk-go/service/redshiftserverless/redshiftserverlessiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type MockRedshiftServerlessClient struct {
	ExecutionResult         *redshiftdataapiservice.ExecuteStatementOutput
	DescribeStatementOutput *redshiftdataapiservice.DescribeStatementOutput
	ListStatementsOutput    *redshiftdataapiservice.ListStatementsOutput
	// Schemas > Tables > Columns
	Resources  map[string]map[string][]string
	Secrets    []string
	Secret     string
	Workgroups []string

	redshiftdataapiservice.RedshiftDataAPIService
	redshiftserverlessiface.RedshiftServerlessAPI
}

type MockRedshiftServerlessClientError struct {
	redshiftserverlessiface.RedshiftServerlessAPI
}

type MockRedshiftServerlessClientNil struct {
	redshiftserverlessiface.RedshiftServerlessAPI
}

func (m *MockRedshiftServerlessClient) ExecuteStatementWithContext(ctx aws.Context, input *redshiftdataapiservice.ExecuteStatementInput, opts ...request.Option) (*redshiftdataapiservice.ExecuteStatementOutput, error) {
	return m.ExecutionResult, nil
}

func (m *MockRedshiftServerlessClient) DescribeStatementWithContext(_ aws.Context, input *redshiftdataapiservice.DescribeStatementInput, _ ...request.Option) (*redshiftdataapiservice.DescribeStatementOutput, error) {
	return m.DescribeStatementOutput, nil
}

func (m *MockRedshiftServerlessClient) ListStatementsWithContext(_ aws.Context, input *redshiftdataapiservice.ListStatementsInput, _ ...request.Option) (*redshiftdataapiservice.ListStatementsOutput, error) {
	return m.ListStatementsOutput, nil
}

func (m *MockRedshiftServerlessClient) ListSchemasWithContext(ctx aws.Context, input *redshiftdataapiservice.ListSchemasInput, opts ...request.Option) (*redshiftdataapiservice.ListSchemasOutput, error) {
	res := &redshiftdataapiservice.ListSchemasOutput{}
	for sc := range m.Resources {
		res.Schemas = append(res.Schemas, aws.String(sc))
	}
	return res, nil
}

func (m *MockRedshiftServerlessClient) ListTablesWithContext(ctx aws.Context, input *redshiftdataapiservice.ListTablesInput, opts ...request.Option) (*redshiftdataapiservice.ListTablesOutput, error) {
	res := &redshiftdataapiservice.ListTablesOutput{}
	for t := range m.Resources[*input.SchemaPattern] {
		res.Tables = append(res.Tables, &redshiftdataapiservice.TableMember{Name: aws.String(t)})
	}
	return res, nil
}

func (m *MockRedshiftServerlessClient) DescribeTableWithContext(ctx aws.Context, input *redshiftdataapiservice.DescribeTableInput, opts ...request.Option) (*redshiftdataapiservice.DescribeTableOutput, error) {
	res := &redshiftdataapiservice.DescribeTableOutput{}
	tables := m.Resources[*input.Schema]
	for _, c := range tables[*input.Table] {
		res.ColumnList = append(res.ColumnList, &redshiftdataapiservice.ColumnMetadata{Name: aws.String(c)})
	}
	return res, nil
}

func (m *MockRedshiftServerlessClient) ListSecretsWithContext(ctx aws.Context, input *secretsmanager.ListSecretsInput, opts ...request.Option) (*secretsmanager.ListSecretsOutput, error) {
	r := &secretsmanager.ListSecretsOutput{}
	for _, c := range m.Secrets {
		r.SecretList = append(r.SecretList, &secretsmanager.SecretListEntry{ARN: aws.String(fmt.Sprintf("arn:%s", c)), Name: aws.String(c)})
	}
	return r, nil
}

func (m *MockRedshiftServerlessClient) GetSecretValueWithContext(ctx aws.Context, input *secretsmanager.GetSecretValueInput, opts ...request.Option) (*secretsmanager.GetSecretValueOutput, error) {
	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(m.Secret),
	}, nil
}

func (m *MockRedshiftServerlessClient) ListWorkgroups(input *redshiftserverless.ListWorkgroupsInput) (*redshiftserverless.ListWorkgroupsOutput, error) {
	r := []*redshiftserverless.Workgroup{}
	for _, c := range m.Workgroups {
		r = append(r, &redshiftserverless.Workgroup{
			WorkgroupName: aws.String(c),
			Endpoint: &redshiftserverless.Endpoint{
				Address: aws.String(c),
				Port:    aws.Int64(123),
			},
		})
	}
	res := redshiftserverless.ListWorkgroupsOutput{
		Workgroups: r,
	}
	return &res, nil
}

func (m *MockRedshiftServerlessClientError) ListWorkgroups(input *redshiftserverless.ListWorkgroupsInput) (*redshiftserverless.ListWorkgroupsOutput, error) {
	return nil, fmt.Errorf("Boom")
}

func (m *MockRedshiftServerlessClientNil) ListWorkgroups(input *redshiftserverless.ListWorkgroupsInput) (*redshiftserverless.ListWorkgroupsOutput, error) {
	return nil, nil
}
