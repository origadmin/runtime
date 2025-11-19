package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	sourcev1 "github.com/origadmin/runtime/api/gen/go/config/source/v1"
	"github.com/origadmin/runtime/interfaces/options"
)

// SourceFactory is the interface for creating configuration sources.
// It defines a single method, NewSource, which creates a new config source
// based on the provided configuration and options.
type SourceFactory interface {
	// NewSource creates a new config source.
	NewSource(*sourcev1.SourceConfig, ...options.Option) (kratosconfig.Source, error)
}
