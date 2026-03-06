/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package database implements the functions, types, and contracts for the module.
package database

import (
	"cmp"
	"database/sql"
	"time"

	databasev1 "github.com/origadmin/runtime/api/gen/go/config/data/database/v1"
	"github.com/origadmin/runtime/contracts"
	"github.com/origadmin/runtime/contracts/options"
	storageiface "github.com/origadmin/runtime/contracts/storage"
	runtimeerrors "github.com/origadmin/runtime/errors"
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

func (d *databaseImpl) Name() string    { return d.name }
func (d *databaseImpl) Dialect() string { return d.dialect }
func (d *databaseImpl) DB() *sql.DB     { return d.db }
func (d *databaseImpl) Close() error    { return d.db.Close() }

func New(cfg *databasev1.DatabaseConfig, opts ...options.Option) (storageiface.Database, error) {
	if cfg == nil {
		return nil, ErrDatabaseConfigNil
	}

	db, err := sql.Open(cfg.GetDialect(), cfg.GetSource())
	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to open database").WithCaller()
	}

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
		name:    cmp.Or(cfg.GetName(), cfg.GetDialect(), contracts.GlobalDefaultKey),
		dialect: cfg.GetDialect(),
		db:      db,
	}, nil
}
