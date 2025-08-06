/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/factory"
)

var DefaultBuilder = NewBuilder()

type buildImpl struct {
	factory.Registry[Factory]
}

func (b *buildImpl) NewRegistrar(registry interfaces.DiscoveryConfig, opts ...interface{}) (KRegistrar, error) {
	f, ok := b.Get(registry.GetType())
	if ok {
		return f.NewRegistrar(registry, opts...)
	}
	return nil, ErrRegistryNotFound
}

func (b *buildImpl) NewDiscovery(registry interfaces.DiscoveryConfig, opts ...interface{}) (KDiscovery, error) {
	f, ok := b.Get(registry.GetType())
	if ok {
		return f.NewDiscovery(registry, opts...)
	}
	return nil, ErrRegistryNotFound
}

func NewBuilder() RegistryBuilder {
	return &buildImpl{
		Registry: factory.New[Factory](),
	}
}

// Ensure buildImpl implements interfaces.RegistryBuilder
var _ interfaces.RegistryBuilder = (*buildImpl)(nil)
