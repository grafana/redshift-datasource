package mock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	redshiftserverlesstypes "github.com/aws/aws-sdk-go-v2/service/redshiftserverless/types"
)

type MockRedshiftServerlessClient struct {
	// Schemas > Tables > Columns
	Resources  map[string]map[string][]string
	Secrets    []string
	Secret     string
	Workgroups []string
}

type MockRedshiftServerlessClientError struct{}

type MockRedshiftServerlessClientNil struct{}

func (m *MockRedshiftServerlessClient) ListWorkgroups(_ context.Context, _ *redshiftserverless.ListWorkgroupsInput, _ ...func(*redshiftserverless.Options)) (*redshiftserverless.ListWorkgroupsOutput, error) {
	r := []redshiftserverlesstypes.Workgroup{}
	for _, c := range m.Workgroups {
		r = append(r, redshiftserverlesstypes.Workgroup{
			WorkgroupName: aws.String(c),
			Endpoint: &redshiftserverlesstypes.Endpoint{
				Address: aws.String(c),
				Port:    aws.Int32(123),
			},
		})
	}
	res := redshiftserverless.ListWorkgroupsOutput{
		Workgroups: r,
	}
	return &res, nil
}

func (m *MockRedshiftServerlessClientError) ListWorkgroups(_ context.Context, _ *redshiftserverless.ListWorkgroupsInput, _ ...func(*redshiftserverless.Options)) (*redshiftserverless.ListWorkgroupsOutput, error) {
	return nil, fmt.Errorf("Boom")
}
func (m *MockRedshiftServerlessClientNil) ListWorkgroups(_ context.Context, _ *redshiftserverless.ListWorkgroupsInput, _ ...func(*redshiftserverless.Options)) (*redshiftserverless.ListWorkgroupsOutput, error) {
	return nil, nil
}
