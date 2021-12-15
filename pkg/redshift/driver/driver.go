package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"sync"

	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
)

const DriverName string = "redshift"

var (
	openFromSessionMutex sync.Mutex
	openFromSessionCount int
)

// Driver is a sql.Driver
type Driver struct {
	api *api.API
}

// Open returns a new driver.Conn using already existing settings
func (d *Driver) Open(_ string) (driver.Conn, error) {
	return newConnection(d.api), nil
}

// Open registers a new driver with a unique name
func Open(dsAPI sqlAPI.AWSAPI) (*sql.DB, error) {
	// The API is stored as a generic object but we need to parse it as a Athena API
	if reflect.TypeOf(dsAPI) != reflect.TypeOf(&api.API{}) {
		return nil, fmt.Errorf("wrong API type")
	}
	openFromSessionMutex.Lock()
	openFromSessionCount++
	name := fmt.Sprintf("%s-%d", DriverName, openFromSessionCount)
	openFromSessionMutex.Unlock()
	sql.Register(name, &Driver{api: dsAPI.(*api.API)})
	return sql.Open(name, "")
}
