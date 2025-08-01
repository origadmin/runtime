/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/registry"
)

// NewRegistrar creates a new KRegistrar object based on the given RegistryConfig.
func (b *builder) NewRegistrar(cfg *configv1.Discovery, ss ...registry.Option) (registry.KRegistrar, error) {
	return b.RegistryBuilder.NewRegistrar(cfg, ss...)
}

// NewDiscovery creates a new discovery object based on the given RegistryConfig.
func (b *builder) NewDiscovery(cfg *configv1.Discovery, ss ...registry.Option) (registry.KDiscovery, error) {
	return b.RegistryBuilder.NewDiscovery(cfg, ss...)
}

// RegisterRegistryBuilder registers a new RegistryBuilder with the given name.
func (b *builder) RegisterRegistryBuilder(name string, factory registry.Factory) {
	b.RegistryBuilder.Register(name, factory)
}

// RegisterRegistryFunc registers a new RegistryBuilder with the given name and functions.
func (b *builder) RegisterRegistryFunc(name string, registryBuilder registry.RegistrarBuildFunc, discoveryBuilder registry.DiscoveryBuildFunc) {
	b.RegisterRegistryBuilder(name, registry.WrapFactory(registryBuilder, discoveryBuilder))
}
