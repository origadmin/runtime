/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"sync"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type builder struct {
	factoryMux sync.RWMutex
	factories  map[string]Factory
}

// RegisterConfigBuilder registers a new ConfigBuilder with the given name.
func (b *builder) RegisterConfigBuilder(name string, factory Factory) {
	b.factoryMux.Lock()
	defer b.factoryMux.Unlock()
	b.factories[name] = factory
}

// RegisterConfigFunc registers a new ConfigBuilder with the given name and function.
func (b *builder) RegisterConfigFunc(name string, buildFunc BuildFunc) {
	b.RegisterConfigBuilder(name, buildFunc)
}

// BuildFunc is a function type that takes a KConfig and a list of Options and returns a Selector and an error.
type BuildFunc func(*configv1.SourceConfig, ...Option) (KConfig, error)

// NewConfig is a method that implements the ConfigBuilder interface for ConfigBuildFunc.
func (fn BuildFunc) NewConfig(cfg *configv1.SourceConfig, ss ...Option) (KConfig, error) {
	// Call the function with the given KConfig and a list of Options.
	return fn(cfg, ss...)
}

// NewConfig creates a new Selector object based on the given KConfig and options.
func (b *builder) NewConfig(cfg *configv1.SourceConfig, ss ...Option) (KConfig, error) {
	b.factoryMux.RLock()
	defer b.factoryMux.RUnlock()
	configBuilder, ok := b.factories[cfg.Type]
	if !ok {
		return nil, ErrNotFound
	}

	return configBuilder.NewConfig(cfg, ss...)
}
