/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package cache implements the functions, types, and interfaces for the module.
package cache

import (
	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"

	"github.com/origadmin/runtime/data/storage/cache/memory" // Import the memory cache implementation
)

const (
	Module            = "storage.cache"
	ErrCacheConfigNil = errors.String("cache: config is nil")
)

// New creates a new cache instance based on the provided configuration.
func New(cfg *cachev1.CacheConfig) (storageiface.Cache, error) {
	if cfg == nil {
		return nil, ErrCacheConfigNil
	}

	switch cfg.GetDriver() {
	case "memory":
		return memory.NewCache(cfg.GetMemory()), nil // Call the NewMemoryCache from the memory subpackage
	// case "redis":
	//     return redis.New(cfg.GetRedis()), nil
	// case "memcached":
	//     return memcached.New(cfg.GetMemcached()), nil
	default:
		return nil, runtimeerrors.NewStructured(Module, "unsupported cache driver: %s", cfg.GetDriver()).WithCaller()
	}
}
