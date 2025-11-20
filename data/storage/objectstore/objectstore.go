/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package objectstore provides a factory function to create ObjectStore instances.
package objectstore

import (
	ossv1 "github.com/origadmin/runtime/api/gen/go/config/data/oss/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/options"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

const (
	DefaultDriver = "local" // Define a default driver
)

// Register registers a new object store factory with the default factory registry.
func Register(name string, factory Factory) {
	defaultFactory.Register(name, factory)
}

// New creates a new ObjectStore instance based on the provided configuration.
// It uses the internal factory registry to find the appropriate provider.
// To use a specific provider (e.g., "minio"), ensure its package
// is imported for its side effects (e.g., `import _ "path/to/minio/provider"`),
// which will register the provider's factory.
func New(cfg *ossv1.ObjectStoreConfig, option ...options.Option) (storageiface.ObjectStore, error) {
	if cfg == nil {
		return nil, runtimeerrors.NewStructured(Module, "object store config is nil").WithCaller()
	}

	driver := cfg.GetDriver()
	if driver == "" {
		driver = DefaultDriver
	}

	// Get the factory from the registry.
	factory, ok := defaultFactory.Get(driver)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "unsupported object store driver: %s", driver).WithCaller()
	}

	// Use the factory to create the new ObjectStore instance.
	return factory.New(cfg)
}
