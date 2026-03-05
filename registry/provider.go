/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"context"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
)

// Provider defines the interface for a registry service provider.
type Provider interface {
	Registrar(name string) (KRegistrar, error)
	DefaultRegistrar() (KRegistrar, error)

	Discovery(name string) (KDiscovery, error)
	DefaultDiscovery() (KDiscovery, error)
}

// providerImpl implements the Provider interface.
type providerImpl struct {
	handle component.Handle
}

func (p *providerImpl) Registrar(name string) (KRegistrar, error) {
	return engine.Get[KRegistrar](context.Background(), p.handle.In(component.CategoryRegistry), name)
}

func (p *providerImpl) DefaultRegistrar() (KRegistrar, error) {
	return engine.GetDefault[KRegistrar](context.Background(), p.handle.In(component.CategoryRegistry))
}

func (p *providerImpl) Discovery(name string) (KDiscovery, error) {
	return engine.Get[KDiscovery](context.Background(), p.handle.In(component.CategoryRegistry), name)
}

func (p *providerImpl) DefaultDiscovery() (KDiscovery, error) {
	return engine.GetDefault[KDiscovery](context.Background(), p.handle.In(component.CategoryRegistry))
}

// NewProvider creates a new registry provider instance.
func NewProvider(handle component.Handle) Provider {
	return &providerImpl{handle: handle}
}

func init() {
	engine.Register(component.CategoryRegistry, DefaultProvider)
}

// RegistrarProvider is the engine-compatible provider for registrar components.
var RegistrarProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	cfg, err := engine.AsConfig[discoveryv1.Discovery](h)
	if err != nil {
		return nil, err
	}
	return NewRegistrar(cfg, opts...)
}

// DiscoveryProvider is the engine-compatible provider for discovery components.
var DiscoveryProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	cfg, err := engine.AsConfig[discoveryv1.Discovery](h)
	if err != nil {
		return nil, err
	}
	return NewDiscovery(cfg, opts...)
}

// DefaultProvider acts as a general provider.
var DefaultProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	cfg, err := engine.AsConfig[discoveryv1.Discovery](h)
	if err != nil {
		return nil, err
	}
	// Many registry implementations (like Consul, ETCD) implement both interfaces in one object.
	return NewRegistrar(cfg, opts...)
}
