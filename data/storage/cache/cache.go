/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package cache provides a factory function to create Cache instances.
package cache

import (
	"fmt" // Added for error formatting
	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"

	// Import the memory provider for its side effects (registration).
	_ "github.com/origadmin/runtime/data/storage/cache/memory"
)

const (
	Module            = "storage.cache"
	ErrCacheConfigNil = errors.String("cache: config is nil")
	DefaultDriver     = "memory" // Define a default driver
)

// New creates a new Cache instance based on the provided configuration.
// It uses a registry of builders to find the appropriate provider.
// To use a specific provider (e.g., "redis"), ensure its package
// is imported for its side effects (e.g., `import _ "path/to/redis/provider"`).
func New(cfg *cachev1.CacheConfig) (storageiface.Cache, error) {
	if cfg == nil {
		return nil, ErrCacheConfigNil
	}

	driver := cfg.GetDriver()
	if driver == "" {
		driver = DefaultDriver
	}

	// Get the builder from the registry.
	builder, ok := storageiface.GetCacheBuilder(driver)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "unsupported cache driver: %s", driver).WithCaller()
	}

	// Use the builder to create the new Cache instance.
	return builder.New(cfg)
}
