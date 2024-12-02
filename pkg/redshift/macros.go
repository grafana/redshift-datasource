package redshift

import (
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/gtime"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/pkg/errors"
)

func macroTimeEpoch(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.WithMessagef(sqlutil.ErrorBadArgumentCount, "expected 1 argument, received %d", len(args))
	}

	return fmt.Sprintf("extract(epoch from %s) as \"time\"", args[0]), nil
}

func macroTimeFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.WithMessagef(sqlutil.ErrorBadArgumentCount, "expected 1 argument, received %d", len(args))
	}

	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().Format(time.RFC3339)
		to     = query.TimeRange.To.UTC().Format(time.RFC3339)
	)

	return fmt.Sprintf("%s BETWEEN '%s' AND '%s'", column, from, to), nil
}

func macroTimeFrom(query *sqlutil.Query, args []string) (string, error) {
	return fmt.Sprintf("'%s'", query.TimeRange.From.UTC().Format(time.RFC3339)), nil

}

func macroTimeTo(query *sqlutil.Query, args []string) (string, error) {
	return fmt.Sprintf("'%s'", query.TimeRange.To.UTC().Format(time.RFC3339)), nil
}

func macroTimeGroup(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.WithMessagef(sqlutil.ErrorBadArgumentCount, "macro $__timeGroup needs time column and interval")
	}

	interval, err := gtime.ParseInterval(strings.Trim(args[1], `'`))
	if err != nil {
		return "", fmt.Errorf("error parsing interval %v", args[1])
	}

	return fmt.Sprintf("floor(extract(epoch from %s)/%v)*%v AS \"time\"", args[0], interval.Seconds(), interval.Seconds()), nil
}

func macroSchema(query *sqlutil.Query, args []string) (string, error) {
	return query.Schema, nil
}

func macroTable(query *sqlutil.Query, args []string) (string, error) {
	return query.Table, nil
}

func macroColumn(query *sqlutil.Query, args []string) (string, error) {
	return query.Column, nil
}

func macroUnixEpochFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.WithMessagef(sqlutil.ErrorBadArgumentCount, "expected 1 argument, received %d", len(args))
	}

	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().Unix()
		to     = query.TimeRange.To.UTC().Unix()
	)

	return fmt.Sprintf("%s >= %d AND %s <= %d", column, from, args[0], to), nil
}

func macroUnixEpochGroup(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.WithMessagef(sqlutil.ErrorBadArgumentCount, "macro $__unixEpochGroup needs time column and interval")
	}

	interval, err := gtime.ParseInterval(strings.Trim(args[1], `'`))
	if err != nil {
		return "", fmt.Errorf("error parsing interval %v", args[1])
	}

	return fmt.Sprintf(`floor(%s/%v)*%v AS "time"`, args[0], interval.Seconds(), interval.Seconds()), nil
}

var macros = map[string]sqlutil.MacroFunc{
	"timeEpoch":       macroTimeEpoch,
	"timeFilter":      macroTimeFilter,
	"timeFrom":        macroTimeFrom,
	"timeTo":          macroTimeTo,
	"timeGroup":       macroTimeGroup,
	"schema":          macroSchema,
	"table":           macroTable,
	"column":          macroColumn,
	"unixEpochFilter": macroUnixEpochFilter,
	"unixEpochGroup":  macroUnixEpochGroup,
}

func (s *RedshiftDatasource) Macros() sqlutil.Macros {
	return macros
}
