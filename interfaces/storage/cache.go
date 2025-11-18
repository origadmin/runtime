/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
)

var (
	// cacheBuilders is a global registry for Cache builders.
	cacheBuilders = make(map[string]CacheBuilder)
	// cacheMu protects the global cache registry.
	cacheMu sync.RWMutex
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
	// New builds a new Cache instance from the given configuration.
	New(cfg *cachev1.CacheConfig) (Cache, error)
	// Removed Name() string from the interface.
}

// RegisterCache registers a new CacheBuilder with a given name.
// This function is typically called from the init() function of a storage provider package.
// If a builder with the same name is already registered, it will panic.
func RegisterCache(name string, b CacheBuilder) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if _, exists := cacheBuilders[name]; exists {
		panic(fmt.Sprintf("storage: Cache builder named %s already registered", name))
	}
	cacheBuilders[name] = b
}

// GetCacheBuilder retrieves a registered CacheBuilder by name.
// It returns the builder and true if found, otherwise nil and false.
func GetCacheBuilder(name string) (CacheBuilder, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	b, ok := cacheBuilders[name]
	return b, ok
}
