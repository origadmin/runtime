/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage provides a unified interface for various storage solutions
// including cache, database, and file storage.
package storage

import (
	"context"

	"github.com/origadmin/runtime/contracts/component"
	storageiface "github.com/origadmin/runtime/contracts/storage"
	"github.com/origadmin/runtime/engine"
)

const (
	Module = "storage"
)

// Provider defines the interface for a storage service provider.
// It acts as a bridge to the runtime engine's component container.
type Provider interface {
	Cache(name string) (storageiface.Cache, error)
	DefaultCache() (storageiface.Cache, error)

	Database(name string) (storageiface.Database, error)
	DefaultDatabase() (storageiface.Database, error)

	ObjectStore(name string) (storageiface.ObjectStore, error)
	DefaultObjectStore() (storageiface.ObjectStore, error)
}

// providerImpl implements the Provider interface by delegating to the engine handle.
type providerImpl struct {
	handle component.Handle
}

// Cache retrieves a cache instance by name from the engine.
func (p *providerImpl) Cache(name string) (storageiface.Cache, error) {
	return engine.Get[storageiface.Cache](context.Background(), p.handle.In(component.CategoryCache), name)
}

// DefaultCache retrieves the default cache instance from the engine.
func (p *providerImpl) DefaultCache() (storageiface.Cache, error) {
	return engine.GetDefault[storageiface.Cache](context.Background(), p.handle.In(component.CategoryCache))
}

// Database retrieves a database instance by name from the engine.
func (p *providerImpl) Database(name string) (storageiface.Database, error) {
	return engine.Get[storageiface.Database](context.Background(), p.handle.In(component.CategoryDatabase), name)
}

// DefaultDatabase retrieves the default database instance from the engine.
func (p *providerImpl) DefaultDatabase() (storageiface.Database, error) {
	return engine.GetDefault[storageiface.Database](context.Background(), p.handle.In(component.CategoryDatabase))
}

// ObjectStore retrieves an object store instance by name from the engine.
func (p *providerImpl) ObjectStore(name string) (storageiface.ObjectStore, error) {
	return engine.Get[storageiface.ObjectStore](context.Background(), p.handle.In(component.CategoryObjectStore), name)
}

// DefaultObjectStore retrieves the default object store instance from the engine.
func (p *providerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	return engine.GetDefault[storageiface.ObjectStore](context.Background(), p.handle.In(component.CategoryObjectStore))
}

// NewProvider creates a new storage provider instance.
// In the engine-driven architecture, it simply wraps the root component handle.
func NewProvider(handle component.Handle) Provider {
	return &providerImpl{
		handle: handle,
	}
}
