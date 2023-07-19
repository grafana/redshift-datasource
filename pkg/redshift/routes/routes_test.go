package routes

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/redshift-datasource/pkg/redshift/fake"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/stretchr/testify/assert"
)

var ds = &fake.RedshiftFakeDatasource{
	SecretList: []models.ManagedSecret{
		{Name: "secret1", ARN: "arn:secret1"},
	},
	RSecret: models.RedshiftSecret{ClusterIdentifier: "clu", DBUser: "user"},
	RClusters: []models.RedshiftCluster{
		{
			ClusterIdentifier: "foo",
			Endpoint: models.RedshiftEndpoint{
				Address: "foo.a.b.c",
				Port:    123,
			},
			Database: "db-foo",
		},
	},
	RWorkgroups: []models.RedshiftWorkgroup{
		{
			WorkgroupName: "bar",
			Endpoint: models.RedshiftEndpoint{
				Address: "bar.a.b.c",
				Port:    456,
			},
			Database: "db-bar",
		},
	},
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		description    string
		route          string
		expectedCode   int
		expectedResult string
	}{
		{
			description:    "return secrets",
			route:          "secrets",
			expectedCode:   http.StatusOK,
			expectedResult: `[{"name":"secret1","arn":"arn:secret1"}]`,
		},
		{
			description:    "return secret",
			route:          "secret",
			expectedCode:   http.StatusOK,
			expectedResult: `{"dbClusterIdentifier":"clu","username":"user"}`,
		},
		{
			description:    "return clusters",
			route:          "clusters",
			expectedCode:   http.StatusOK,
			expectedResult: `[{"clusterIdentifier":"foo","endpoint":{"address":"foo.a.b.c","port":123},"database":"db-foo"}]`,
		},
		{
			description:    "return workgroups",
			route:          "workgroups",
			expectedCode:   http.StatusOK,
			expectedResult: `[{"workgroupName":"bar","endpoint":{"address":"bar.a.b.c","port":456},"database":"db-bar"}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/foo", bytes.NewReader([]byte("{}")))
			rw := httptest.NewRecorder()
			rh := RedshiftResourceHandler{redshift: ds}
			switch tt.route {
			case "secrets":
				rh.secrets(rw, req)
			case "secret":
				rh.secret(rw, req)
			case "clusters":
				rh.clusters(rw, req)
			case "workgroups":
				rh.workgroups(rw, req)
			default:
				t.Fatalf("unexpected route %s", tt.route)
			}

			resp := rw.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expecting code %v got %v. Body: %v", tt.expectedCode, resp.StatusCode, string(body))
			}
			if resp.StatusCode == http.StatusOK && !cmp.Equal(string(body), tt.expectedResult) {
				t.Errorf("unexpected response: %v", cmp.Diff(string(body), tt.expectedResult))
			}
		})
	}
}

func Test_Routes(t *testing.T) {
	rh := RedshiftResourceHandler{redshift: ds}
	r := rh.Routes()
	assert.Contains(t, r, "/secrets")
	assert.Contains(t, r, "/secret")
	assert.Contains(t, r, "/workgroups")
	assert.Contains(t, r, "/clusters")
}
