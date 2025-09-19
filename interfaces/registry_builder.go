package interfaces

import (
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/interfaces/factory"
)

type RegistryBuilder interface {
	factory.Registry[RegistryFactory]
	// NewRegistrar creates a new registrar using the provided ConfigDecoder and config path.
	NewRegistrar(decoder ConfigDecoder, path string, opts ...interface{}) (registry.Registrar, error)
	// NewDiscovery creates a new discovery client using the provided ConfigDecoder and config path.
	NewDiscovery(decoder ConfigDecoder, path string, opts ...interface{}) (registry.Discovery, error)
}

type RegistryFactory interface {
	// NewRegistrar creates a new registrar using the provided ConfigDecoder and config path.
	NewRegistrar(decoder ConfigDecoder, path string, opts ...interface{}) (registry.Registrar, error)
	// NewDiscovery creates a new discovery client using the provided ConfigDecoder and config path.
	NewDiscovery(decoder ConfigDecoder, path string, opts ...interface{}) (registry.Discovery, error)
}
