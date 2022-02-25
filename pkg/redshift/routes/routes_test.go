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
)

var ds = &fake.RedshiftFakeDatasource{
	SecretList: []models.ManagedSecret{
		{Name: "secret1", ARN: "arn:secret1"},
	},
	RSecret: models.RedshiftSecret{ClusterIdentifier: "clu", DBUser: "user"},
	RCluster: models.RedshiftCluster{
		Endpoint: models.RedshiftEndpoint{
			Address: "foo.a.b.c",
			Port: 123,
		},
		Database: "db-foo",
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
			description:    "return cluster",
			route:          "cluster",
			expectedCode:   http.StatusOK,
			expectedResult: `{"endpoint":{"address":"foo.a.b.c","port":123},"database":"db-foo"}`,
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
			case "cluster":
				rh.cluster(rw, req)
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
