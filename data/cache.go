/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package data implements the functions, types, and interfaces for the module.
package data

import (
	storagev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/storage/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"

	"github.com/origadmin/runtime/data/cache"
)

const (
	CacheModule       = "storage.cache"
	ErrCacheConfigNil = errors.String("cache: config is nil")
)

// New creates a new cache instance based on the provided configuration.
func New(cfg *storagev1.CacheConfig) (storageiface.Cache, error) {
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
		return nil, runtimeerrors.NewStructured(CacheModule, "unsupported cache driver: %s", cfg.GetDriver()).WithCaller()
	}
}
