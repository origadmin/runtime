package registry

import (
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
)

type Builder interface {
	factory.Registry[Factory]
	Factory
}

// --- Top-Level API ---

// Register registers a new registry factory with the DefaultBuilder.
// It is a convenience wrapper around the builder's Register method.
func Register(name string, factory Factory) {
	defaultBuilder.Register(name, factory)
}

// NewRegistrar creates a new KRegistrar instance using the DefaultBuilder.
func NewRegistrar(cfg *discoveryv1.Discovery, opts ...options.Option) (KRegistrar, error) {
	return defaultBuilder.NewRegistrar(cfg, opts...)
}

// NewDiscovery creates a new KDiscovery instance using the DefaultBuilder.
func NewDiscovery(cfg *discoveryv1.Discovery, opts ...options.Option) (KDiscovery, error) {
	return defaultBuilder.NewDiscovery(cfg, opts...)
}
