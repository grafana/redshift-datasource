package driver

// import (
// 	"context"
// 	"database/sql/driver"
// 	"fmt"

// 	sqlAPI "github.com/grafana/grafana-aws-sdk/pkg/sql/api"
// 	"github.com/grafana/redshift-datasource/pkg/redshift/api"
// )

// type conn struct {
// 	api    *api.API
// 	closed bool
// }

// func newConnection(api *api.API) *conn {
// 	return &conn{api: api}
// }

// func (c *conn) StartQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
// 	output, err := c.api.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: query})
// 	if err != nil {
// 		return "", err
// 	}
// 	return output.ID, nil
// }

// func (c *conn) QueryStatus(ctx context.Context, queryID string) (bool, string, error) {
// 	status, err := c.api.Status(ctx, &sqlAPI.ExecuteQueryOutput{ID: queryID})
// 	if err != nil {
// 		return false, "", err
// 	}
// 	return status.Finished, status.State, nil
// }

// func (c *conn) CancelQuery(ctx context.Context, queryID string) error {
// 	return c.api.Stop(&sqlAPI.ExecuteQueryOutput{ID: queryID})
// }

// func (c *conn) GetRows(ctx context.Context, queryID string) (driver.Rows, error) {
// 	return newRows(c.api.DataClient, queryID)
// }

// func (c *conn) Ping() error {
// 	_, err := c.api.Execute(context.Background(), &sqlAPI.ExecuteQueryInput{Query: "SELECT 1"})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *conn) PingContext(ctx context.Context) error {
// 	_, err := c.api.Execute(ctx, &sqlAPI.ExecuteQueryInput{Query: "SELECT 1"})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *conn) Begin() (driver.Tx, error) {
// 	return nil, fmt.Errorf("redshift driver doesn't support begin statements")
// }

// func (c *conn) Prepare(query string) (driver.Stmt, error) {
// 	return nil, fmt.Errorf("redshift driver doesn't support prepared statements")
// }

// func (c *conn) Close() error {
// 	c.closed = true
// 	return nil
// }
