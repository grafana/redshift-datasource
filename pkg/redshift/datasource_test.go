package redshift

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/sqlds/v4"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	wasCalledWith sqlds.Options
}

func (m *mockClient) Init(config backend.DataSourceInstanceSettings) {}
func (m *mockClient) GetDB(ctx context.Context, id int64, options sqlds.Options) (*sql.DB, error) {
	m.wasCalledWith = options
	return nil, nil
}
func (m *mockClient) GetAsyncDB(ctx context.Context, id int64, options sqlds.Options) (awsds.AsyncDB, error) {
	m.wasCalledWith = options
	return nil, nil
}
func (m *mockClient) GetAPI(ctx context.Context, id int64, options sqlds.Options) (sqlAPI.AWSAPI, error) {
	m.wasCalledWith = options
	return nil, errors.New("fake api error")
}

func TestConnection(t *testing.T) {
	t.Run("it should call getDB with the updated time", func(t *testing.T) {
		mc := mockClient{}
		ds := RedshiftDatasource{
			awsDS: &mc,
		}

		updatedTime := time.Now()
		fakeConfig := backend.DataSourceInstanceSettings{
			JSONData: json.RawMessage{},
			Updated:  updatedTime,
		}
		_, err := ds.Connect(context.Background(), fakeConfig, json.RawMessage(`{}`))

		assert.Nil(t, err)
		assert.Equal(t, updatedTime.String(), mc.wasCalledWith["updated"])
	})

	t.Run("it should call getAsyncDB with the updated time", func(t *testing.T) {
		mc := mockClient{}
		ds := RedshiftDatasource{
			awsDS: &mc,
		}

		updatedTime := time.Now()
		fakeConfig := backend.DataSourceInstanceSettings{
			JSONData: json.RawMessage{},
			Updated:  updatedTime,
		}
		_, err := ds.GetAsyncDB(context.Background(), fakeConfig, json.RawMessage(`{}`))

		assert.Nil(t, err)
		assert.Equal(t, updatedTime.String(), mc.wasCalledWith["updated"])
	})
}

func TestDatabases(t *testing.T) {
	t.Run("it should call getAPI with the region passed in from args", func(t *testing.T) {
		mc := mockClient{}
		ds := RedshiftDatasource{
			awsDS: &mc,
		}
		_, err := ds.Databases(context.Background(), sqlds.Options{
			"region":   "us-east1",
			"catalog":  "cat",
			"database": "db",
			"table":    "thing",
		})

		assert.Error(t, err, "fake api error", "unexpected error: %v", err)
		assert.Equal(t, "us-east1", mc.wasCalledWith["region"])
		// We can not set the config in the context, but we can confirm that updated is being added
		assert.Equal(t, "", mc.wasCalledWith["updated"])
	})
}
