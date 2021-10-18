package api

import (
	"context"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/google/go-cmp/cmp"
	redshiftclientmock "github.com/grafana/redshift-datasource/pkg/redshift/api/mock"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

func Test_Execute(t *testing.T) {
	expectedResult := &redshiftdataapiservice.ExecuteStatementOutput{Id: aws.String("foo")}
	c := &API{
		settings: &models.RedshiftDataSourceSettings{},
		Client:   &redshiftclientmock.MockRedshiftClient{ExecutionResult: expectedResult},
	}
	res, err := c.Execute(context.TODO(), "select * from foo")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if !cmp.Equal(expectedResult, res) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedResult, res))
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
		Client:   &redshiftclientmock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.ListSchemas(context.TODO())
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
		Client:   &redshiftclientmock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.ListTables(context.TODO(), "foo")
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
		Client:   &redshiftclientmock.MockRedshiftClient{Resources: resources},
	}
	res, err := c.ListColumns(context.TODO(), "public", "foo")
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
	secrets, err := c.ListSecrets(context.TODO())
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
	secret, err := c.GetSecret(context.TODO(), "arn")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	expectedSecret := &models.RedshiftSecret{ClusterIdentifier: "foo", DBUser: "bar"}
	if !cmp.Equal(expectedSecret, secret) {
		t.Errorf("unexpected result: %v", cmp.Diff(expectedSecret, secret))
	}
}
