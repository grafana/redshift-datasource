package schema_test

import (
	_ "embed"
	"testing"

	"github.com/grafana/dsconfig/schema"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

//go:embed dsconfig.json
var configSchemaJSON []byte

//go:generate go test -run TestPlugin -generateArtifacts
func TestPlugin(t *testing.T) {
	schema.RunPluginTests(t, schema.PluginUnderTest{
		ID:                "grafana-redshift-datasource",
		ConfigSchemaJSON:  configSchemaJSON,
		SettingsJSONModel: models.RedshiftDataSourceSettings{},
		SecureKeys:        []string{"accessKey", "secretKey", "sessionToken"},
	})
}
