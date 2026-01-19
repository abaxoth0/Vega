package query

import (
	"context"

	common "github.com/abaxoth0/Vega/libs/go/packages"
	errs "github.com/abaxoth0/Vega/libs/go/packages/erorrs"
	"github.com/abaxoth0/Vega/libs/go/packages/logger"
)

var queryLogger = logger.NewSource("QUERY", logger.Default)

type Query struct {
	SQL  string
	Args []any
}

func New(sql string, args ...any) *Query {
	return &Query{
		SQL:  sql,
		Args: args,
	}
}

func (q *Query) ConvertAndLogError(err error) *errs.Status {
	defer queryLogger.Debug("Failed query: "+q.SQL, nil)

	if err == context.DeadlineExceeded {
		queryLogger.Error("Query failed", "Operation timeout", nil)
		return errs.StatusTimeout
	}

	queryLogger.Error("Query failed", err.Error(), nil)
	return errs.StatusInternalServerError
}

// Used for positional parameters that can be zero/empty and must be set to NULL in this case (e.g. empty strings).
//
// If specified value is zero, then this function will return value which will be inserted in query as NULL.
//
// For example: "NULL" for strings and nil for empty slices.
//
// P.S. Built-in NULLIF() function is not used cuz it involves extra cost for query execution.
// (affects query planing, SARGability, caching, requires type casting et cetera)
func Nullif[T any](value T) T {
	switch v := any(value).(type) {
	case string:
		return common.Ternary(v == "", any("NULL").(T), value)
	case []string:
		if len(v) == 0 {
			var zero T
			return zero
		}
	}
	return value
}
