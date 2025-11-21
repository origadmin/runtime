/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package storage provides a unified interface for various storage solutions
// including cache, database, and file storage.
package storage

import (
	"cmp"
	"fmt"
	"sync"

	"github.com/goexts/generic/maps"

	cachev1 "github.com/origadmin/runtime/api/gen/go/config/data/cache/v1"
	databasev1 "github.com/origadmin/runtime/api/gen/go/config/data/database/v1"
	ossv1 "github.com/origadmin/runtime/api/gen/go/config/data/oss/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/data/storage/objectstore"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

const (
	Module = "storage"
)

type Provider interface {
	Cache(name string) (storageiface.Cache, error)
	DefaultCache() (storageiface.Cache, error)

	Database(name string) (storageiface.Database, error)
	DefaultDatabase() (storageiface.Database, error)

	ObjectStore(name string) (storageiface.ObjectStore, error)
	DefaultObjectStore() (storageiface.ObjectStore, error)
}

// providerImpl implements the storageiface.Provider interface.
type providerImpl struct {
	caches       map[string]storageiface.Cache
	databases    map[string]storageiface.Database
	objectstores map[string]storageiface.ObjectStore

	defaultCache       string
	defaultDatabase    string
	defaultObjectStore string

	mu sync.RWMutex
}

// New creates a new storage provider instance based on the provided structured configuration.
// This function is primarily for backward compatibility or when a custom interfaces.StructuredConfig
// implementation is used. For most cases, NewProvider is recommended.
func New(sc interfaces.StructuredConfig) (Provider, error) {
	dataConfig, err := sc.DecodeData()
	if err != nil {
		return nil, fmt.Errorf("failed to decode structured config: %w", err)
	}
	return NewProvider(dataConfig)
}

// NewProvider creates a new storage provider instance based on the provided decoded DataConfig.
// This is the recommended way to initialize a full storage provider when you have the
// configuration already decoded.
func NewProvider(dataConfig *datav1.Data) (Provider, error) {
	p := &providerImpl{
		caches:       make(map[string]storageiface.Cache),
		databases:    make(map[string]storageiface.Database),
		objectstores: make(map[string]storageiface.ObjectStore),
	}

	var err error

	// Initialize Caches
	p.caches, p.defaultCache, err = NewCaches(dataConfig.GetCaches())
	if err != nil {
		return nil, err
	}

	// Initialize Databases
	p.databases, p.defaultDatabase, err = NewDatabases(dataConfig.GetDatabases())
	if err != nil {
		return nil, err
	}

	// Initialize Filestores
	p.objectstores, p.defaultObjectStore, err = NewObjectStores(dataConfig.GetObjectstores())
	if err != nil {
		return nil, err
	}

	return p, nil
}

// NewCaches creates a map of cache instances and determines the default cache name
// from a datav1.Caches configuration object.
func NewCaches(cachesConfig *datav1.Caches) (map[string]storageiface.Cache, string, error) {
	if cachesConfig == nil {
		return nil, "", nil
	}

	cacheConfigsMap := maps.FromSlice(cachesConfig.GetConfigs(),
		func(cfg *cachev1.CacheConfig) (string, *cachev1.CacheConfig) {
			return cmp.Or(cfg.GetName(), cfg.GetDriver()), cfg
		})

	caches, err := NewCachesFromConfigs(cacheConfigsMap)
	if err != nil {
		return nil, "", err
	}

	defaultCacheName := cmp.Or(cachesConfig.GetActive(), cachesConfig.GetDefault(), interfaces.GlobalDefaultKey)
	if _, ok := caches[defaultCacheName]; !ok && len(caches) > 0 {
		// If default not found, but caches exist, pick the first one
		for name := range caches {
			defaultCacheName = name
			break
		}
	}
	return caches, defaultCacheName, nil
}

// NewDatabases creates a map of database instances and determines the default database name
// from a datav1.Databases configuration object.
func NewDatabases(databasesConfig *datav1.Databases) (map[string]storageiface.Database, string, error) {
	if databasesConfig == nil {
		return nil, "", nil
	}

	databaseConfigsMap := maps.FromSlice(databasesConfig.GetConfigs(),
		func(cfg *databasev1.DatabaseConfig) (string, *databasev1.DatabaseConfig) {
			return cmp.Or(cfg.GetName(), cfg.GetDialect()), cfg
		})

	databases, err := NewDatabasesFromConfigs(databaseConfigsMap)
	if err != nil {
		return nil, "", err
	}

	defaultDatabaseName := cmp.Or(databasesConfig.GetActive(), databasesConfig.GetDefault(), interfaces.GlobalDefaultKey)
	if _, ok := databases[defaultDatabaseName]; !ok && len(databases) > 0 {
		// If default not found, but databases exist, pick the first one
		for name := range databases {
			defaultDatabaseName = name
			break
		}
	}
	return databases, defaultDatabaseName, nil
}

// NewObjectStores creates a map of object store instances and determines the default object store name
// from a datav1.ObjectStores configuration object.
func NewObjectStores(objectstoresConfig *datav1.ObjectStores) (map[string]storageiface.ObjectStore, string, error) {
	if objectstoresConfig == nil {
		return nil, "", nil
	}

	objectstoreConfigsMap := maps.FromSlice(objectstoresConfig.GetConfigs(),
		func(cfg *ossv1.ObjectStoreConfig) (string, *ossv1.ObjectStoreConfig) {
			return cmp.Or(cfg.GetName(), cfg.GetDriver()), cfg
		})

	objectstores, err := NewObjectStoresFromConfigs(objectstoreConfigsMap)
	if err != nil {
		return nil, "", err
	}

	defaultObjectstoreName := cmp.Or(objectstoresConfig.GetActive(), objectstoresConfig.GetDefault(), interfaces.GlobalDefaultKey)
	if _, ok := objectstores[defaultObjectstoreName]; !ok && len(objectstores) > 0 {
		// If default not found, but objectstores exist, pick the first one
		for name := range objectstores {
			defaultObjectstoreName = name
			break
		}
	}
	return objectstores, defaultObjectstoreName, nil
}

// NewCachesFromConfigs creates a map of cache instances from a map of cache configurations.
func NewCachesFromConfigs(configs map[string]*cachev1.CacheConfig) (map[string]storageiface.Cache, error) {
	caches := make(map[string]storageiface.Cache)
	for name, cfg := range configs {
		c, err := cache.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache '%s': %w", name, err)
		}
		caches[name] = c
	}
	return caches, nil
}

// NewDatabasesFromConfigs creates a map of database instances from a map of database configurations.
func NewDatabasesFromConfigs(configs map[string]*databasev1.DatabaseConfig) (map[string]storageiface.Database, error) {
	databases := make(map[string]storageiface.Database)
	for name, cfg := range configs {
		db, err := database.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create database '%s': %w", name, err)
		}
		databases[name] = db
	}
	return databases, nil
}

// NewObjectStoresFromConfigs creates a map of object store instances from a map of object store configurations.
func NewObjectStoresFromConfigs(configs map[string]*ossv1.ObjectStoreConfig) (map[string]storageiface.ObjectStore, error) {
	objectstores := make(map[string]storageiface.ObjectStore)
	for name, cfg := range configs {
		fs, err := objectstore.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create objectstore '%s': %w", name, err)
		}
		objectstores[name] = fs
	}
	return objectstores, nil
}

// ObjectStore returns the configured file storage service by name.
func (p *providerImpl) ObjectStore(name string) (storageiface.ObjectStore, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if fs, ok := p.objectstores[name]; ok {
		return fs, nil
	}
	return nil, runtimeerrors.NewStructured(Module, "objectstore %v not found", name).WithCaller()
}

// DefaultObjectStore returns the default object storage service.
func (p *providerImpl) DefaultObjectStore() (storageiface.ObjectStore, error) {
	return p.ObjectStore(p.defaultObjectStore)
}

// Cache returns the configured cache service by name.
func (p *providerImpl) Cache(name string) (storageiface.Cache, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if c, ok := p.caches[name]; ok {
		return c, nil
	}
	return nil, runtimeerrors.NewStructured(Module, "cache %v not found", name).WithCaller()
}

// DefaultCache returns the default cache service.
func (p *providerImpl) DefaultCache() (storageiface.Cache, error) {
	return p.Cache(p.defaultCache)
}

// Database returns the configured database service by name.
func (p *providerImpl) Database(name string) (storageiface.Database, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if db, ok := p.databases[name]; ok {
		return db, nil
	}
	return nil, runtimeerrors.NewStructured(Module, "database %v not found", name).WithCaller()
}

// DefaultDatabase returns the default database service.
func (p *providerImpl) DefaultDatabase() (storageiface.Database, error) {
	return p.Database(p.defaultDatabase)
}
