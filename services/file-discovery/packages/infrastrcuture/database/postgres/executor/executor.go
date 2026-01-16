package executor

import (
	"context"
	"errors"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
	"vega_file_discovery/common/config"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/connection"
	dblog "vega_file_discovery/packages/infrastrcuture/database/postgres/db-logger"
	"vega_file_discovery/packages/infrastrcuture/database/postgres/query"

	errs "github.com/abaxoth0/Vega/libs/go/packages/erorrs"
	"github.com/jackc/pgx/v5"
)

var conManager *connection.Manager

func Init(manager *connection.Manager) {
	if manager == nil {
		dblog.Logger.Panic(
			"Failed to initlized database query executor",
			"Connetion manager is nil",
			nil,
		)
	}
	conManager = manager
}

// Creates new execution context.
// Instead of newExecutionContext() which just creates it, this function make it ready-to-use.
func initExecutionContext(conType connection.Type, q *query.Query) (*executionContext, context.CancelFunc, *errs.Status) {
	con, err := conManager.AcquireConnection(conType)
	if err != nil {
		return nil, nil, err
	}

	if config.Debug.Enabled && config.Debug.LogDbQueries {
		args := make([]string, len(q.Args))

		for i, arg := range q.Args {
			switch a := arg.(type) {
			case string:
				args[i] = a
			case []string:
				args[i] = strings.Join(a, ", ")
			case int:
				args[i] = strconv.FormatInt(int64(a), 10)
			case int64:
				args[i] = strconv.FormatInt(a, 10)
			case int32:
				args[i] = strconv.FormatInt(int64(a), 10)
			case float32:
				args[i] = strconv.FormatFloat(float64(a), 'f', 8, 32)
			case float64:
				args[i] = strconv.FormatFloat(float64(a), 'f', 11, 64)
			case time.Time:
				args[i] = a.String()
			case *time.Time:
				args[i] = a.String()
			case bool:
				args[i] = strconv.FormatBool(a)
			case net.IP:
				args[i] = a.To4().String()
			}
		}

		dblog.Logger.Debug("Running query:\n"+q.SQL+"\n * Query args: "+strings.Join(args, "; "), nil)
	}

	ctx, cancel := newExecutionContext(context.Background(), time.Second*5, con)

	return ctx, cancel, nil
}

func Rows(conType connection.Type, query *query.Query) (pgx.Rows, *errs.Status) {
	ctx, cancel, err := initExecutionContext(conType, query)
	if err != nil {
		return nil, err
	}
	defer cancel()

	r, e := ctx.Connection.Query(ctx, query.SQL, query.Args...)
	if e != nil {
		return nil, query.ConvertAndLogError(e)
	}

	return r, nil
}

// Scans a row into the given destinations.
// All dests must be pointers.
// By default, dests validation is disabled,
// to enable this add "debug-safe-db-scans: true" to the config.
// (works only if app launched in debug mode)
type rowScanner = func(dests ...any) *errs.Status

// Wrapper for '*pgxpool.Con.QueryRow'
func Row(conType connection.Type, query *query.Query) (rowScanner, *errs.Status) {
	ctx, cancel, err := initExecutionContext(conType, query)
	if err != nil {
		return nil, err
	}
	defer cancel()

	row := ctx.Connection.QueryRow(ctx, query.SQL, query.Args...)

	return func(dests ...any) *errs.Status {
		dblog.Logger.Trace("Scanning row...", nil)

		if config.Debug.Enabled && config.Debug.SafeDatabaseScans {
			for _, dest := range dests {
				typeof := reflect.TypeOf(dest)

				if typeof.Kind() != reflect.Ptr {
					dblog.Logger.Panic(
						"Query scan failed",
						"Destination for scanning must be a pointer, but got '"+typeof.String()+"'",
						nil,
					)
				}
			}
		}

		if e := row.Scan(dests...); e != nil {
			if errors.Is(e, pgx.ErrNoRows) {
				return errs.StatusNotFound
			}
			return query.ConvertAndLogError(e)
		}

		dblog.Logger.Trace("Scanning row: OK", nil)

		return nil
	}, nil
}

// Wrapper for '*pgxpool.Con.Exec'
func Exec(conType connection.Type, query *query.Query) *errs.Status {
	ctx, cancel, err := initExecutionContext(conType, query)
	if err != nil {
		return err
	}
	defer cancel()

	if _, err := ctx.Connection.Exec(ctx, query.SQL, query.Args...); err != nil {
		return query.ConvertAndLogError(err)
	}

	return nil
}
