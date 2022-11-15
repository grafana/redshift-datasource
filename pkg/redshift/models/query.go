package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

// FormatQueryOption defines how the user has chosen to represent the data
type FormatQueryOption uint32

const (
	// FormatOptionTimeSeries formats the query results as a timeseries using "WideToLong"
	FormatOptionTimeSeries FormatQueryOption = iota
	// FormatOptionTable formats the query results as a table using "LongToWide"
	FormatOptionTable
	// FormatOptionLogs sets the preferred visualization to logs
	FormatOptionLogs
)

// Query is the model that represents the query that users submit from the panel / queryeditor.
// For the sake of backwards compatibility, when making changes to this type, ensure that changes are
// only additive.
type Query struct {
	RawSQL         string            `json:"rawSql"`
	Format         FormatQueryOption `json:"format"`
	ConnectionArgs json.RawMessage   `json:"connectionArgs"`
	Region         string            `json:"region"`

	RefID         string            `json:"-"`
	Interval      time.Duration     `json:"-"`
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int64             `json:"-"`
	FillMissing   *data.FillMissing `json:"fillMode,omitempty"`

	// Macros
	Schema string `json:"schema,omitempty"`
	Table  string `json:"table,omitempty"`
	Column string `json:"column,omitempty"`
}

// WithSQL copies the Query, but with a different RawSQL value.
// This is mostly useful in the Interpolate function, where the RawSQL value is modified in a loop
func (q *Query) WithSQL(query string) *Query {
	return &Query{
		RawSQL:         query,
		ConnectionArgs: q.ConnectionArgs,
		RefID:          q.RefID,
		Interval:       q.Interval,
		TimeRange:      q.TimeRange,
		MaxDataPoints:  q.MaxDataPoints,
		FillMissing:    q.FillMissing,
		Schema:         q.Schema,
		Table:          q.Table,
		Column:         q.Column,
	}
}

// GetQuery returns a Query object given a backend.DataQuery using json.Unmarshal
func GetQuery(query backend.DataQuery) (*Query, error) {
	model := &Query{}

	if err := json.Unmarshal(query.JSON, &model); err != nil {
		return nil, nil
	}

	// Copy directly from the well typed query
	return &Query{
		RawSQL:         model.RawSQL,
		Format:         model.Format,
		ConnectionArgs: model.ConnectionArgs,
		RefID:          query.RefID,
		Interval:       query.Interval,
		TimeRange:      query.TimeRange,
		MaxDataPoints:  query.MaxDataPoints,
		FillMissing:    model.FillMissing,
		Schema:         model.Schema,
		Table:          model.Table,
		Column:         model.Column,
	}, nil
}

// getErrorFrameFromQuery returns a error frames with empty data and meta fields
func getErrorFrameFromQuery(query *Query) data.Frames {
	frames := data.Frames{}
	frame := data.NewFrame(query.RefID)
	frame.Meta = &data.FrameMeta{
		ExecutedQueryString: query.RawSQL,
	}
	frames = append(frames, frame)
	return frames
}

func getFrames(rows *sql.Rows, limit int64, converters []sqlutil.Converter, fillMode *data.FillMissing, query *Query) (data.Frames, error) {
	frame, err := sqlutil.FrameFromRows(rows, limit, converters...)
	if err != nil {
		return nil, err
	}
	frame.Name = query.RefID
	if frame.Meta == nil {
		frame.Meta = &data.FrameMeta{}
	}

	frame.Meta.ExecutedQueryString = query.RawSQL
	frame.Meta.PreferredVisualization = data.VisTypeGraph

	if query.Format == FormatOptionTable {
		frame.Meta.PreferredVisualization = data.VisTypeTable
		return data.Frames{frame}, nil
	}

	if query.Format == FormatOptionLogs {
		frame.Meta.PreferredVisualization = data.VisTypeLogs
		return data.Frames{frame}, nil
	}

	count, err := frame.RowLen()

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	if frame.TimeSeriesSchema().Type == data.TimeSeriesTypeLong {
		frame, err := data.LongToWide(frame, fillMode)
		if err != nil {
			return nil, err
		}
		return data.Frames{frame}, nil
	}

	return data.Frames{frame}, nil
}
