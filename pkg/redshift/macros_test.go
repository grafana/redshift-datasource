package redshift

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/pkg/errors"
)

func Test_macros(t *testing.T) {
	tests := []struct {
		description string
		macro       string
		query       *sqlutil.Query
		args        []string
		expected    string
		expectedErr error
	}{
		{
			"adds time as Unix",
			"timeEpoch",
			&sqlutil.Query{},
			[]string{"starttime"},
			`extract(epoch from starttime) as "time"`,
			nil,
		},
		{
			"creates time filter",
			"timeFilter",
			&sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2021, 6, 23, 0, 0, 0, 0, &time.Location{}),
					To:   time.Date(2021, 6, 23, 1, 0, 0, 0, &time.Location{}),
				},
			},
			[]string{"starttime"},
			`starttime BETWEEN '2021-06-23T00:00:00Z' AND '2021-06-23T01:00:00Z'`,
			nil,
		},
		{
			"wrong args for time filter",
			"timeFilter",
			&sqlutil.Query{},
			[]string{},
			"",
			sqlutil.ErrorBadArgumentCount,
		},
		{
			"creates time from filter",
			"timeFrom",
			&sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2021, 6, 23, 0, 0, 0, 0, &time.Location{}),
					To:   time.Date(2021, 6, 23, 1, 0, 0, 0, &time.Location{}),
				},
			},
			[]string{},
			`'2021-06-23T00:00:00Z'`,
			nil,
		},
		{
			"creates time to filter",
			"timeTo",
			&sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2021, 6, 23, 0, 0, 0, 0, &time.Location{}),
					To:   time.Date(2021, 6, 23, 1, 0, 0, 0, &time.Location{}),
				},
			},
			[]string{},
			`'2021-06-23T01:00:00Z'`,
			nil,
		},
		{
			"creates time group",
			"timeGroup",
			&sqlutil.Query{},
			[]string{"starttime", "'1m'"},
			`floor(extract(epoch from starttime)/60)*60 AS "time"`,
			nil,
		},
		{
			"wrong args for time group",
			"timeGroup",
			&sqlutil.Query{},
			[]string{},
			"",
			sqlutil.ErrorBadArgumentCount,
		},
		{
			"adds a schema",
			"schema",
			&sqlutil.Query{Schema: "foo"},
			[]string{},
			`foo`,
			nil,
		},
		{
			"adds a table",
			"table",
			&sqlutil.Query{Table: "foo"},
			[]string{},
			`foo`,
			nil,
		},
		{
			"adds a column",
			"column",
			&sqlutil.Query{Column: "foo"},
			[]string{},
			`foo`,
			nil,
		},
		{
			"unix epoch filter",
			"unixEpochFilter",
			&sqlutil.Query{
				TimeRange: backend.TimeRange{
					From: time.Date(2021, 6, 23, 0, 0, 0, 0, &time.Location{}),
					To:   time.Date(2021, 6, 23, 1, 0, 0, 0, &time.Location{}),
				},
			},
			[]string{"starttime"},
			`starttime >= 1624406400 AND starttime <= 1624410000`,
			nil,
		},
		{
			"unix epoch time group",
			"unixEpochGroup",
			&sqlutil.Query{},
			[]string{"starttime", "1h"},
			`floor(starttime/3600)*3600 AS "time"`,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			res, err := macros[tt.macro](tt.query, tt.args)
			if (err != nil || tt.expectedErr != nil) && !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error %v, expecting %v", err, tt.expectedErr)
			}
			if res != tt.expected {
				t.Errorf("unexpected result %v, expecting %v", res, tt.expected)
			}
		})
	}
}
