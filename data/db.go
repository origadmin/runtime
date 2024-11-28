/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package database implements the functions, types, and interfaces for the module.
package data

import (
	"database/sql"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/errors"
)

const (
	ErrDatabaseConfigNil = errors.String("database: config is nil")
)

func OpenDB(database *configv1.Data_Database) (*sql.DB, error) {
	if database == nil {
		return nil, ErrDatabaseConfigNil
	}

	db, err := sql.Open(database.Driver, database.Source)
	if err != nil {
		return nil, err
	}
	if database.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(int(database.MaxOpenConnections))
	}
	if database.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(int(database.MaxIdleConnections))
	}
	if t := database.ConnectionMaxLifetime; t != nil && t.AsDuration() > 0 {
		db.SetConnMaxLifetime(t.AsDuration())
	}
	if t := database.ConnectionMaxIdleTime; t != nil && t.AsDuration() > 0 {
		db.SetConnMaxIdleTime(t.AsDuration())
	}

	return db, nil
}
