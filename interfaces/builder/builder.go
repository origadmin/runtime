/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package builder implements the functions, types, and interfaces for the module.
package builder

import (
	"maps"
	"sync"
)

type Builder[F any] interface {
	Get(name string) (F, bool)
	Register(name string, factory F)
	RegisteredFactories() map[string]F
}

type builder[F any] struct {
	factories map[string]F
	mu        sync.RWMutex
}

func (b *builder[F]) Get(name string) (F, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	factory, ok := b.factories[name]
	return factory, ok
}

func (b *builder[F]) Register(name string, factory F) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.factories[name] = factory
}

func (b *builder[F]) RegisteredFactories() map[string]F {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return maps.Clone(b.factories)
}

func New[F any]() Builder[F] {
	return &builder[F]{
		factories: make(map[string]F),
	}
}
