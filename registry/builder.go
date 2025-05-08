/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/builder"
)

type buildImpl struct {
	builder.Builder[Factory]
}

func (b *buildImpl) NewRegistrar(registry *configv1.Registry, opts ...Option) (KRegistrar, error) {
	factory, ok := b.Get(registry.Type)
	if ok {
		return factory.NewRegistrar(registry, opts...)
	}
	return nil, ErrRegistryNotFound
}

func (b *buildImpl) NewDiscovery(registry *configv1.Registry, opts ...Option) (KDiscovery, error) {
	factory, ok := b.Get(registry.Type)
	if ok {
		return factory.NewDiscovery(registry, opts...)
	}
	return nil, ErrRegistryNotFound
}

func NewBuilder() Builder {
	return &buildImpl{
		Builder: builder.New[Factory](),
	}
}
