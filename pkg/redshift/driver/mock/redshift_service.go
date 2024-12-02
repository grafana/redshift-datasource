package redshiftservicemock

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	redshiftdataV2 "github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	redshiftdataV2types "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
)

const SinglePageResponseQueryId = "singlePageResponse"
const MultiPageResponseQueryId = "multiPageResponse"

var columnMetaData = []redshiftdataV2types.ColumnMetadata{
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

var twoRecords = [][]redshiftdataV2types.Field{
	{
		&redshiftdataV2types.FieldMemberStringValue{Value: "row1col1"},
		&redshiftdataV2types.FieldMemberStringValue{Value: "row1col2"},
	},
	{
		&redshiftdataV2types.FieldMemberStringValue{Value: "row2col1"},
		&redshiftdataV2types.FieldMemberStringValue{Value: "row2col2"},
	},
}

type RedshiftService struct {
	CalledTimesCounter   int
	CalledTimesCountDown int
}

// GetStatementResult returns a GetStatementResultOutput
// When mockRedshiftService.calledTimesCountDown is more than 0, the GetStatementResultOutput will have a next token
func (s *RedshiftService) GetStatementResult(_ context.Context, _ *redshiftdataV2.GetStatementResultInput, _ ...func(*redshiftdataV2.Options)) (*redshiftdataV2.GetStatementResultOutput, error) {
	s.CalledTimesCountDown--
	s.CalledTimesCounter++

	if s.CalledTimesCountDown == 0 {
		return &redshiftdataV2.GetStatementResultOutput{
			ColumnMetadata: columnMetaData,
			Records:        twoRecords,
		}, nil
	}

	return &redshiftdataV2.GetStatementResultOutput{
		ColumnMetadata: columnMetaData,
		Records:        twoRecords,
		NextToken:      aws.String("nexttoken"),
	}, nil
}

const DESCRIBE_STATEMENT_FAILED = "DESCRIBE_STATEMENT_FAILED"

// DescribeStatement returns a DescribeStatementOutput
// When DescribeStatementInput.Id == DESCRIBE_STATEMENT_FAILED, an the output will include an error message that is equal to the input id
// When DescribeStatementInput.Id == DESCRIBE_STATEMENT_FINISHED, the output will have a status redshiftdataV2.StatusStringFinished once mockRedshiftService.calledTimesCountDown == 0
func (s *RedshiftService) DescribeStatement(_ context.Context, input *redshiftdataV2.DescribeStatementInput, _ ...func(options *redshiftdataV2.Options)) (*redshiftdataV2.DescribeStatementOutput, error) {
	s.CalledTimesCountDown--
	s.CalledTimesCounter++

	output := &redshiftdataV2.DescribeStatementOutput{}
	if s.CalledTimesCountDown == 0 {
		if *input.Id == DESCRIBE_STATEMENT_FAILED {
			output.Status = redshiftdataV2types.StatusStringFailed
			output.Error = aws.String(DESCRIBE_STATEMENT_FAILED)
		} else {
			output.Status = redshiftdataV2types.StatusStringFinished
		}
	} else {
		output.Status = redshiftdataV2types.StatusStringStarted
	}
	return output, nil
}
