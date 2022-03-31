package redshiftservicemock

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
)

const SinglePageResponseQueryId = "singlePageResponse"
const MultiPageResponseQueryId = "multiPageResponse"

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

type RedshiftService struct {
	CalledTimesCounter   int
	CalledTimesCountDown int
}

func NewMockRedshiftService() *RedshiftService {
	return &RedshiftService{CalledTimesCounter: 0}
}

// GetStatementResult returns a GetStatementResultOutput
// When mockRedshiftService.calledTimesCountDown is more than 0, the GetStatementResultOutput will have a next token
func (s *RedshiftService) GetStatementResult(input *redshiftdataapiservice.GetStatementResultInput) (*redshiftdataapiservice.GetStatementResultOutput, error) {
	s.CalledTimesCountDown--
	s.CalledTimesCounter++

	if s.CalledTimesCountDown == 0 {
		return &redshiftdataapiservice.GetStatementResultOutput{
			ColumnMetadata: columnMetaData,
			Records:        twoRecords,
		}, nil
	}

	return &redshiftdataapiservice.GetStatementResultOutput{
		ColumnMetadata: columnMetaData,
		Records:        twoRecords,
		NextToken:      aws.String("nexttoken"),
	}, nil
}

const DESCRIBE_STATEMENT_FAILED = "DESCRIBE_STATEMENT_FAILED"
const DESCRIBE_STATEMENT_SUCCEEDED = "DESCRIBE_STATEMENT_FINISHED"

// DescribeStatementWithContext returns a DescribeStatementOutput
// When DescribeStatementInput.Id == DESCRIBE_STATEMENT_FAILED, an the output will include an error message that is equal to the input id
// When DescribeStatementInput.Id == DESCRIBE_STATEMENT_FINISHED, the output will have a status redshiftdataapiservice.StatusStringFinished once mockRedshiftService.calledTimesCountDown == 0
func (s *RedshiftService) DescribeStatementWithContext(_ aws.Context, input *redshiftdataapiservice.DescribeStatementInput, _ ...request.Option) (*redshiftdataapiservice.DescribeStatementOutput, error) {
	s.CalledTimesCountDown--
	s.CalledTimesCounter++

	output := &redshiftdataapiservice.DescribeStatementOutput{}
	if s.CalledTimesCountDown == 0 {
		if *input.Id == DESCRIBE_STATEMENT_FAILED {
			output.Status = aws.String(redshiftdataapiservice.StatusStringFailed)
			output.Error = aws.String(DESCRIBE_STATEMENT_FAILED)
		} else {
			output.Status = aws.String(redshiftdataapiservice.StatusStringFinished)
		}
	} else {
		output.Status = aws.String(redshiftdataapiservice.StatusStringStarted)
	}
	return output, nil
}

func (s *RedshiftService) CancelStatement(*redshiftdataapiservice.CancelStatementInput) (*redshiftdataapiservice.CancelStatementOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) CancelStatementWithContext(aws.Context, *redshiftdataapiservice.CancelStatementInput, ...request.Option) (*redshiftdataapiservice.CancelStatementOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) CancelStatementRequest(*redshiftdataapiservice.CancelStatementInput) (*request.Request, *redshiftdataapiservice.CancelStatementOutput) {
	panic("not implemented")
}

func (s *RedshiftService) DescribeStatement(*redshiftdataapiservice.DescribeStatementInput) (*redshiftdataapiservice.DescribeStatementOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) DescribeStatementRequest(*redshiftdataapiservice.DescribeStatementInput) (*request.Request, *redshiftdataapiservice.DescribeStatementOutput) {
	panic("not implemented")
}

func (s *RedshiftService) DescribeTable(*redshiftdataapiservice.DescribeTableInput) (*redshiftdataapiservice.DescribeTableOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) DescribeTableWithContext(aws.Context, *redshiftdataapiservice.DescribeTableInput, ...request.Option) (*redshiftdataapiservice.DescribeTableOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) DescribeTableRequest(*redshiftdataapiservice.DescribeTableInput) (*request.Request, *redshiftdataapiservice.DescribeTableOutput) {
	panic("not implemented")
}

func (s *RedshiftService) DescribeTablePages(*redshiftdataapiservice.DescribeTableInput, func(*redshiftdataapiservice.DescribeTableOutput, bool) bool) error {
	panic("not implemented")
}

func (s *RedshiftService) DescribeTablePagesWithContext(aws.Context, *redshiftdataapiservice.DescribeTableInput, func(*redshiftdataapiservice.DescribeTableOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *RedshiftService) ExecuteStatement(*redshiftdataapiservice.ExecuteStatementInput) (*redshiftdataapiservice.ExecuteStatementOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ExecuteStatementWithContext(aws.Context, *redshiftdataapiservice.ExecuteStatementInput, ...request.Option) (*redshiftdataapiservice.ExecuteStatementOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ExecuteStatementRequest(*redshiftdataapiservice.ExecuteStatementInput) (*request.Request, *redshiftdataapiservice.ExecuteStatementOutput) {
	panic("not implemented")
}

func (s *RedshiftService) GetStatementResultWithContext(aws.Context, *redshiftdataapiservice.GetStatementResultInput, ...request.Option) (*redshiftdataapiservice.GetStatementResultOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) GetStatementResultRequest(*redshiftdataapiservice.GetStatementResultInput) (*request.Request, *redshiftdataapiservice.GetStatementResultOutput) {
	panic("not implemented")
}

func (s *RedshiftService) GetStatementResultPages(*redshiftdataapiservice.GetStatementResultInput, func(*redshiftdataapiservice.GetStatementResultOutput, bool) bool) error {
	panic("not implemented")
}

func (s *RedshiftService) GetStatementResultPagesWithContext(aws.Context, *redshiftdataapiservice.GetStatementResultInput, func(*redshiftdataapiservice.GetStatementResultOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *RedshiftService) ListDatabases(*redshiftdataapiservice.ListDatabasesInput) (*redshiftdataapiservice.ListDatabasesOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListDatabasesWithContext(aws.Context, *redshiftdataapiservice.ListDatabasesInput, ...request.Option) (*redshiftdataapiservice.ListDatabasesOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListDatabasesPagesWithContext(aws.Context, *redshiftdataapiservice.ListDatabasesInput, func(*redshiftdataapiservice.ListDatabasesOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *RedshiftService) ListDatabasesRequest(*redshiftdataapiservice.ListDatabasesInput) (*request.Request, *redshiftdataapiservice.ListDatabasesOutput) {
	panic("not implemented")
}

func (s *RedshiftService) ListDatabasesPages(*redshiftdataapiservice.ListDatabasesInput, func(*redshiftdataapiservice.ListDatabasesOutput, bool) bool) error {
	panic("not implemented")
}

func (s *RedshiftService) ListSchemas(*redshiftdataapiservice.ListSchemasInput) (*redshiftdataapiservice.ListSchemasOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListSchemasWithContext(aws.Context, *redshiftdataapiservice.ListSchemasInput, ...request.Option) (*redshiftdataapiservice.ListSchemasOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListSchemasRequest(*redshiftdataapiservice.ListSchemasInput) (*request.Request, *redshiftdataapiservice.ListSchemasOutput) {
	panic("not implemented")
}

func (s *RedshiftService) ListSchemasPages(*redshiftdataapiservice.ListSchemasInput, func(*redshiftdataapiservice.ListSchemasOutput, bool) bool) error {
	panic("not implemented")
}

func (s *RedshiftService) ListSchemasPagesWithContext(aws.Context, *redshiftdataapiservice.ListSchemasInput, func(*redshiftdataapiservice.ListSchemasOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *RedshiftService) ListStatements(*redshiftdataapiservice.ListStatementsInput) (*redshiftdataapiservice.ListStatementsOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListStatementsWithContext(aws.Context, *redshiftdataapiservice.ListStatementsInput, ...request.Option) (*redshiftdataapiservice.ListStatementsOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListStatementsRequest(*redshiftdataapiservice.ListStatementsInput) (*request.Request, *redshiftdataapiservice.ListStatementsOutput) {
	panic("not implemented")
}

func (s *RedshiftService) ListStatementsPages(*redshiftdataapiservice.ListStatementsInput, func(*redshiftdataapiservice.ListStatementsOutput, bool) bool) error {
	panic("not implemented")
}

func (s *RedshiftService) ListStatementsPagesWithContext(aws.Context, *redshiftdataapiservice.ListStatementsInput, func(*redshiftdataapiservice.ListStatementsOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *RedshiftService) ListTables(*redshiftdataapiservice.ListTablesInput) (*redshiftdataapiservice.ListTablesOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListTablesWithContext(aws.Context, *redshiftdataapiservice.ListTablesInput, ...request.Option) (*redshiftdataapiservice.ListTablesOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) ListTablesRequest(*redshiftdataapiservice.ListTablesInput) (*request.Request, *redshiftdataapiservice.ListTablesOutput) {
	panic("not implemented")
}

func (s *RedshiftService) ListTablesPages(*redshiftdataapiservice.ListTablesInput, func(*redshiftdataapiservice.ListTablesOutput, bool) bool) error {
	panic("not implemented")
}

func (s *RedshiftService) ListTablesPagesWithContext(aws.Context, *redshiftdataapiservice.ListTablesInput, func(*redshiftdataapiservice.ListTablesOutput, bool) bool, ...request.Option) error {
	panic("not implemented")
}

func (s *RedshiftService) BatchExecuteStatement(input *redshiftdataapiservice.BatchExecuteStatementInput) (*redshiftdataapiservice.BatchExecuteStatementOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) BatchExecuteStatementWithContext(context aws.Context, input *redshiftdataapiservice.BatchExecuteStatementInput, option ...request.Option) (*redshiftdataapiservice.BatchExecuteStatementOutput, error) {
	panic("not implemented")
}

func (s *RedshiftService) BatchExecuteStatementRequest(input *redshiftdataapiservice.BatchExecuteStatementInput) (*request.Request, *redshiftdataapiservice.BatchExecuteStatementOutput) {
	panic("not implemented")
}
