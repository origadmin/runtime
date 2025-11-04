/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package database implements the functions, types, and interfaces for the module.
package database

import (
	"cmp"
	"database/sql"
	"time"

	databasev1 "github.com/origadmin/runtime/api/gen/go/runtime/data/database/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	storageiface "github.com/origadmin/runtime/interfaces/storage"
	"github.com/origadmin/toolkits/errors"
)

const (
	Module               = "storage.database"
	ErrDatabaseConfigNil = errors.String("database: config is nil")
)

// databaseImpl implements the storageiface.Database interface.
type databaseImpl struct {
	db      *sql.DB
	dialect string
	name    string
}

func (d *databaseImpl) Name() string {
	return d.name
}

func (d *databaseImpl) Dialect() string {
	return d.dialect
}

// DB returns the underlying *sql.DB instance.
func (d *databaseImpl) DB() *sql.DB {
	return d.db
}

// Close closes the database connection.
func (d *databaseImpl) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// New creates a new database instance based on the provided configuration.
// It uses Go's native database/sql package to open a connection.
func New(cfg *databasev1.DatabaseConfig) (storageiface.Database, error) {
	if cfg == nil {
		return nil, ErrDatabaseConfigNil
	}

	driver := cfg.GetDialect()
	if driver == "" {
		return nil, runtimeerrors.NewStructured(Module, "database driver cannot be empty").WithCaller()
	}
	source := cfg.GetSource()
	if source == "" {
		return nil, runtimeerrors.NewStructured(Module, "database source (DSN) cannot be empty").WithCaller()
	}

	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, runtimeerrors.NewStructured(Module, "failed to open database connection: %w", err).WithCaller()
	}

	// Optional: Ping the database to verify the connection is alive
	if err = db.Ping(); err != nil {
		_ = db.Close() // Close the connection if ping fails
		return nil, runtimeerrors.NewStructured(Module, "failed to ping database: %w", err).WithCaller()
	}

	// Set connection pool settings if provided in config
	if cfg.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(int(cfg.MaxOpenConnections))
	}
	if cfg.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(int(cfg.MaxIdleConnections))
	}
	if cfg.ConnectionMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnectionMaxLifetime) * time.Second)
	}
	if cfg.ConnectionMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(time.Duration(cfg.ConnectionMaxIdleTime) * time.Second)
	}

	return &databaseImpl{
		name:    cmp.Or(cfg.GetName(), cfg.GetDialect(), interfaces.GlobalDefaultKey),
		dialect: cfg.GetDialect(),
		db:      db,
	}, nil
}
