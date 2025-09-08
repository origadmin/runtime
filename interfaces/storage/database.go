package storage

import (
	"database/sql"
)

// Database defines the interface for a database service.
type Database interface {
	// DB returns the underlying *sql.DB instance.
	DB() *sql.DB

	// Close closes the database connection.
	Close() error
}
