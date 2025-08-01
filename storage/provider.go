package storage

import (
	"fmt"
	"sync"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/runtime/storage/filestore"
)

// provider implements the storage.Provider interface.
type provider struct {
	cfg *configv1.Storage
	// 存储配置，而不是实例
	fileStoreConfigs map[string]*configv1.FileStore
	cacheConfigs     map[string]*configv1.Cache
	databaseConfigs  map[string]*configv1.Database

	// 缓存已初始化的实例
	initializedFileStores map[string]storageiface.FileStore
	initializedCaches     map[string]storageiface.Cache
	initializedDatabases  map[string]storageiface.Database

	defaultFileStore string
	defaultCache     string
	defaultDatabase  string

	mu sync.Mutex // 用于保护实例创建的并发访问
}

// NewProvider creates a new storage provider based on the given configuration.
func NewProvider(cfg *configv1.Storage) (storageiface.Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("storage config cannot be nil")
	}

	p := &provider{
		cfg:                   cfg,
		fileStoreConfigs:      cfg.GetFilestores(),
		cacheConfigs:          cfg.GetCaches(),
		databaseConfigs:       cfg.GetDatabases(),
		initializedFileStores: make(map[string]storageiface.FileStore),
		initializedCaches:     make(map[string]storageiface.Cache),
		initializedDatabases:  make(map[string]storageiface.Database),
		defaultFileStore:      cfg.GetDefaultFilestore(),
		defaultCache:          cfg.GetDefaultCache(),
		defaultDatabase:       cfg.GetDefaultDatabase(),
	}

	return p, nil
}

// FileStore returns the configured file storage service by name.
func (p *provider) FileStore(name string) (storageiface.FileStore, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if already initialized
	if fs, ok := p.initializedFileStores[name]; ok {
		return fs, nil
	}

	// Get configuration
	fsCfg, ok := p.fileStoreConfigs[name]
	if !ok {
		return nil, fmt.Errorf("filestore '%s' not found in configuration", name)
	}

	// Initialize and cache
	fs, err := filestore.New(fsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create filestore '%s': %w", name, err)
	}
	p.initializedFileStores[name] = fs

	return fs, nil
}

// DefaultFileStore returns the default file storage service.
func (p *provider) DefaultFileStore() (storageiface.FileStore, error) {
	if p.defaultFileStore != "" {
		return p.FileStore(p.defaultFileStore)
	}

	if len(p.fileStoreConfigs) == 1 {
		for name := range p.fileStoreConfigs {
			return p.FileStore(name) // Return the only instance if no default is specified
		}
	}

	return nil, fmt.Errorf("no default filestore configured and multiple instances exist")
}

// Cache returns the configured cache service by name.
func (p *provider) Cache(name string) (storageiface.Cache, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if already initialized
	if c, ok := p.initializedCaches[name]; ok {
		return c, nil
	}

	// Get configuration
	cacheCfg, ok := p.cacheConfigs[name]
	if !ok {
		return nil, fmt.Errorf("cache '%s' not found in configuration", name)
	}

	// Initialize and cache (TODO: Implement actual cache initialization)
	_ = cacheCfg // Suppress unused variable warning
	// c, err := cache.New(cacheCfg)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create cache '%s': %w", name, err)
	// }
	// p.initializedCaches[name] = c

	return nil, fmt.Errorf("cache initialization not yet implemented")
}

// DefaultCache returns the default cache service.
func (p *provider) DefaultCache() (storageiface.Cache, error) {
	if p.defaultCache != "" {
		return p.Cache(p.defaultCache)
	}

	if len(p.cacheConfigs) == 1 {
		for name := range p.cacheConfigs {
			return p.Cache(name)
		}
	}

	return nil, fmt.Errorf("no default cache configured and multiple instances exist")
}

// Database returns the configured database service by name.
func (p *provider) Database(name string) (storageiface.Database, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if already initialized
	if db, ok := p.initializedDatabases[name]; ok {
		return db, nil
	}

	// Get configuration
	dbCfg, ok := p.databaseConfigs[name]
	if !ok {
		return nil, fmt.Errorf("database '%s' not found in configuration", name)
	}

	// Initialize and cache (TODO: Implement actual database initialization)
	_ = dbCfg // Suppress unused variable warning
	// db, err := database.New(dbCfg)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create database '%s': %w", name, err)
	// }
	// p.initializedDatabases[name] = db

	return nil, fmt.Errorf("database initialization not yet implemented")
}

// DefaultDatabase returns the default database service.
func (p *provider) DefaultDatabase() (storageiface.Database, error) {
	if p.defaultDatabase != "" {
		return p.Database(p.defaultDatabase)
	}

	if len(p.databaseConfigs) == 1 {
		for name := range p.databaseConfigs {
			return p.Database(name)
		}
	}

	return nil, fmt.Errorf("no default database configured and multiple instances exist")
}
