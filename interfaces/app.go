package interfaces

import (
	"time"
)

// AppInfo defines the contract for accessing application metadata.
// It provides a set of read-only methods to retrieve the application's
// static, immutable identity information. An instance of this interface
// is considered immutable once created.
type AppInfo interface {
	// ID returns the unique identifier of the application instance.
	ID() string
	// Name returns the name of the application.
	Name() string
	// Version returns the version of the application.
	Version() string
	// Env returns the environment the application is running in (e.g., "dev", "test", "prod").
	Env() string
	// StartTime returns the time when the application started.
	StartTime() time.Time
	// Metadata returns a collection of arbitrary key-value pairs.
	Metadata() map[string]string
}
