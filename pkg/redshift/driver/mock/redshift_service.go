package redshiftservicemock

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	redshiftdatatypes "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
)

const SinglePageResponseQueryId = "singlePageResponse"
const MultiPageResponseQueryId = "multiPageResponse"

var columnMetaData = []redshiftdatatypes.ColumnMetadata{
	{
		Name:     aws.String("col1"),
		Nullable: 1,
		TypeName: aws.String("varchar"),
	},
	{
		Name:     aws.String("col2"),
		Nullable: 1,
		TypeName: aws.String("varchar"),
	},
}

var twoRecords = [][]redshiftdatatypes.Field{
	{
		&redshiftdatatypes.FieldMemberStringValue{Value: "row1col1"},
		&redshiftdatatypes.FieldMemberStringValue{Value: "row1col2"},
	},
	{
		&redshiftdatatypes.FieldMemberStringValue{Value: "row2col1"},
		&redshiftdatatypes.FieldMemberStringValue{Value: "row2col2"},
	},
}

type RedshiftService struct {
	CalledTimesCounter   int
	CalledTimesCountDown int
}

// GetStatementResult returns a GetStatementResultOutput
// When mockRedshiftService.calledTimesCountDown is more than 0, the GetStatementResultOutput will have a next token
func (s *RedshiftService) GetStatementResult(_ context.Context, _ *redshiftdata.GetStatementResultInput, _ ...func(*redshiftdata.Options)) (*redshiftdata.GetStatementResultOutput, error) {
	s.CalledTimesCountDown--
	s.CalledTimesCounter++

	if s.CalledTimesCountDown == 0 {
		return &redshiftdata.GetStatementResultOutput{
			ColumnMetadata: columnMetaData,
			Records:        twoRecords,
		}, nil
	}

	return &redshiftdata.GetStatementResultOutput{
		ColumnMetadata: columnMetaData,
		Records:        twoRecords,
		NextToken:      aws.String("nexttoken"),
	}, nil
}

const DESCRIBE_STATEMENT_FAILED = "DESCRIBE_STATEMENT_FAILED"

// DescribeStatement returns a DescribeStatementOutput
// When DescribeStatementInput.Id == DESCRIBE_STATEMENT_FAILED, an the output will include an error message that is equal to the input id
// When DescribeStatementInput.Id == DESCRIBE_STATEMENT_FINISHED, the output will have a status redshiftdata.StatusStringFinished once mockRedshiftService.calledTimesCountDown == 0
func (s *RedshiftService) DescribeStatement(_ context.Context, input *redshiftdata.DescribeStatementInput, _ ...func(options *redshiftdata.Options)) (*redshiftdata.DescribeStatementOutput, error) {
	s.CalledTimesCountDown--
	s.CalledTimesCounter++

	output := &redshiftdata.DescribeStatementOutput{}
	if s.CalledTimesCountDown == 0 {
		if *input.Id == DESCRIBE_STATEMENT_FAILED {
			output.Status = redshiftdatatypes.StatusStringFailed
			output.Error = aws.String(DESCRIBE_STATEMENT_FAILED)
		} else {
			output.Status = redshiftdatatypes.StatusStringFinished
		}
	} else {
		output.Status = redshiftdatatypes.StatusStringStarted
	}
	return output, nil
}
