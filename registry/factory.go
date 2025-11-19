/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
	internalfactory "github.com/origadmin/runtime/internal/factory"
)

const Module = "registry"

type (
	// Factory is the interface for creating new registrar and discovery components.
	Factory interface {
		NewRegistrar(*discoveryv1.Discovery, ...options.Option) (KRegistrar, error)
		NewDiscovery(*discoveryv1.Discovery, ...options.Option) (KDiscovery, error)
	}
)

// defaultFactory is a private variable to prevent accidental modification from other packages.
var defaultFactory = newFactory()

func init() {
	// The factories will be registered here once they are updated to the new interface.
	// For example:
	// Register("consul", &consulFactory{})
	// Register("nacos", &nacosFactory{})
}

// factoryImpl is the concrete implementation of the factory registry.
type factoryImpl struct {
	factory.Registry[Factory]
}

// newFactory creates a new factoryImpl.
func newFactory() *factoryImpl {
	return &factoryImpl{
		Registry: internalfactory.New[Factory](),
	}
}

// NewRegistrar creates a new KRegistrar instance using the registered factory.
func (f *factoryImpl) NewRegistrar(cfg *discoveryv1.Discovery, opts ...options.Option) (KRegistrar, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, runtimeerrors.NewStructured(Module, "registry configuration or type is missing").WithCaller()
	}

	factory, ok := f.Get(cfg.Type)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "no registry factory found for type: %s", cfg.Type).WithCaller()
	}
	registrar, err := factory.NewRegistrar(cfg, opts...)
	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to create registrar for type %s", cfg.Type).WithCaller()
	}
	return registrar, nil
}

// NewDiscovery creates a new KDiscovery instance using the registered factory.
func (f *factoryImpl) NewDiscovery(cfg *discoveryv1.Discovery, opts ...options.Option) (KDiscovery, error) {
	if cfg == nil || cfg.Type == "" {
		return nil, runtimeerrors.NewStructured(Module, "registry configuration or type is missing").WithCaller()
	}

	factory, ok := f.Get(cfg.Type)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "no registry factory found for type: %s", cfg.Type).WithMetadata(map[string]string{"type": cfg.Type}).WithCaller()
	}
	discovery, err := factory.NewDiscovery(cfg, opts...)
	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to create discovery for type %s", cfg.Type).WithCaller()
	}
	return discovery, nil
}

// Register registers a new registry factory with the default factory.
// It is a convenience wrapper around the internal factory's Register method.
func Register(name string, factory Factory) {
	defaultFactory.Register(name, factory)
}

// NewRegistrar creates a new KRegistrar instance using the default factory.
func NewRegistrar(cfg *discoveryv1.Discovery, opts ...options.Option) (KRegistrar, error) {
	return defaultFactory.NewRegistrar(cfg, opts...)
}

// NewDiscovery creates a new KDiscovery instance using the default factory.
func NewDiscovery(cfg *discoveryv1.Discovery, opts ...options.Option) (KDiscovery, error) {
	return defaultFactory.NewDiscovery(cfg, opts...)
}
