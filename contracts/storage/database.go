/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package storage

import (
	"database/sql"

	databasev1 "github.com/origadmin/runtime/api/gen/go/config/data/database/v1"
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
}
