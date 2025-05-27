/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package discovery implements the functions, types, and interfaces for the module.
package registry

import (
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

type buildImpl struct {
	factory.Registry[Factory]
}

func (b *buildImpl) NewRegistrar(registry *configv1.Discovery, opts ...Option) (KRegistrar, error) {
	f, ok := b.Get(registry.Type)
	if ok {
		return f.NewRegistrar(registry, opts...)
	}
	return nil, ErrRegistryNotFound
}

func (b *buildImpl) NewDiscovery(registry *configv1.Discovery, opts ...Option) (KDiscovery, error) {
	factory, ok := b.Get(registry.Type)
	if ok {
		return factory.NewDiscovery(registry, opts...)
	}
	return nil, ErrRegistryNotFound
}

func NewBuilder() Builder {
	return &buildImpl{
		Registry: factory.New[Factory](),
	}
}
