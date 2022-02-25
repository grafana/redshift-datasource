package api

import (
	"context"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	redshiftclientmock "github.com/grafana/redshift-datasource/pkg/redshift/api/mock"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v2"
)

func Test_apiInput(t *testing.T) {
	tests := []struct {
		description string
		settings    *models.RedshiftDataSourceSettings
		expected    apiInput
	}{
		{
			"using temporary creds",
			&models.RedshiftDataSourceSettings{
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
			"using managed secret",
			&models.RedshiftDataSourceSettings{
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
		settings: &models.RedshiftDataSourceSettings{},
		DataClient:   &redshiftclientmock.MockRedshiftClient{ExecutionResult: &redshiftdataapiservice.ExecuteStatementOutput{Id: aws.String("foo")}},
	}
	res, err := c.Execute(context.TODO(), &api.ExecuteQueryInput{Query: "select * from foo"})
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
		status      string
		err         string
		finished    bool
	}{
		{
			description: "success",
			status:      redshiftdataapiservice.StatusStringFinished,
			finished:    true,
		},
		{
			description: "error",
			status:      redshiftdataapiservice.StatusStringFailed,
			err:         "boom",
			finished:    true,
		},
		{
			description: "pending",
			status:      redshiftdataapiservice.StatusStringStarted,
			finished:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			c := &API{
				settings: &models.RedshiftDataSourceSettings{},
				DataClient: &redshiftclientmock.MockRedshiftClient{
					DescribeStatementOutput: &redshiftdataapiservice.DescribeStatementOutput{
						Id:     aws.String("foo"),
						Status: aws.String(tt.status),
						Error:  aws.String(tt.err),
					},
				},
			}
			status, err := c.Status(context.TODO(), &api.ExecuteQueryOutput{ID: "foo"})
			if err != nil && tt.err == "" {
				t.Errorf("unexpected error %v", err)
			}
			if status.Finished != tt.finished {
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
		settings: &models.RedshiftDataSourceSettings{},
		DataClient:   &redshiftclientmock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.Schemas(context.TODO(), sqlds.Options{})
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
		settings: &models.RedshiftDataSourceSettings{},
		DataClient:   &redshiftclientmock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.Tables(context.TODO(), sqlds.Options{"schema": "foo"})
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
		settings: &models.RedshiftDataSourceSettings{},
		DataClient:   &redshiftclientmock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.Columns(context.TODO(), sqlds.Options{"schema": "public", "table": "foo"})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if !cmp.Equal(expectedResult, res) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedResult, res))
	}
}
func Test_ListSecrets(t *testing.T) {
	expectedSecrets := []models.ManagedSecret{{Name: "foo", ARN: "arn:foo"}}
	c := &API{SecretsClient: &redshiftclientmock.MockRedshiftClient{Secrets: []string{"foo"}}}
	secrets, err := c.Secrets(context.TODO())
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if !cmp.Equal(expectedSecrets, secrets) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedSecrets, secrets))
	}
}

func Test_GetSecret(t *testing.T) {
	secretContent := `{"dbClusterIdentifier":"foo","username":"bar"}`
	c := &API{SecretsClient: &redshiftclientmock.MockRedshiftClient{Secret: secretContent}}
	secret, err := c.Secret(context.TODO(), sqlds.Options{"secretARN": "arn"})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expectedSecret := &models.RedshiftSecret{ClusterIdentifier: "foo", DBUser: "bar"}
	if !cmp.Equal(expectedSecret, secret) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedSecret, secret))
	}
}

func Test_GetCluster(t *testing.T) {
	c := &API{ManagementClient: &redshiftclientmock.MockRedshiftClient{Clusters: []string{"foo"}}}
	cluster, err := c.Cluster(sqlds.Options{"clusterIdentifier": "foo"})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expectedCluster := &models.RedshiftCluster{
		Endpoint: models.RedshiftEndpoint{ 
			Address: "foo", 
			Port: 123,
		}, 
		Database: "foo",
	}
	if !cmp.Equal(expectedCluster, cluster) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedCluster, cluster))
	}
}
