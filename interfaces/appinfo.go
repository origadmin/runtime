package interfaces

import (
	"time"
)

// AppInfo represents the application's metadata.
// It is now a native Go struct, not a protobuf message.
type AppInfo struct {
	// ID is the unique identifier of the application instance.
	ID string `json:"id" yaml:"id" mapstructure:"id"`
	// Name is the name of the application.
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// Version is the version of the application.
	Version string `json:"version" yaml:"version" mapstructure:"version"`
	// Env is the environment the application is running in (e.g., "dev", "test", "prod").
	Env string `json:"env" yaml:"env" mapstructure:"env"`
	// StartTime is the time when the application started. This is a runtime generated value.
	StartTime time.Time `json:"-" yaml:"-" mapstructure:"-"` // Mark as non-configurable
	// Metadata is a collection of arbitrary key-value pairs.
	Metadata map[string]string `json:"metadata" yaml:"metadata" mapstructure:"metadata"`
}
