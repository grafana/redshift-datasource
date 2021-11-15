package driver

import (
	"context"
	"database/sql/driver"
	"fmt"

	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	"github.com/grafana/redshift-datasource/pkg/redshift/api"
)

type conn struct {
	api    *api.API
	closed bool
}

func newConnection(api *api.API) *conn {
	return &conn{api: api}
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	output, err := c.api.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: query})
	if err != nil {
		return nil, err
	}

	if err := sqlAPI.WaitOnQuery(ctx, c.api, output); err != nil {
		return nil, err
	}

	return newRows(c.api.Client, output.ID)
}

func (c *conn) Ping(ctx context.Context) error {
	rows, err := c.QueryContext(ctx, "SELECT 1", nil)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("redshift driver doesn't support begin statements")
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("redshift driver doesn't support prepared statements")
}

func (c *conn) Close() error {
	c.closed = true
	return nil
}
