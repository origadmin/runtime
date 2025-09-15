package interfaces

import (
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/interfaces/factory"
)

type RegistryBuilder interface {
	factory.Registry[RegistryFactory]
	NewRegistrar(DiscoveryConfig, ...interface{}) (registry.Registrar, error)
	NewDiscovery(DiscoveryConfig, ...interface{}) (registry.Discovery, error)
}

type RegistryFactory interface {
	NewRegistrar(DiscoveryConfig, ...interface{}) (registry.Registrar, error)
	NewDiscovery(DiscoveryConfig, ...interface{}) (registry.Discovery, error)
}
