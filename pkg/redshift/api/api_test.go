package api

import (
	"context"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/api/mock"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v4"
)

func Test_apiInput(t *testing.T) {
	tests := []struct {
		description string
		settings    *models.RedshiftDataSourceSettings
		expected    apiInput
	}{
		{
			"serverless using temporary creds",
			&models.RedshiftDataSourceSettings{
				UseServerless:    true,
				UseManagedSecret: false,
				WorkgroupName:    "workgroup",
				Database:         "db",
				// ignored
				DBUser: "user",
			},
			apiInput{
				WorkgroupName: aws.String("workgroup"),
				Database:      aws.String("db"),
			},
		},
		{
			"serverless using managed secret",
			&models.RedshiftDataSourceSettings{
				UseServerless:    true,
				UseManagedSecret: true,
				WorkgroupName:    "workgroup",
				Database:         "db",
				ManagedSecret:    models.ManagedSecret{ARN: "arn:..."},
				// ignored
				DBUser: "user",
			},
			apiInput{
				WorkgroupName: aws.String("workgroup"),
				Database:      aws.String("db"),
				SecretARN:     aws.String("arn:..."),
			},
		},
		{
			"provisioned using temporary creds",
			&models.RedshiftDataSourceSettings{
				UseServerless:     false,
				UseManagedSecret:  false,
				ClusterIdentifier: "cluster",
				Database:          "db",
				DBUser:            "user",
			},
			apiInput{
				ClusterIdentifier: aws.String("cluster"),
				Database:          aws.String("db"),
				DbUser:            aws.String("user"),
			},
		},
		{
			"provisioned using managed secret",
			&models.RedshiftDataSourceSettings{
				UseServerless:     false,
				UseManagedSecret:  true,
				ClusterIdentifier: "cluster",
				Database:          "db",
				ManagedSecret:     models.ManagedSecret{ARN: "arn:..."},
				// ignored
				DBUser: "user",
			},
			apiInput{
				ClusterIdentifier: aws.String("cluster"),
				Database:          aws.String("db"),
				SecretARN:         aws.String("arn:..."),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			api := &API{settings: tt.settings}
			res := api.apiInput()
			if !cmp.Equal(res, tt.expected) {
				t.Errorf("unexpected result: %v", cmp.Diff(res, tt.expected))
			}
		})
	}
}

func Test_Execute(t *testing.T) {
	c := &API{
		settings:   &models.RedshiftDataSourceSettings{},
		DataClient: &mock.MockRedshiftClient{ExecutionResult: &redshiftdata.ExecuteStatementOutput{Id: aws.String("foo")}},
	}
	res, err := c.Execute(context.Background(), &api.ExecuteQueryInput{Query: "select * from foo"})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expectedResult := &api.ExecuteQueryOutput{ID: "foo"}
	if !cmp.Equal(expectedResult, res) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedResult, res))
	}
}

func Test_Status(t *testing.T) {
	tests := []struct {
		description string
		status      types.StatusString
		err         string
		finished    bool
	}{
		{
			description: "success",
			status:      types.StatusStringFinished,
			finished:    true,
		},
		{
			description: "error",
			status:      types.StatusStringFailed,
			err:         "boom",
			finished:    true,
		},
		{
			description: "pending",
			status:      types.StatusStringStarted,
			finished:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			c := &API{
				settings: &models.RedshiftDataSourceSettings{},
				DataClient: &mock.MockRedshiftClient{
					DescribeStatementOutput: &redshiftdata.DescribeStatementOutput{
						Id:     aws.String("foo"),
						Status: tt.status,
						Error:  aws.String(tt.err),
					},
				},
			}
			status, err := c.Status(context.Background(), &api.ExecuteQueryOutput{ID: "foo"})
			if err != nil && tt.err == "" {
				t.Errorf("unexpected error %v", err)
			}
			if status != nil && status.Finished != tt.finished {
				t.Errorf("expecting status.Finished to be %v but got %v", tt.finished, status.Finished)
			}
		})
	}
}

func Test_ListSchemas(t *testing.T) {
	resources := map[string]map[string][]string{
		"foo": {},
		"bar": {},
	}
	expectedResult := []string{"bar", "foo"}
	c := &API{
		settings:   &models.RedshiftDataSourceSettings{},
		DataClient: &mock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.Schemas(context.Background(), sqlds.Options{})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sort.Strings(res)
	if !cmp.Equal(expectedResult, res) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedResult, res))
	}
}

func Test_ListTables(t *testing.T) {
	resources := map[string]map[string][]string{
		"foo": {
			"foofoo": {},
		},
		"bar": {
			"barbar": {},
		},
	}
	expectedResult := []string{"foofoo"}
	c := &API{
		settings:   &models.RedshiftDataSourceSettings{},
		DataClient: &mock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.Tables(context.Background(), sqlds.Options{"schema": "foo"})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if !cmp.Equal(expectedResult, res) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedResult, res))
	}
}

func Test_ListColumns(t *testing.T) {
	resources := map[string]map[string][]string{
		"public": {
			"foo": {
				"col1",
				"col2",
			},
			"bar": {
				"col3",
			},
		},
	}
	expectedResult := []string{"col1", "col2"}
	c := &API{
		settings:   &models.RedshiftDataSourceSettings{},
		DataClient: &mock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.Columns(context.Background(), sqlds.Options{"schema": "public", "table": "foo"})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if !cmp.Equal(expectedResult, res) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedResult, res))
	}
}
func Test_ListSecrets(t *testing.T) {
	expectedSecrets := []models.ManagedSecret{{Name: "foo", ARN: "arn:foo"}}
	c := &API{SecretsClient: &mock.MockRedshiftSecretsManager{Secrets: []string{"foo"}}}
	secrets, err := c.Secrets(context.Background())
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if !cmp.Equal(expectedSecrets, secrets) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedSecrets, secrets))
	}
}

func Test_GetSecret(t *testing.T) {
	secretContent := `{"dbClusterIdentifier":"foo","username":"bar"}`
	c := &API{SecretsClient: &mock.MockRedshiftSecretsManager{Secret: secretContent}}
	secret, err := c.Secret(context.Background(), sqlds.Options{"secretARN": "arn"})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expectedSecret := &models.RedshiftSecret{ClusterIdentifier: "foo", DBUser: "bar"}
	if !cmp.Equal(expectedSecret, secret) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedSecret, secret))
	}
}

func Test_GetClusters(t *testing.T) {
	c := &API{ManagementClient: &mock.MockRedshiftClient{Clusters: []string{"foo", "bar"}}}
	errC := &API{ManagementClient: &mock.MockRedshiftClientError{}}
	nilC := &API{ManagementClient: &mock.MockRedshiftClientNil{}}
	expectedCluster1 := &models.RedshiftCluster{
		ClusterIdentifier: "foo",
		Endpoint: models.RedshiftEndpoint{
			Address: "foo",
			Port:    123,
		},
		Database: "foo",
	}
	expectedCluster2 := &models.RedshiftCluster{
		ClusterIdentifier: "bar",
		Endpoint: models.RedshiftEndpoint{
			Address: "bar",
			Port:    123,
		},
		Database: "bar",
	}
	tests := []struct {
		c                *API
		desc             string
		errMsg           string
		expectedClusters []models.RedshiftCluster
	}{
		{
			c:                c,
			desc:             "Happy Path",
			expectedClusters: []models.RedshiftCluster{*expectedCluster1, *expectedCluster2},
		},
		{
			c:      errC,
			desc:   "Error with DescribeCluster",
			errMsg: "Boom",
		},
		{
			c:      nilC,
			desc:   "DescribeCluster returned nil",
			errMsg: "missing clusters content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			clusters, err := tt.c.Clusters(context.Background())
			if tt.errMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedClusters, clusters)
			} else {
				assert.Nil(t, clusters)
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}

func Test_GetWorkgroups(t *testing.T) {
	c := &API{ServerlessManagementClient: &mock.MockRedshiftServerlessClient{Workgroups: []string{"foo", "bar"}}}
	errC := &API{ServerlessManagementClient: &mock.MockRedshiftServerlessClientError{}}
	nilC := &API{ServerlessManagementClient: &mock.MockRedshiftServerlessClientNil{}}
	expectedWorkgroup1 := &models.RedshiftWorkgroup{
		WorkgroupName: "foo",
		Endpoint: models.RedshiftEndpoint{
			Address: "foo",
			Port:    123,
		},
	}
	expectedWorkgroup2 := &models.RedshiftWorkgroup{
		WorkgroupName: "bar",
		Endpoint: models.RedshiftEndpoint{
			Address: "bar",
			Port:    123,
		},
	}
	tests := []struct {
		c                  *API
		desc               string
		errMsg             string
		expectedWorkgroups []models.RedshiftWorkgroup
	}{
		{
			c:                  c,
			desc:               "Happy Path",
			expectedWorkgroups: []models.RedshiftWorkgroup{*expectedWorkgroup1, *expectedWorkgroup2},
		},
		{
			c:      errC,
			desc:   "Error with DescribeWorkgroup",
			errMsg: "Boom",
		},
		{
			c:      nilC,
			desc:   "DescribeWorkgroup returned nil",
			errMsg: "missing workgroups content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			workgroups, err := tt.c.Workgroups(context.Background())
			if tt.errMsg == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedWorkgroups, workgroups)
			} else {
				assert.Nil(t, workgroups)
				assert.EqualError(t, err, tt.errMsg)
			}
		})
	}
}
