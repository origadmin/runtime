/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"fmt"
	"sync"

	configenv "github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/goexts/generic/settings"
	"google.golang.org/protobuf/proto"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/builder"
)

type (
	// Builder is an interface that defines a method for registering a config Builder.
	Builder interface {
		builder.Builder[Factory]
		Factory
	}
	// Factory is an interface that defines a method for creating a new config.
	Factory interface {
		// NewConfig creates a new config using the given KConfig and a list of Options.
		NewConfig(*configv1.SourceConfig, ...Option) (KConfig, error)
	}

	// Syncer is an interface that defines a method for synchronizing a config.
	Syncer interface {
		SyncConfig(*configv1.SourceConfig, string, any, ...Option) error
	}

	ProtoSyncer interface {
		SyncConfig(*configv1.SourceConfig, string, proto.Message, ...Option) error
	}

	ResolverRegistry interface {
		Register(name string, resolver Resolver)
		Get(name string) (Resolver, bool)
		MustGet(name string) Resolver
		Default(Resolver)
	}
)

type Config struct {
	builder     Builder
	sources     map[string]KConfig
	resolvers   ResolverRegistry
	resolvedMap map[string]Resolved
	path        string
	mu          sync.RWMutex
}

func (c *Config) Load(name string, cfg *configv1.SourceConfig, opts ...Option) error {
	if c.sources != nil {
		return nil
	}
	if _, ok := c.sources[name]; ok {
		return nil
	}

	option := settings.ApplyZero(opts)
	var sources = []KSource{
		file.NewSource(c.path),
	}
	if option.Prefixes != nil {
		sources = append(sources, configenv.NewSource(option.Prefixes...))
		option.SourceOptions = append(option.SourceOptions, WithSource(sources...))
	}
	config, err := c.builder.NewConfig(cfg, opts...)
	if err != nil {
		return err
	}
	c.sources[name] = config
	return config.Load()
}

func (c *Config) Register(name string, resolver Resolver) error {
	c.mu.RLock()
	config, ok := c.sources[name]
	c.mu.RUnlock()

	if !ok {
		return ErrNotFound
	}

	resolved, err := resolver.Resolve(config)
	if err != nil {
		return fmt.Errorf("resolve config: %w", err)
	}

	if err := config.Watch(name, resolver.Observer); err != nil {
		return fmt.Errorf("watch config: %w", err)
	}

	// 保存解析结果
	c.mu.Lock()
	if c.resolvedMap == nil {
		c.resolvedMap = make(map[string]Resolved)
	}
	c.resolvedMap[name] = resolved
	c.mu.Unlock()

	// 仍需保留注册逻辑（除非确认不再需要解析器实例）
	c.resolvers.Register(name, resolver)
	return nil
}

func (c *Config) GetResolved(name string) (Resolved, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	resolved, ok := c.resolvedMap[name]
	if !ok {
		return nil, ErrNotFound
	}
	return resolved, nil
}

func New(path string) *Config {
	return &Config{
		path:    path,
		builder: DefaultBuilder,
	}
}

func NewConfigWithBuilder(path string, builder Builder) *Config {
	return &Config{
		path:    path,
		builder: builder,
	}
}
