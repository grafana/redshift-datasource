package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/sqlds/v4"
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
	Port    int32  `json:"port"`
}

type RedshiftCluster struct {
	ClusterIdentifier string           `json:"clusterIdentifier"`
	Endpoint          RedshiftEndpoint `json:"endpoint"`
	Database          string           `json:"database"`
}

type RedshiftWorkgroup struct {
	WorkgroupName string           `json:"workgroupName"`
	Endpoint      RedshiftEndpoint `json:"endpoint"`
	Database      string           `json:"database"`
}

type RedshiftDataSourceSettings struct {
	awsds.AWSDatasourceSettings
	Config            backend.DataSourceInstanceSettings
	ClusterIdentifier string `json:"clusterIdentifier"`
	WorkgroupName     string `json:"workgroupName"`
	Database          string `json:"database"`
	UseServerless     bool   `json:"useServerless"`
	UseManagedSecret  bool   `json:"useManagedSecret"`
	WithEvent         bool   `json:"withEvent"`
	DBUser            string `json:"dbUser"`
	ManagedSecret     ManagedSecret
}

func New(_ context.Context) models.Settings {
	return &RedshiftDataSourceSettings{}
}

func (s *RedshiftDataSourceSettings) Load(config backend.DataSourceInstanceSettings) error {
	if len(config.JSONData) > 1 {
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
