/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package registry

import (
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
)

const (
	Module = "registry.factory"
)

// Factory is the interface for creating new registrar and discovery components.
type Factory interface {
	NewRegistrar(*discoveryv1.Discovery, ...options.Option) (KRegistrar, error)
	NewDiscovery(*discoveryv1.Discovery, ...options.Option) (KDiscovery, error)
}

// buildImpl is the concrete implementation of the Builder.
type buildImpl struct {
	factory.Registry[Factory]
}

func (b *buildImpl) NewRegistrar(cfg *discoveryv1.Discovery, opts ...options.Option) (KRegistrar, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, runtimeerrors.NewStructured(Module, "registry configuration or type is missing").WithCaller()
	}

	f, ok := b.Get(cfg.Type)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "no registry factory found for type: %s", cfg.Type).WithMetadata(map[string]string{"type": cfg.Type}).WithCaller()
	}
	registrar, err := f.NewRegistrar(cfg, opts...)
	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to create registrar for type %s", cfg.Type).WithCaller()
	}
	return registrar, nil
}

func (b *buildImpl) NewDiscovery(cfg *discoveryv1.Discovery, opts ...options.Option) (KDiscovery, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, runtimeerrors.NewStructured(Module, "registry configuration or type is missing").WithCaller()
	}

	f, ok := b.Get(cfg.Type)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "no registry factory found for type: %s", cfg.Type).WithMetadata(map[string]string{"type": cfg.Type}).WithCaller()
	}
	discovery, err := f.NewDiscovery(cfg, opts...)
	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to create discovery for type %s", cfg.Type).WithCaller()
	}
	return discovery, nil
}

// defaultBuilder is a private variable to prevent accidental modification from other packages.
var defaultBuilder = &buildImpl{
	Registry: factory.New[Factory](),
}

// DefaultBuilder returns the shared instance of the registry builder.
func DefaultBuilder() Builder {
	return defaultBuilder
}
