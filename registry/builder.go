/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"sync"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/grpc"
	"github.com/origadmin/runtime/service/http"
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

func (b *builder) NewRegistrar(registry *configv1.Registry, httpSetting []service.HTTPOptionSetting, grpcSetting []service.GRPCOptionSetting) (KRegistrar, error) {
	b.factoryMux.RLock()
	defer b.factoryMux.RUnlock()

	// 调用 Configure 接口的 FromConfig 方法
	for _, opt := range httpSetting {
		if opt != nil {
			option := &http.Options{}
			opt(option)
			if option.Configure != nil {
				if err := option.Configure.FromConfig(registry); err != nil {
					return nil, err
				}
			}
		}
	}

	for _, opt := range grpcSetting {
		if opt != nil {
			option := &grpc.Options{}
			opt(option)
			if option.Configure != nil {
				if err := option.Configure.FromConfig(registry); err != nil {
					return nil, err
				}
			}
		}
	}

	if r, ok := b.factories[registry.Type]; ok {
		return r.NewRegistrar(registry, httpSetting, grpcSetting...)
	}
	return nil, ErrRegistryNotFound
}

func (b *builder) NewDiscovery(registry *configv1.Registry, httpSetting ...service.HTTPOptionSetting, grpcSetting ...service.GRPCOptionSetting) (KDiscovery, error) {
	b.factoryMux.RLock()
	defer b.factoryMux.RUnlock()

	// 调用 Configure 接口的 FromConfig 方法
	for _, opt := range httpSetting {
		if opt != nil {
			option := &service.HTTPOption{}
			opt(option)
			if option.Configure != nil {
				if err := option.Configure.FromConfig(registry); err != nil {
					return nil, err
				}
			}
		}
	}

	for _, opt := range grpcSetting {
		if opt != nil {
			option := &service.GRPCOption{}
			opt(option)
			if option.Configure != nil {
				if err := option.Configure.FromConfig(registry); err != nil {
					return nil, err
				}
			}
		}
	}

	if r, ok := b.factories[registry.Type]; ok {
		return r.NewDiscovery(registry, httpSetting, grpcSetting...)
	}
	return nil, ErrRegistryNotFound
}

func NewBuilder() Builder {
	return &builder{
		factories: make(map[string]Factory),
	}
}
