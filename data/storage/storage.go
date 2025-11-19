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
	filestorev1 "github.com/origadmin/runtime/api/gen/go/config/data/file/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/data/storage/filestore"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
)

const (
	Module = "storage"
)

// providerImpl implements the storageiface.Provider interface.
type providerImpl struct {
	caches     map[string]storageiface.Cache
	databases  map[string]storageiface.Database
	filestores map[string]storageiface.ObjectStore

	defaultCache     string
	defaultDatabase  string
	defaultFilestore string

	mu sync.RWMutex
}

// New creates a new storage provider instance based on the provided structured configuration.
// This function is primarily for backward compatibility or when a custom interfaces.StructuredConfig
// implementation is used. For most cases, NewProvider is recommended.
func New(sc interfaces.StructuredConfig) (storageiface.Provider, error) {
	dataConfig, err := sc.DecodeData()
	if err != nil {
		return nil, fmt.Errorf("failed to decode structured config: %w", err)
	}
	return NewProvider(dataConfig)
}

// NewProvider creates a new storage provider instance based on the provided decoded DataConfig.
// This is the recommended way to initialize a full storage provider when you have the
// configuration already decoded.
func NewProvider(dataConfig *datav1.Data) (storageiface.Provider, error) {
	p := &providerImpl{
		caches:     make(map[string]storageiface.Cache),
		databases:  make(map[string]storageiface.Database),
		filestores: make(map[string]storageiface.FileStore),
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
	p.filestores, p.defaultFilestore, err = NewFilestores(dataConfig.GetFilestores())
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

// NewFilestores creates a map of file store instances and determines the default file store name
// from a datav1.FileStores configuration object.
func NewFilestores(filestoresConfig *datav1.Filestores) (map[string]storageiface.FileStore, string, error) {
	if filestoresConfig == nil {
		return nil, "", nil
	}

	filestoreConfigsMap := maps.FromSlice(filestoresConfig.GetConfigs(),
		func(cfg *filestorev1.FilestoreConfig) (string, *filestorev1.FilestoreConfig) {
			return cmp.Or(cfg.GetName(), cfg.GetDriver()), cfg
		})

	filestores, err := NewFileStoresFromConfigs(filestoreConfigsMap)
	if err != nil {
		return nil, "", err
	}

	defaultFilestoreName := cmp.Or(filestoresConfig.GetActive(), filestoresConfig.GetDefault(), interfaces.GlobalDefaultKey)
	if _, ok := filestores[defaultFilestoreName]; !ok && len(filestores) > 0 {
		// If default not found, but filestores exist, pick the first one
		for name := range filestores {
			defaultFilestoreName = name
			break
		}
	}
	return filestores, defaultFilestoreName, nil
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

// NewFileStoresFromConfigs creates a map of file store instances from a map of file store configurations.
func NewFileStoresFromConfigs(configs map[string]*filestorev1.FilestoreConfig) (map[string]storageiface.FileStore, error) {
	filestores := make(map[string]storageiface.FileStore)
	for name, cfg := range configs {
		fs, err := filestore.New(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create filestore '%s': %w", name, err)
		}
		filestores[name] = fs
	}
	return filestores, nil
}

// FileStore returns the configured file storage service by name.
func (p *providerImpl) FileStore(name string) (storageiface.FileStore, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if fs, ok := p.filestores[name]; ok {
		return fs, nil
	}
	return nil, runtimeerrors.NewStructured(Module, "filestore %v not found", name).WithCaller()
}

// DefaultFileStore returns the default file storage service.
func (p *providerImpl) DefaultFileStore() (storageiface.FileStore, error) {
	return p.FileStore(p.defaultFilestore)
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
