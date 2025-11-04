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

	cachev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/cache/v1"
	databasev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/database/v1"
	filestorev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/file/v1"
	"github.com/origadmin/runtime/data/storage/filestore"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"

	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/data/storage/database"
)

const (
	Module      = "storage"
	ErrNotFound = errors.String("storage: %s not found")
)

// providerImpl implements the storageiface.Provider interface.
type providerImpl struct {
	caches     map[string]storageiface.Cache
	databases  map[string]storageiface.Database
	filestores map[string]storageiface.FileStore

	defaultCache     string
	defaultDatabase  string
	defaultFilestore string

	mu sync.RWMutex
}

// New creates a new storage provider instance based on the provided structured configuration.
func New(sc interfaces.StructuredConfig) (storageiface.Provider, error) {
	p := &providerImpl{
		caches:     make(map[string]storageiface.Cache),
		databases:  make(map[string]storageiface.Database),
		filestores: make(map[string]storageiface.FileStore),
	}

	// Read default names from the top-level 'storage' configuration
	var defaults struct {
		DefaultCache     string `json:"defaultCache" yaml:"defaultCache"`
		DefaultDatabase  string `json:"defaultDatabase" yaml:"defaultDatabase"`
		DefaultFilestore string `json:"defaultFilestore" yaml:"defaultFilestore"`
	}

	dataConfig, err := sc.DecodeData()
	if err != nil {
		return nil, err
	}

	// Initialize Caches
	var cacheConfigs map[string]*cachev1.CacheConfig
	cacheConfigs = maps.FromSlice(dataConfig.GetCaches().GetConfigs(),
		func(cfg *cachev1.CacheConfig) (string, *cachev1.CacheConfig) {
			key := cfg.GetName()
			if key == "" {
				key = cfg.GetDriver()
			}
			return key, cfg
		})
	defaults.DefaultCache = cmp.Or(dataConfig.GetCaches().GetActive(), dataConfig.GetCaches().GetDefault(), interfaces.GlobalDefaultKey)

	for name, cacheCfg := range cacheConfigs {
		c, err := cache.New(cacheCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache '%s': %w", name, err)
		}
		p.defaultCache = name
		p.caches[name] = c
	}
	if _, ok := p.caches[defaults.DefaultCache]; ok {
		p.defaultCache = defaults.DefaultCache
	}

	// Initialize Databases
	var databaseConfigs map[string]*databasev1.DatabaseConfig
	databaseConfigs = maps.FromSlice(dataConfig.GetDatabases().GetConfigs(),
		func(cfg *databasev1.DatabaseConfig) (string, *databasev1.DatabaseConfig) {
			key := cfg.GetName()
			if key == "" {
				key = cfg.GetDialect()
			}
			return key, cfg
		})
	defaults.DefaultDatabase = cmp.Or(dataConfig.GetDatabases().GetActive(), dataConfig.GetDatabases().GetDefault(), interfaces.GlobalDefaultKey)

	for name, dbCfg := range databaseConfigs {
		db, err := database.New(dbCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create database '%s': %w", name, err)
		}
		p.defaultDatabase = name
		p.databases[name] = db
	}

	if _, ok := p.databases[defaults.DefaultDatabase]; ok {
		p.defaultDatabase = defaults.DefaultDatabase
	}

	// Initialize FileStores
	var filestoreConfigs map[string]*filestorev1.FileStoreConfig
	filestoreConfigs = maps.FromSlice(dataConfig.GetFilestores().GetConfigs(),
		func(cfg *filestorev1.FileStoreConfig) (string, *filestorev1.FileStoreConfig) {
			key := cfg.GetName()
			if key == "" {
				key = cfg.GetDriver()
			}
			return key, cfg
		})
	defaults.DefaultFilestore = cmp.Or(dataConfig.GetFilestores().GetActive(), dataConfig.GetFilestores().GetDefault(), interfaces.GlobalDefaultKey)

	for name, fsCfg := range filestoreConfigs {
		fs, err := filestore.New(fsCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create filestore '%s': %w", name, err)
		}
		p.defaultFilestore = name
		p.filestores[name] = fs
	}

	if _, ok := p.filestores[defaults.DefaultFilestore]; ok {
		p.defaultFilestore = defaults.DefaultFilestore
	}

	return p, nil
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
