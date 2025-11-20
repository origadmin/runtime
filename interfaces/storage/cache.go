/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package storage

import (
	"context"
	"time"

	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
)

// Cache defines the interface for a key-value cache.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetAndDelete(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value string, exp ...time.Duration) error
	Delete(ctx context.Context, key string) error
	Close(ctx context.Context) error
	Clear(ctx context.Context) error
}

// CacheBuilder defines the interface for building a Cache instance.
type CacheBuilder interface {
	// NewCache builds a new Cache instance from the given configuration.
	NewCache(cfg *cachev1.CacheConfig) (Cache, error)
}
