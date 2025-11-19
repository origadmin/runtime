/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package cache provides a factory function to create Cache instances.
package cache

import (
	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	"github.com/origadmin/runtime/data/storage/cache/memory"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

const (
	DefaultDriver = "memory" // Define a default driver
)

func init() {
	// Register the memory cache factory.
	Register(DefaultDriver, FactoryFunc(memory.New))
}

// Register registers a new cache factory with the default factory registry.
func Register(name string, factory Factory) {
	defaultFactory.Register(name, factory)
}

// New creates a new Cache instance based on the provided configuration.
// It uses the internal factory registry to find the appropriate provider.
// To use a specific provider (e.g., "redis"), ensure its package
// is imported for its side effects (e.g., `import _ "path/to/redis/provider"`),
// which will register the provider's factory.
func New(cfg *cachev1.CacheConfig, opts ...options.Option) (storageiface.Cache, error) {
	if cfg == nil {
		return nil, runtimeerrors.NewStructured(Module, "cache config is nil").WithCaller()
	}

	driver := cfg.GetDriver()
	if driver == "" {
		driver = DefaultDriver
	}

	// Get the factory from the registry.
	factory, ok := defaultFactory.Get(driver)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "unsupported cache driver: %s", driver).WithCaller()
	}

	// Use the factory to create the new Cache instance.
	return factory.NewCache(cfg, opts...)
}
