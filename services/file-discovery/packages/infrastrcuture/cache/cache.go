package cache

import (
	"time"

	"github.com/abaxoth0/Vega/libs/go/packages/erorrs"
)

type Client interface {
	Connect()
	Close() *errs.Status
	Get(key string) (string, bool)
	// Uses default cache TTL (config.Cache.TTL())
	Set(key string, value any) *errs.Status
	// Same as Set(), but uses custom TTL instead of default
	SetWithTTL(key string, value any, ttl time.Duration) *errs.Status
	Delete(keys ...string) *errs.Status
}
