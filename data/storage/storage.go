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
	"github.com/origadmin/runtime/helpers/comp"
)

const (
	// CategoryDatabase is the category for database components.
	CategoryDatabase component.Category = "database"
	// CategoryCache is the category for cache components.
	CategoryCache component.Category = "cache"
	// CategoryObjectStore is the category for objectstore components.
	CategoryObjectStore component.Category = "objectstore"
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

// providerImpl implements the Provider interface by delegating to the engine locator.
type providerImpl struct {
	l component.Locator
}

// Cache retrieves a cache instance by name from the engine.
func (p *providerImpl) Cache(name string) (storageiface.Cache, error) {
	return comp.Get[storageiface.Cache](context.Background(), p.l.In(CategoryCache), name)
}

// DefaultCache retrieves the default cache instance from the engine.
func (p *providerImpl) DefaultCache() (storageiface.Cache, error) {
	return comp.Get[storageiface.Cache](context.Background(), p.l.In(CategoryCache))
}

// Database retrieves a database instance by name from the engine.
func (p *providerImpl) Database(name string) (storageiface.Database, error) {
	return comp.Get[storageiface.Database](context.Background(), p.l.In(CategoryDatabase), name)
}

// DefaultDatabase retrieves the default database instance from the engine.
func (p *providerImpl) DefaultDatabase() (storageiface.Database, error) {
	return comp.Get[storageiface.Database](context.Background(), p.l.In(CategoryDatabase))
}

// ObjectStore retrieves an object store instance by name from the engine.
func (p *providerImpl) ObjectStore(name string) (storageiface.ObjectStore, error) {
	return comp.Get[storageiface.ObjectStore](context.Background(), p.l.In(CategoryObjectStore), name)
}

// DefaultObjectStore retrieves the default object store instance from the engine.
func (p *providerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	return comp.Get[storageiface.ObjectStore](context.Background(), p.l.In(CategoryObjectStore))
}

// NewProvider creates a new storage provider instance.
func NewProvider(l component.Locator) Provider {
	return &providerImpl{
		l: l,
	}
}
