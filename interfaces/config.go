package interfaces

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

// Resolved is the main interface for accessing resolved configuration values.
// It provides a flexible way to scan configuration sections into custom Go structs.
type Resolved interface {
	// Decode decodes the configuration section identified by 'key' into the 'target' Go struct.
	// The 'key' can be a dot-separated path (e.g., "service.http", "data.database").
	// 'target' must be a pointer to a Go struct.
	Decode(key string, target interface{}) error
}

// ConfigLoader defines the interface for loading application configuration.
// DEPRECATED: This interface is being phased out in favor of the new bootstrap.Load mechanism.
type ConfigLoader interface {
	Load(configPath string, bootstrapConfig interface{}) (kratosconfig.Config, error)
}
