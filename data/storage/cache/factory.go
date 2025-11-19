/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package cache

import (
	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	internalfactory "github.com/origadmin/runtime/internal/factory"
)

const Module = "storage.cache"

// FactoryFunc is a function type that implements the Factory interface.
type FactoryFunc func(cfg *cachev1.CacheConfig, opts ...options.Option) (storageiface.Cache, error)

// NewCache creates a new Cache component based on the provided configuration.
func (f FactoryFunc) NewCache(cfg *cachev1.CacheConfig, opts ...options.Option) (storageiface.Cache, error) {
	return f(cfg, opts...)
}

// Factory is the interface for creating new Cache components.
type Factory interface {
	NewCache(cfg *cachev1.CacheConfig, opts ...options.Option) (storageiface.Cache, error)
}

// defaultFactory is the default, package-level instance of the cache factory registry.
var defaultFactory = internalfactory.New[Factory]()
