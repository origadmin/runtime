/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"sync"

	"github.com/goexts/generic/settings"

	"github.com/origadmin/runtime/config"
)

// builder is a struct that holds a map of ConfigBuilders and a map of RegistryBuilders.
type builder struct {
	configMux     sync.RWMutex
	config        *config.RuntimeConfig
	builderMux    sync.RWMutex
	builders      map[string]ConfigBuilder
	syncMux       sync.RWMutex
	syncs         map[string]ConfigSyncer
	registryMux   sync.RWMutex
	registries    map[string]RegistryBuilder
	serviceMux    sync.RWMutex
	services      map[string]ServiceBuilder
	middlewareMux sync.RWMutex
	middlewares   map[string]MiddlewareBuilder
}

func (b *builder) Config() *config.RuntimeConfig {
	b.configMux.RLock()
	defer b.configMux.RUnlock()
	return b.config
}

func (b *builder) ApplyConfig(ss ...config.RuntimeConfigSetting) {
	b.configMux.Lock()
	defer b.configMux.Unlock()
	b.config = settings.Apply(b.config, ss)
}

func (b *builder) init() {
	b.builders = make(map[string]ConfigBuilder)
	b.syncs = make(map[string]ConfigSyncer)
	b.registries = make(map[string]RegistryBuilder)
	b.services = make(map[string]ServiceBuilder)
	b.middlewares = make(map[string]MiddlewareBuilder)
}

func newBuilder(c *config.RuntimeConfig) *builder {
	b := &builder{
		config: c,
	}
	return b
}
