package redis

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
	"vega_file_discovery/common/config"

	"github.com/abaxoth0/Vega/libs/go/packages/erorrs"
	"github.com/abaxoth0/Vega/libs/go/packages/logger"
	"github.com/redis/go-redis/v9"
)

var log = logger.NewSource("CACHE", logger.Default)

type driver struct {
	client      *redis.Client
	isConnected bool
}

func New() *driver {
	return new(driver)
}

func (d *driver) Connect() {
	if d.isConnected {
		log.Panic("DB connection failed", "Connection already established", nil)
	}

	log.Info("Connecting to DB...", nil)

	d.client = redis.NewClient(&redis.Options{
		PoolSize:     20 * runtime.NumCPU(),
		MinIdleConns: 10,
		Addr:         config.Secret.CacheURI,
		Password:     config.Secret.CachePassword,
		DB:           config.Secret.CacheDB,
		ReadTimeout:  config.Cache.PoolTimeout() / 2,
		PoolTimeout:  config.Cache.PoolTimeout(),
	})

	ctx, cancel := defaultTimeoutContext()
	defer cancel()

	if err := d.client.Ping(ctx).Err(); err != nil {
		log.Panic("DB connection failed", err.Error(), nil)
	}

	log.Info("Connecting to DB: OK", nil)

	d.isConnected = true
}

func (d *driver) Close() *errs.Status {
	if !d.isConnected {
		return errs.NewStatusError(
			"connection not established",
			http.StatusInternalServerError,
		)
	}

	log.Info("Disconnecting from DB...", nil)

	if err := d.client.Close(); err != nil {
		return errs.NewStatusError(
			err.Error(),
			http.StatusInternalServerError,
		)
	}

	log.Info("Disconnecting from DB: OK", nil)

	d.isConnected = false

	return nil
}

func defaultTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), config.Cache.OperationTimeout())
}

// Performs logging.
// Returns err converted to *Error.Status if err is not nil.
// If err is nil then just logs given action (with trace level).
func handleError(action string, err error) *errs.Status {
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Error(
				"Request failed",
				"TIMEOUT: "+action,
				nil,
			)
		} else {
			log.Error(
				"Request failed",
				"Failed to "+action+": "+err.Error(),
				nil,
			)
		}
		return errs.StatusInternalServerError
	}

	log.Trace(action, nil)

	return nil
}

func (d *driver) Get(key string) (string, bool) {
	ctx, cancel := defaultTimeoutContext()
	defer cancel()

	cachedData, err := d.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Trace("Miss: "+key, nil)
		return "", false
	}

	return cachedData, handleError("Get: "+key, err) == nil
}

const maxRetries = 4

func (d *driver) retry(fn func(ctx context.Context) error) error {
	var lastErr error

	for i := range maxRetries {
		ctx, cancel := defaultTimeoutContext()
		defer cancel()

		err := fn(ctx)
		if err == nil {
			log.Trace("Operation succeeded with "+strconv.Itoa(i)+" retry(-s)", nil)
			return nil
		}

		lastErr = err

		// Exponential time gap between retries
		backoff := time.Duration(math.Pow(2, float64(i))) * time.Millisecond
		// Random value that will be added to backoff to get total time gap
		jitter := time.Duration(rand.Intn(10)) * time.Millisecond

		time.Sleep(backoff + jitter)
	}

	return lastErr
}

// IMPORTANT:
// go-redis driver can handle only this types:
// string, bool, []byte, int, int64, float64, time.Time
func (d *driver) set(key string, value any, ttl time.Duration) *errs.Status {
	// Alas, generics can't be used in methods
	// (it can be passed to a struct, but thats kinda strange and
	//  even so i failed to make it works as i want, so using type switch instead)
	switch v := value.(type) {
	case string, bool, []byte, int, int64, float64, time.Time:
	// Type allowed, do nothing and just go forward
	case uint32:
		if uint64(v) > uint64(math.MaxInt64) {
			return handleError("Set: ", fmt.Errorf("value overflows int64: %v", value))
		}
		value = int64(v)
	case uint64:
		if v > uint64(math.MaxInt64) {
			return handleError("Set: ", fmt.Errorf("value overflows int64: %v", value))
		}
		value = int64(v)
	default:
		return handleError("Set: ", fmt.Errorf("invalid cache value type: %T", value))
	}

	err := d.retry(func(ctx context.Context) error {
		return d.client.Set(ctx, key, value, ttl).Err()
	})

	return handleError("Set: "+key, err)
}

func (d *driver) Set(key string, value any) *errs.Status {
	return d.set(key, value, config.Cache.TTL())
}

func (d *driver) SetWithTTL(key string, value any, ttl time.Duration) *errs.Status {
	return d.set(key, value, ttl)
}

func (d *driver) Delete(keys ...string) *errs.Status {
	err := d.retry(func(ctx context.Context) error {
		return d.client.Unlink(ctx, keys...).Err()
	})
	return handleError("Delete: "+strings.Join(keys, ","), err)
}
