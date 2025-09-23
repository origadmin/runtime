package interfaces

import (
	"time"
)

// AppInfo represents the application's metadata.
// It is now a native Go struct, not a protobuf message.
type AppInfo struct {
	// ID is the unique identifier of the application instance.
	ID string
	// Name is the name of the application.
	Name string
	// Version is the version of the application.
	Version string
	// Env is the environment the application is running in (e.g., "dev", "test", "prod").
	Env string
	// StartTime is the time when the application started.
	StartTime time.Time
	// Metadata is a collection of arbitrary key-value pairs.
	Metadata map[string]string
}
