package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"sync"

	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	sqlDriver "github.com/grafana/grafana-aws-sdk/pkg/sql/driver"
	asyncSQLDriver "github.com/grafana/grafana-aws-sdk/pkg/sql/driver/async"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
	"github.com/grafana/sqlds/v2"
)

const DriverName string = "redshift"

var (
	openFromSessionMutex sync.Mutex
	openFromSessionCount int
)

// Driver is a sql.Driver
type Driver struct {
	name       string
	api        *api.API
	connection *conn
	asyncDB    *DB
}

// Open returns a new driver.Conn using already existing settings
func (d *Driver) Open(_ string) (driver.Conn, error) {
	d.connection = newConnection(d.api)
	return d.connection, nil
}

func (d *Driver) Closed() bool {
	return d.connection == nil || d.connection.closed
}

func (d *Driver) OpenDB() (*sql.DB, error) {
	return sql.Open(d.name, "")
}

func (d *Driver) GetAsyncDB() (sqlds.AsyncDB, error) {
	return d.asyncDB, nil
}

// New registers a new driver with a unique name
func New(dsAPI sqlAPI.AWSAPI) (sqlDriver.Driver, error) {
	return NewAsync(dsAPI)
}

func NewAsync(dsAPI sqlAPI.AWSAPI) (asyncSQLDriver.AsyncDriver, error) {
	// The API is stored as a generic object but we need to parse it as a Redshift API
	if reflect.TypeOf(dsAPI) != reflect.TypeOf(&api.API{}) {
		return nil, fmt.Errorf("wrong API type")
	}
	openFromSessionMutex.Lock()
	openFromSessionCount++
	name := fmt.Sprintf("%s-%d", DriverName, openFromSessionCount)
	openFromSessionMutex.Unlock()
	d := &Driver{api: dsAPI.(*api.API), name: name}
	d.asyncDB = &DB{api: d.api}
	sql.Register(name, d)
	return d, nil
}
