/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configenv "github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/toolkits/env"
)

type (
	// Builder is an interface that defines a method for registering a config builder.
	Builder interface {
		Factory
		// RegisterConfigBuilder registers a config builder with the given name.
		RegisterConfigBuilder(string, Factory)
	}
	// configSyncRegistry is an interface that defines a method for synchronizing a config.
	configSyncRegistry interface {
		SyncConfig(*configv1.SourceConfig, any) error
	}

	// Factory is an interface that defines a method for creating a new config.
	Factory interface {
		// NewConfig creates a new config using the given KConfig and a list of Options.
		NewConfig(*configv1.SourceConfig, ...OptionSetting) (KConfig, error)
	}

	// Syncer is an interface that defines a method for synchronizing a config.
	Syncer interface {
		SyncConfig(*configv1.SourceConfig, any, ...OptionSetting) error
	}
)

type EnvVars struct {
	KeyValues map[string]string
}

func (v *EnvVars) Setup(prefix string) error {
	for k, v := range v.KeyValues {
		if prefix != "" {
			k = env.Var(prefix, k)
		}
		err := env.SetEnv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *EnvVars) Set(key string, value string) {
	v.KeyValues[key] = value
}

func (v *EnvVars) Get(key string) (string, bool) {
	if v, ok := v.KeyValues[key]; ok {
		return v, true
	}
	return "", false
}

type Config struct {
	cfg         any
	envVars     EnvVars
	source      KConfig
	Path        string
	EnvPrefixes []string
	Builder     Builder
	registry    func(source any, serviceName string) (*configv1.Registry, error)
	service     func(source any, serviceName string) (*configv1.Service, error)
}

func (c *Config) LoadFromFile(path string, opts ...KOption) error {
	if c.source != nil {
		return nil
	}
	var sources = []KSource{file.NewSource(path)}
	if c.EnvPrefixes != nil {
		sources = append(sources, configenv.NewSource(c.EnvPrefixes...))
		opts = append(opts, WithSource(sources...))
	}
	c.source = NewSourceConfig(opts...)
	return c.source.Load()
}

func (c *Config) LoadFromSource(cfg *configv1.SourceConfig, opts ...OptionSetting) error {
	if c.source != nil {
		return nil
	}

	config, err := c.Builder.NewConfig(cfg, opts...)
	if err != nil {
		return err
	}
	c.source = config
	return c.source.Load()
}

func (c *Config) Scan() error {
	return c.source.Scan(c.cfg)
}

func (c *Config) Bind(cfg any) error {
	c.cfg = cfg
	return c.Scan()
}

func (c *Config) Watch(key string, ob KObserver) error {
	return c.source.Watch(key, ob)
}

func (c *Config) SetEnv(key, value string) {
	c.envVars.Set(key, value)
}

func (c *Config) GetEnv(key string) (string, bool) {
	if v, ok := c.envVars.Get(key); ok {
		return v, true
	}
	return "", false
}

func (c *Config) Setup(prefix string) error {
	return c.envVars.Setup(prefix)
}

func (c *Config) BindRegistry(fn func(source any, serviceName string) (*configv1.Registry, error)) {
	c.registry = fn
}

func (c *Config) Registry(serviceName string) (*configv1.Registry, error) {
	if c.registry != nil {
		return c.registry(c.cfg, serviceName)
	}
	return nil, ErrNotFound
}

func (c *Config) Service(serviceName string) (*configv1.Service, error) {
	if c.service != nil {
		return c.service(c.cfg, serviceName)
	}
	return nil, ErrNotFound
}

func NewBuilder() Builder {
	return &builder{
		factories: make(map[string]Factory),
	}
}
