/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"fmt"
	"sync"

	configenv "github.com/go-kratos/kratos/v2/config/env"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Loader struct {
	builder  Builder
	source   KConfig
	resolver Resolver
	resolved Resolved
	mu       sync.RWMutex
}

func (c *Loader) Load(cfg *configv1.SourceConfig, opts ...Option) error {
	options := settings.ApplyZero(opts)
	if c.source != nil && !options.ForceReload {
		return nil
	}

	if err != nil {
		return err
	}

	if err := config.Load(); err != nil {
		return err
	}

	if resolver := options.Resolver; resolver != nil {
		err := c.Resolve(config, resolver)
		if err != nil {
			return err
		}
	}
	c.source = config
	return nil
}

func (c *Loader) Resolve(config KConfig, resolver Resolver) error {
	resolved, err := resolver.Resolve(config)
	if err != nil {
		return fmt.Errorf("resolve config: %w", err)
	}

	c.mu.Lock()
	c.resolved = resolved
	c.resolver = resolver
	c.mu.Unlock()
	return nil
}

func (c *Loader) GetSource() (KConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.source, nil
}

func (c *Loader) GetResolved() (Resolved, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.resolved, nil
}

func New() *Loader {
	return &Loader{
		builder: DefaultBuilder,
	}
}

func NewWithBuilder(builder Builder) *Loader {
	return &Loader{
		builder: builder,
	}
}
