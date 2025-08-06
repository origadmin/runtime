package interfaces

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

// Resolver defines the interface for resolving configuration.
type Resolver interface {
	Resolve(kratosconfig.Config) (Resolved, error)
}

// ResolveObserver is an alias for Resolver, used for observing config changes.
type ResolveObserver = Resolver

// Resolved defines the interface for accessing resolved configuration values.
type Resolved interface {
	Services() ServiceConfig
	Discovery() DiscoveryConfig
	Middleware() MiddlewareConfig
	Logger() LoggerConfig
	Value(name string) (any, error)
	WithDecode(name string, v any, decode func([]byte, any) error) error
}
