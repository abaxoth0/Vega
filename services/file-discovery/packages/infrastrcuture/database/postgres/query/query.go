package query

import (
	"context"

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
