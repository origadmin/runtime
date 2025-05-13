/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package factory implements the functions, types, and interfaces for the module.
package factory

import (
	"maps"
	"sync"
)

// Registry defines the interface for managing factory functions.
// It allows for registering and retrieving factory functions by name.
type Registry[F any] interface {
	Get(name string) (F, bool)
	Register(name string, factory F)
	RegisteredFactories() map[string]F
}

// registry is a thread-safe implementation of the Registry interface.
// It manages the lifecycle and access of factory functions.
type registry[F any] struct {
	factories map[string]F
	mu        sync.RWMutex
}

// Get retrieves a factory function by name.
// It returns the factory function and a boolean indicating its presence.
func (f *registry[F]) Get(name string) (F, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	factory, ok := f.factories[name]
	return factory, ok
}

// Register adds a new factory function to the registry.
// If a factory with the same name already exists, it will be overwritten.
func (f *registry[F]) Register(name string, factory F) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.factories[name] = factory
}

// RegisteredFactories returns a copy of the current factory function map.
// This prevents modification of the registry's internal state from the outside.
func (f *registry[F]) RegisteredFactories() map[string]F {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return maps.Clone(f.factories)
}

// New creates and returns a new instance of the registry.
// This is the entry point for creating a registry for factory functions.
func New[F any]() Registry[F] {
	return &registry[F]{
		factories: make(map[string]F),
	}
}
