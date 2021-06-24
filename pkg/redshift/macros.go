package redshift

import (
	"fmt"
	"time"

	"github.com/grafana/sqlds"
	"github.com/pkg/errors"
)

func macroTimeFilter(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.WithMessagef(sqlds.ErrorBadArgumentCount, "expected 1 argument, received %d", len(args))
	}

	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().Format(time.RFC3339)
		to     = query.TimeRange.To.UTC().Format(time.RFC3339)
	)

	return fmt.Sprintf("%s BETWEEN '%s' AND '%s'", column, from, to), nil
}

func macroTimeFrom(query *sqlds.Query, args []string) (string, error) {
	return fmt.Sprintf("'%s'", query.TimeRange.From.UTC().Format(time.RFC3339)), nil

}

func macroTimeTo(query *sqlds.Query, args []string) (string, error) {
	return fmt.Sprintf("'%s'", query.TimeRange.To.UTC().Format(time.RFC3339)), nil
}

func macroTimeGroup(query *sqlds.Query, args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.WithMessagef(sqlds.ErrorBadArgumentCount, "macro timeGroup expects 2 arguments, received %d", len(args))
	}

	return fmt.Sprintf("date_trunc(%s, %s)", args[1], args[0]), nil
}

var macros = map[string]sqlds.MacroFunc{
	"timeFilter": macroTimeFilter,
	"timeFrom":   macroTimeFrom,
	"timeTo":     macroTimeTo,
	"timeGroup":  macroTimeGroup,
}

func (s *RedshiftDatasource) Macros() sqlds.Macros {
	return macros
}
