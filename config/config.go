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
	"google.golang.org/protobuf/proto"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
)

type (
	// Builder is an interface that defines a method for registering a config Builder.
	Builder interface {
		factory.Registry[Factory]
		NewConfig(*configv1.SourceConfig, ...Option) (KConfig, error)
	}
	// Factory is an interface that defines a method for creating a new config.
	Factory interface {
		// NewSource creates a new config using the given KConfig and a list of Options.
		NewSource(*configv1.SourceConfig, ...Option) (KSource, error)
	}

	// Syncer is an interface that defines a method for synchronizing a config.
	Syncer interface {
		SyncConfig(*configv1.SourceConfig, string, any, ...Option) error
	}

	// ProtoSyncer is an interface that defines a method for synchronizing a protobuf message.
	ProtoSyncer interface {
		SyncConfig(*configv1.SourceConfig, string, proto.Message, ...Option) error
	}
)

type Config struct {
	builder  Builder
	source   KConfig
	resolver Resolver
	resolved Resolved
	mu       sync.RWMutex
}

func (c *Config) Load(cfg *configv1.SourceConfig, opts ...Option) error {
	option := settings.ApplyZero(opts)
	if c.source != nil && !option.ForceReload {
		return nil
	}
	sources := option.Sources
	if sources == nil {
		sources = make([]KSource, 0)
	}

	if option.EnvPrefixes != nil {
		sources = append(sources, configenv.NewSource(option.EnvPrefixes...))
	}
	option.ConfigOptions = append(option.ConfigOptions, WithSource(sources...))

	config, err := c.builder.NewConfig(cfg, opts...)
	if err != nil {
		return err
	}

	if err := config.Load(); err != nil {
		return err
	}

	if resolver := option.Resolver; resolver != nil {
		err := c.Resolve(config, resolver)
		if err != nil {
			return err
		}
	}
	c.source = config
	return nil
}

func (c *Config) Resolve(config KConfig, resolver Resolver) error {
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

func (c *Config) GetSource() (KConfig, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.source, nil
}

func (c *Config) GetResolved() (Resolved, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.resolved, nil
}

func New() *Config {
	return &Config{
		builder: DefaultBuilder,
	}
}

func NewWithBuilder(builder Builder) *Config {
	return &Config{
		builder: builder,
	}
}
