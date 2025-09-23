package registry

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/registry"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
)

// NewDiscovery creates a new service discovery instance based on the provided configuration.
func NewDiscovery(cfg *discoveryv1.Discovery) (registry.Discovery, error) {
	if cfg == nil {
		return nil, fmt.Errorf("discovery config cannot be nil")
	}

	// This is where different discovery implementations (like etcd, consul, nacos) would be handled.
	switch cfg.GetType() {
	// case "etcd":
	// 	 return etcd.New(cfg.GetEtcd())
	default:
		return nil, fmt.Errorf("unsupported discovery type: %s", cfg.GetType())
	}
}

// NewRegistrar creates a new service registrar instance based on the provided configuration.
func NewRegistrar(cfg *discoveryv1.Discovery) (registry.Registrar, error) {
	if cfg == nil {
		return nil, fmt.Errorf("registrar config cannot be nil")
	}

	// This is where different registrar implementations (like etcd, consul, nacos) would be handled.
	switch cfg.GetType() {
	// case "etcd":
	// 	 return etcd.New(cfg.GetEtcd())
	default:
		return nil, fmt.Errorf("unsupported registrar type: %s", cfg.GetType())
	}
}
