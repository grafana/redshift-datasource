package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

const DriverName string = "redshift"

var (
	openFromSessionMutex sync.Mutex
	openFromSessionCount int
)

// Driver is a sql.Driver
type Driver struct {
	settings     *models.RedshiftDataSourceSettings
	sessionCache *awsds.SessionCache
}

// Open returns a new driver.Conn using already existing settings
func (d *Driver) Open(_ string) (driver.Conn, error) {
	return newConnection(d.sessionCache, d.settings), nil
}

// Open registers a new driver with a unique name
func Open(settings models.RedshiftDataSourceSettings, sessionCache *awsds.SessionCache) (*sql.DB, error) {
	openFromSessionMutex.Lock()
	openFromSessionCount++
	name := fmt.Sprintf("%s-%d", DriverName, openFromSessionCount)
	openFromSessionMutex.Unlock()
	sql.Register(name, &Driver{&settings, sessionCache})
	return sql.Open(name, "")
}
