package storage

import (
	"fmt"

	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/storage/cache"
)

const (
	ErrCacheConfigNil = errors.String("cache: config is nil")
)

// New creates a new cache instance based on the provided configuration.
func New(cfg *configv1.Cache) (storageiface.Cache, error) {
	if cfg == nil {
		return nil, ErrCacheConfigNil
	}

	switch cfg.GetDriver() {
	case "memory":
		return cache.NewMemoryCache(cfg.GetMemory()), nil // Pass the Memory config
	// case "redis":
	//     return redis.New(cfg.GetRedis()), nil
	// case "memcached":
	//     return memcached.New(cfg.GetMemcached()), nil
	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", cfg.GetDriver())
	}
}
