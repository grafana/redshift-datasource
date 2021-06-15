package driver

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

const singlePageResponseQueryId = "singlePageResponse"
const multiPageResponseQueryId = "multiPageResponse"

var columnMetaData []*redshiftdataapiservice.ColumnMetadata = []*redshiftdataapiservice.ColumnMetadata{
	{
		Name:     aws.String("col1"),
		Nullable: aws.Int64(1),
		TypeName: aws.String("varchar"),
	},
	{
		Name:     aws.String("col2"),
		Nullable: aws.Int64(1),
		TypeName: aws.String("varchar"),
	},
}

var twoRecords = [][]*redshiftdataapiservice.Field{
	{
		&redshiftdataapiservice.Field{
			StringValue: aws.String("row1col1"),
		},
		&redshiftdataapiservice.Field{
			StringValue: aws.String("row1col2"),
		},
	},
	{
		&redshiftdataapiservice.Field{
			StringValue: aws.String("row2col1"),
		},
		&redshiftdataapiservice.Field{
			StringValue: aws.String("row2col2"),
		},
	},
}

type mockRedshiftService struct {
	getStatementResult   *redshiftdataapiservice.GetStatementResultOutput
	calledTimesCounter   int
	calledTimesCountDown int
}

func newMockRedshiftService() *mockRedshiftService {
	return &mockRedshiftService{calledTimesCounter: 0}
}

func (s *mockRedshiftService) GetStatementResult(input *redshiftdataapiservice.GetStatementResultInput) (*redshiftdataapiservice.GetStatementResultOutput, error) {
	s.calledTimesCounter++

	if *input.Id == singlePageResponseQueryId || s.calledTimesCounter == 2 {
		return &redshiftdataapiservice.GetStatementResultOutput{
			ColumnMetadata: columnMetaData,
			Records:        twoRecords,
		}, nil
	}

	if *input.Id == multiPageResponseQueryId && s.calledTimesCounter == 1 {
		return &redshiftdataapiservice.GetStatementResultOutput{
			ColumnMetadata: columnMetaData,
			Records:        twoRecords,
			NextToken:      aws.String("nexttoken"),
		}, nil
	}

	return nil, fmt.Errorf("no test response for this query id")
}

const DESCRIBE_STATEMENT_FAILED = "DESCRIBE_STATEMENT_FAILED"
const DESCRIBE_STATEMENT_SUCCEEDED = "DESCRIBE_STATEMENT_FINISHED"

func (s *mockRedshiftService) DescribeStatementWithContext(_ aws.Context, input *redshiftdataapiservice.DescribeStatementInput, _ ...request.Option) (*redshiftdataapiservice.DescribeStatementOutput, error) {
	s.calledTimesCountDown--
	s.calledTimesCounter++

	if *input.Id == DESCRIBE_STATEMENT_FAILED {
		return &redshiftdataapiservice.DescribeStatementOutput{
			Status: aws.String(redshiftdataapiservice.StatusStringFailed),
			Error:  aws.String(DESCRIBE_STATEMENT_FAILED),
		}, nil
	}

	if s.calledTimesCountDown == 0 {
		return &redshiftdataapiservice.DescribeStatementOutput{
			Status: aws.String(redshiftdataapiservice.StatusStringFinished),
		}, nil
	} else {
		return &redshiftdataapiservice.DescribeStatementOutput{
			Status: aws.String(redshiftdataapiservice.StatusStringStarted),
		}, nil
	}
}

func (s *mockRedshiftService) CancelStatement(*redshiftdataapiservice.CancelStatementInput) (*redshiftdataapiservice.CancelStatementOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) CancelStatementWithContext(aws.Context, *redshiftdataapiservice.CancelStatementInput, ...request.Option) (*redshiftdataapiservice.CancelStatementOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) CancelStatementRequest(*redshiftdataapiservice.CancelStatementInput) (*request.Request, *redshiftdataapiservice.CancelStatementOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) DescribeStatement(*redshiftdataapiservice.DescribeStatementInput) (*redshiftdataapiservice.DescribeStatementOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) DescribeStatementRequest(*redshiftdataapiservice.DescribeStatementInput) (*request.Request, *redshiftdataapiservice.DescribeStatementOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) DescribeTable(*redshiftdataapiservice.DescribeTableInput) (*redshiftdataapiservice.DescribeTableOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) DescribeTableWithContext(aws.Context, *redshiftdataapiservice.DescribeTableInput, ...request.Option) (*redshiftdataapiservice.DescribeTableOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) DescribeTableRequest(*redshiftdataapiservice.DescribeTableInput) (*request.Request, *redshiftdataapiservice.DescribeTableOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) DescribeTablePages(*redshiftdataapiservice.DescribeTableInput, func(*redshiftdataapiservice.DescribeTableOutput, bool) bool) error {
	panic("not implemented")
}

func (s *mockRedshiftService) DescribeTablePagesWithContext(aws.Context, *redshiftdataapiservice.DescribeTableInput, func(*redshiftdataapiservice.DescribeTableOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ExecuteStatement(*redshiftdataapiservice.ExecuteStatementInput) (*redshiftdataapiservice.ExecuteStatementOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ExecuteStatementWithContext(aws.Context, *redshiftdataapiservice.ExecuteStatementInput, ...request.Option) (*redshiftdataapiservice.ExecuteStatementOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ExecuteStatementRequest(*redshiftdataapiservice.ExecuteStatementInput) (*request.Request, *redshiftdataapiservice.ExecuteStatementOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) GetStatementResultWithContext(aws.Context, *redshiftdataapiservice.GetStatementResultInput, ...request.Option) (*redshiftdataapiservice.GetStatementResultOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) GetStatementResultRequest(*redshiftdataapiservice.GetStatementResultInput) (*request.Request, *redshiftdataapiservice.GetStatementResultOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) GetStatementResultPages(*redshiftdataapiservice.GetStatementResultInput, func(*redshiftdataapiservice.GetStatementResultOutput, bool) bool) error {
	panic("not implemented")
}

func (s *mockRedshiftService) GetStatementResultPagesWithContext(aws.Context, *redshiftdataapiservice.GetStatementResultInput, func(*redshiftdataapiservice.GetStatementResultOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListDatabases(*redshiftdataapiservice.ListDatabasesInput) (*redshiftdataapiservice.ListDatabasesOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListDatabasesWithContext(aws.Context, *redshiftdataapiservice.ListDatabasesInput, ...request.Option) (*redshiftdataapiservice.ListDatabasesOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListDatabasesPagesWithContext(aws.Context, *redshiftdataapiservice.ListDatabasesInput, func(*redshiftdataapiservice.ListDatabasesOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListDatabasesRequest(*redshiftdataapiservice.ListDatabasesInput) (*request.Request, *redshiftdataapiservice.ListDatabasesOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListDatabasesPages(*redshiftdataapiservice.ListDatabasesInput, func(*redshiftdataapiservice.ListDatabasesOutput, bool) bool) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListSchemas(*redshiftdataapiservice.ListSchemasInput) (*redshiftdataapiservice.ListSchemasOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListSchemasWithContext(aws.Context, *redshiftdataapiservice.ListSchemasInput, ...request.Option) (*redshiftdataapiservice.ListSchemasOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListSchemasRequest(*redshiftdataapiservice.ListSchemasInput) (*request.Request, *redshiftdataapiservice.ListSchemasOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListSchemasPages(*redshiftdataapiservice.ListSchemasInput, func(*redshiftdataapiservice.ListSchemasOutput, bool) bool) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListSchemasPagesWithContext(aws.Context, *redshiftdataapiservice.ListSchemasInput, func(*redshiftdataapiservice.ListSchemasOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListStatements(*redshiftdataapiservice.ListStatementsInput) (*redshiftdataapiservice.ListStatementsOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListStatementsWithContext(aws.Context, *redshiftdataapiservice.ListStatementsInput, ...request.Option) (*redshiftdataapiservice.ListStatementsOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListStatementsRequest(*redshiftdataapiservice.ListStatementsInput) (*request.Request, *redshiftdataapiservice.ListStatementsOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListStatementsPages(*redshiftdataapiservice.ListStatementsInput, func(*redshiftdataapiservice.ListStatementsOutput, bool) bool) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListStatementsPagesWithContext(aws.Context, *redshiftdataapiservice.ListStatementsInput, func(*redshiftdataapiservice.ListStatementsOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListTables(*redshiftdataapiservice.ListTablesInput) (*redshiftdataapiservice.ListTablesOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListTablesWithContext(aws.Context, *redshiftdataapiservice.ListTablesInput, ...request.Option) (*redshiftdataapiservice.ListTablesOutput, error) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListTablesRequest(*redshiftdataapiservice.ListTablesInput) (*request.Request, *redshiftdataapiservice.ListTablesOutput) {
	panic("not implemented")
}

func (s *mockRedshiftService) ListTablesPages(*redshiftdataapiservice.ListTablesInput, func(*redshiftdataapiservice.ListTablesOutput, bool) bool) error {
	panic("not implemented")
}

func (s *mockRedshiftService) ListTablesPagesWithContext(aws.Context, *redshiftdataapiservice.ListTablesInput, func(*redshiftdataapiservice.ListTablesOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}
