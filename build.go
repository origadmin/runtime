/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"sync"
)

// builder is a struct that holds a map of ConfigBuilders and a map of RegistryBuilders.
type builder struct {
	configMux     sync.RWMutex
	configs       map[string]ConfigBuilder
	syncMux       sync.RWMutex
	syncs         map[string]ConfigSyncer
	registryMux   sync.RWMutex
	registries    map[string]RegistryBuilder
	serviceMux    sync.RWMutex
	services      map[string]ServiceBuilder
	middlewareMux sync.RWMutex
	middlewares   map[string]MiddlewareBuilder
}

func (b *builder) init() {
	b.configs = make(map[string]ConfigBuilder)
	b.syncs = make(map[string]ConfigSyncer)
	b.registries = make(map[string]RegistryBuilder)
	b.services = make(map[string]ServiceBuilder)
	b.middlewares = make(map[string]MiddlewareBuilder)
}

func newBuilder() *builder {
	b := &builder{}
	b.init()
	return b
}
