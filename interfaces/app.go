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

// AppInfoBuilder defines the contract for constructing an AppInfo instance.
// This follows the builder pattern to allow for flexible and readable creation
// of an immutable AppInfo object.
type AppInfoBuilder interface {
	// WithEnv sets the environment for the application.
	WithEnv(env string) AppInfoBuilder
	// WithID sets a custom instance ID. If not called, a default (e.g., UUID) will be used.
	WithID(id string) AppInfoBuilder
	// WithStartTime sets a custom start time. If not called, the time of build will be used.
	WithStartTime(startTime time.Time) AppInfoBuilder
	// WithMetadata adds a key-value pair to the application's metadata.
	// It can be called multiple times.
	WithMetadata(key, value string) AppInfoBuilder
	// Build finalizes the construction and returns an immutable AppInfo instance.
	Build() AppInfo
}
