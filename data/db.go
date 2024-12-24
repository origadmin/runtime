/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package database implements the functions, types, and interfaces for the module.
package data

import (
	"database/sql"
	"time"

	"github.com/origadmin/toolkits/errors"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

const (
	ErrDatabaseConfigNil = errors.String("database: config is nil")
)

func OpenDB(database *configv1.Data_Database) (*sql.DB, error) {
	if database == nil {
		return nil, ErrDatabaseConfigNil
	}

	db, err := sql.Open(database.Dialect, database.Source)
	if err != nil {
		return nil, err
	}
	if database.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(int(database.MaxOpenConnections))
	}
	if database.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(int(database.MaxIdleConnections))
	}
	if t := database.ConnectionMaxLifetime; t != 0 {
		db.SetConnMaxLifetime(time.Duration(t))
	}
	if t := database.ConnectionMaxIdleTime; t != 0 {
		db.SetConnMaxIdleTime(time.Duration(t))
	}
	return db, nil
}
