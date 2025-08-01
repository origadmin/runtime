/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"fmt"
	"sync"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
)

type Loader struct {
	cfg          *configv1.SourceConfig
	builder      Builder
	sourceConfig KConfig
	resolver     Resolver
	resolved     Resolved
	mu           sync.RWMutex
}

func (c *Loader) Load(cfg *configv1.SourceConfig, opts ...Option) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sourceConfig != nil {
		return nil
	}
	if c.resolver == nil {
		return fmt.Errorf("resolver is not set")
	}
	config, err := c.builder.NewConfig(cfg, opts...)
	if err != nil {
		return err
	}
	if err := config.Load(); err != nil {
		return err
	}
	resolved, err := c.resolver.Resolve(config) // Resolve without internal lock
	if err != nil {
		return err
	}
	c.sourceConfig = config
	c.cfg = cfg
	c.resolved = resolved

	return nil
}

func (c *Loader) Resolve(config KConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.resolver == nil {
		return fmt.Errorf("resolver is not set")
	}
	resolved, err := c.resolver.Resolve(config)
	if err != nil {
		return err
	}
	c.resolved = resolved
	return nil
}

func (c *Loader) GetConfig() (*configv1.SourceConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cfg, nil
}

func (c *Loader) GetSource() (KConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sourceConfig, nil
}

func (c *Loader) GetResolved() (Resolved, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.resolved, nil
}

func (c *Loader) GetResolver() (Resolver, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.resolver, nil
}

func (c *Loader) SetResolver(resolver Resolver) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.resolver = resolver
	return nil
}

func New() *Loader {
	return NewWithBuilder(DefaultBuilder)
}

func NewWithBuilder(builder Builder) *Loader {
	return &Loader{
		builder: builder,
	}
}
