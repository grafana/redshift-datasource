package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/sqlds/v2"
)

type ManagedSecret struct {
	Name string `json:"name"`
	ARN  string `json:"arn"`
}

type RedshiftSecret struct {
	ClusterIdentifier string `json:"dbClusterIdentifier"`
	DBUser            string `json:"username"`
}

type RedshiftEndpoint struct {
	Address string `json:"address"`
	Port    int64  `json:"port"`
}

type RedshiftCluster struct {
	ClusterIdentifier string           `json:"clusterIdentifier"`
	Endpoint          RedshiftEndpoint `json:"endpoint"`
	Database          string           `json:"database"`
}

type RedshiftDataSourceSettings struct {
	awsds.AWSDatasourceSettings
	Config            backend.DataSourceInstanceSettings
	ClusterIdentifier string `json:"clusterIdentifier"`
	Database          string `json:"database"`
	UseManagedSecret  bool   `json:"useManagedSecret"`
	DBUser            string `json:"dbUser"`
	ManagedSecret     ManagedSecret
}

func New() models.Settings {
	return &RedshiftDataSourceSettings{}
}

func (s *RedshiftDataSourceSettings) Load(config backend.DataSourceInstanceSettings) error {
	if config.JSONData != nil && len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, s); err != nil {
			return fmt.Errorf("could not unmarshal DatasourceSettings json: %w", err)
		}
	}

	s.AccessKey = config.DecryptedSecureJSONData["accessKey"]
	s.SecretKey = config.DecryptedSecureJSONData["secretKey"]

	s.Config = config

	return nil
}

func (s *RedshiftDataSourceSettings) Apply(args sqlds.Options) {
	region, database := args["region"], args["database"]
	if region != "" {
		if region == models.DefaultKey {
			s.Region = s.DefaultRegion
		} else {
			s.Region = region
		}
	}

	if database != "" && database != models.DefaultKey {
		s.Database = database
	}
}
