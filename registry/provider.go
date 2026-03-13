/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	"context"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/helpers/comp"
)

// Provider defines the interface for a registry service provider.
type Provider interface {
	Registrar(name string) (KRegistrar, error)
	DefaultRegistrar() (KRegistrar, error)

	Discovery(name string) (KDiscovery, error)
	DefaultDiscovery() (KDiscovery, error)
}

// providerImpl implements the Provider interface by delegating to specific categories.
type providerImpl struct {
	locator component.Locator
}

func (p *providerImpl) Registrar(name string) (KRegistrar, error) {
	return comp.Get[KRegistrar](context.Background(), p.locator.In(CategoryRegistrar), name)
}

func (p *providerImpl) DefaultRegistrar() (KRegistrar, error) {
	return comp.GetDefault[KRegistrar](context.Background(), p.locator.In(CategoryRegistrar))
}

func (p *providerImpl) Discovery(name string) (KDiscovery, error) {
	return comp.Get[KDiscovery](context.Background(), p.locator.In(CategoryDiscovery), name)
}

func (p *providerImpl) DefaultDiscovery() (KDiscovery, error) {
	return comp.GetDefault[KDiscovery](context.Background(), p.locator.In(CategoryDiscovery))
}

// GetDiscoveries collects all discovery instances from the given locator.
func GetDiscoveries(ctx context.Context, h component.Locator) (map[string]KDiscovery, error) {
	m := make(map[string]KDiscovery)
	for name, inst := range h.Iter(ctx) {
		if d, ok := inst.(KDiscovery); ok {
			m[name] = d
		}
	}
	return m, nil
}

// NewProvider creates a new registry provider instance.
func NewProvider(locator component.Locator) Provider {
	return &providerImpl{locator: locator}
}

// DefaultRegistrarProvider creates instances for service registration.
var DefaultRegistrarProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	cfg, err := comp.AsConfig[discoveryv1.Discovery](h)
	if err != nil {
		return nil, err
	}
	return NewRegistrar(cfg)
}

// DefaultDiscoveryProvider creates instances for service discovery.
var DefaultDiscoveryProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	cfg, err := comp.AsConfig[discoveryv1.Discovery](h)
	if err != nil {
		return nil, err
	}
	return NewDiscovery(cfg)
}
