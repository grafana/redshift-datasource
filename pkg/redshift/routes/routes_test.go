package routes

import (
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
			expectedResult: `[{"arn":"arn:secret1","name":"secret1"}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			rw := httptest.NewRecorder()
			rh := ResourceHandler{ds: ds}
			switch tt.route {
			case "secretsmanager":
				rh.secrets(rw, req)
			default:
				t.Fatalf("unexpected route %s", tt.route)
			}

			resp := rw.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expecting code %v got %v", tt.expectedCode, resp.StatusCode)
			}
			if resp.StatusCode == http.StatusOK && !cmp.Equal(string(body), tt.expectedResult) {
				t.Errorf("unexpected response: %v", cmp.Diff(string(body), tt.expectedResult))
			}
		})
	}
}
