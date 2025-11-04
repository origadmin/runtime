package storage

import (
	"database/sql"
)

// Database defines the interface for a database service.
type Database interface {
	// Name returns the name of the database service.
	Name() string

	// Dialect returns the database dialect.
	Dialect() string

	// DB returns the underlying *sql.DB instance.
	DB() *sql.DB

	// Close closes the database connection.
	Close() error
}
