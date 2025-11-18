/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package storage

import (
	"database/sql"
	"fmt"
	"sync"

	databasev1 "github.com/origadmin/runtime/api/gen/go/config/data/database/v1"
)

var (
	// databaseBuilders is a global registry for Database builders.
	databaseBuilders = make(map[string]DatabaseBuilder)
	// databaseMu protects the global database registry.
	databaseMu sync.RWMutex
)

// Database defines the interface for a database service.
type Database interface {
	Name() string
	Dialect() string
	DB() *sql.DB
	Close() error
}

// DatabaseBuilder defines the interface for building a Database instance.
type DatabaseBuilder interface {
	// New builds a new Database instance from the given configuration.
	New(cfg *databasev1.DatabaseConfig) (Database, error)
	// Name returns the name of the builder (e.g., "mysql", "postgres").
	// Note: This refers to the dialect/driver name.
	Name() string
}

// RegisterDatabase registers a new DatabaseBuilder.
func RegisterDatabase(b DatabaseBuilder) {
	databaseMu.Lock()
	defer databaseMu.Unlock()

	name := b.Name()
	if _, exists := databaseBuilders[name]; exists {
		panic(fmt.Sprintf("storage: Database builder named %s already registered", name))
	}
	databaseBuilders[name] = b
}

// GetDatabaseBuilder retrieves a registered DatabaseBuilder by name.
func GetDatabaseBuilder(name string) (DatabaseBuilder, bool) {
	databaseMu.RLock()
	defer databaseMu.RUnlock()
	b, ok := databaseBuilders[name]
	return b, ok
}
