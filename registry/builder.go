/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"sync"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type builder struct {
	factoryMux sync.RWMutex
	factories  map[string]Factory
}

func (b *builder) RegisterRegistryBuilder(name string, factory Factory) {
	b.factoryMux.Lock()
	defer b.factoryMux.Unlock()
	b.factories[name] = factory
}

func (b *builder) NewRegistrar(registry *configv1.Registry, setting ...OptionSetting) (Registrar, error) {
	b.factoryMux.RLock()
	defer b.factoryMux.RUnlock()
	if r, ok := b.factories[registry.Type]; ok {
		return r.NewRegistrar(registry, setting...)
	}
	return nil, ErrRegistryNotFound
}

func (b *builder) NewDiscovery(registry *configv1.Registry, setting ...OptionSetting) (Discovery, error) {
	b.factoryMux.RLock()
	defer b.factoryMux.RUnlock()
	if r, ok := b.factories[registry.Type]; ok {
		return r.NewDiscovery(registry, setting...)
	}
	return nil, ErrRegistryNotFound
}
func NewBuilder() Builder {
	return &builder{
		factories: make(map[string]Factory),
	}
}
