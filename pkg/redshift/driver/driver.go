package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

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
func Open(api *api.API) (*sql.DB, error) {
	openFromSessionMutex.Lock()
	openFromSessionCount++
	name := fmt.Sprintf("%s-%d", DriverName, openFromSessionCount)
	openFromSessionMutex.Unlock()
	sql.Register(name, &Driver{api})
	return sql.Open(name, "")
}
