package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type ManagedSecret struct {
	Name string `json:"name"`
	ARN  string `json:"arn"`
}

type RedshiftSecret struct {
	ClusterIdentifier string `json:"dbClusterIdentifier"`
	DBUser            string `json:"username"`
}

type RedshiftDataSourceSettings struct {
	awsds.AWSDatasourceSettings
	ClusterIdentifier string `json:"clusterIdentifier"`
	Database          string `json:"database"`
	UseManagedSecret  bool   `json:"useManagedSecret"`
	DBUser            string `json:"dbUser"`
	ManagedSecret     ManagedSecret
}

func (s *RedshiftDataSourceSettings) Load(config backend.DataSourceInstanceSettings) error {
	if config.JSONData != nil && len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, s); err != nil {
			return fmt.Errorf("could not unmarshal DatasourceSettings json: %w", err)
		}
	}

	s.AccessKey = config.DecryptedSecureJSONData["accessKey"]
	s.SecretKey = config.DecryptedSecureJSONData["secretKey"]

	return nil
}
