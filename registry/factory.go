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

const (
	Module = "registry.factory"
)

type (
	// RegistryFactory is the interface for creating new registrar and discovery components.
	RegistryFactory interface {
		NewRegistrar(*discoveryv1.Discovery, ...options.Option) (KRegistrar, error)
		NewDiscovery(*discoveryv1.Discovery, ...options.Option) (KDiscovery, error)
	}
)

// defaultRegistryBuilder is a private variable to prevent accidental modification from other packages.
var defaultRegistryBuilder = NewBuilder()

func init() {
	// The factories will be registered here once they are updated to the new interface.
	// For example:
	// Register("consul", &consulFactory{})
	// Register("nacos", &nacosFactory{})
}

// registryBuilder is the concrete implementation of the Builder.
type registryBuilder struct {
	factory.Registry[RegistryFactory]
}

// NewRegistrar creates a new KRegistrar instance using the registered factory.
func (b *registryBuilder) NewRegistrar(cfg *discoveryv1.Discovery, opts ...options.Option) (KRegistrar, error) {
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

// NewDiscovery creates a new KDiscovery instance using the registered factory.
func (b *registryBuilder) NewDiscovery(cfg *discoveryv1.Discovery, opts ...options.Option) (KDiscovery, error) {
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

// Register registers a registry factory with the given name.
func Register(name string, factory RegistryFactory) {
	defaultRegistryBuilder.Register(name, factory)
}

// NewRegistrar creates a new KRegistrar instance using the DefaultRegistryBuilder.
func NewRegistrar(cfg *discoveryv1.Discovery, opts ...options.Option) (KRegistrar, error) {
	return defaultRegistryBuilder.NewRegistrar(cfg, opts...)
}

// NewDiscovery creates a new KDiscovery instance using the DefaultRegistryBuilder.
func NewDiscovery(cfg *discoveryv1.Discovery, opts ...options.Option) (KDiscovery, error) {
	return defaultRegistryBuilder.NewDiscovery(cfg, opts...)
}

// NewBuilder creates a new registryBuilder.
func NewBuilder() Builder {
	return &registryBuilder{
		Registry: internalfactory.New[RegistryFactory](),
	}
}

// Builder is an interface that defines a method for registering a buildImpl.
// This interface is kept for backward compatibility but its usage is discouraged.
// All new code should use the package-level functions.
type Builder interface {
	factory.Registry[RegistryFactory]
	NewRegistrar(*discoveryv1.Discovery, ...options.Option) (KRegistrar, error)
	NewDiscovery(*discoveryv1.Discovery, ...options.Option) (KDiscovery, error)
}