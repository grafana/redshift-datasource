package main

import (
	"embed"
	"io/fs"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/pluginschema"
)

//go:embed schema/v0alpha1/*.json
var schemaFS embed.FS

// loadSchema returns the embedded v0alpha1 PluginSchema. It uses fs.Sub to
// strip the leading "schema/" so NewCompositeFileSchemaProvider can find
// "v0alpha1/..." paths at the FS root.
func loadSchema() (*pluginschema.PluginSchema, error) {
	sub, err := fs.Sub(schemaFS, "schema")
	if err != nil {
		return nil, err
	}
	return pluginschema.NewCompositeFileSchemaProvider(sub).Get("v0alpha1")
}
