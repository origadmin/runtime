/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configenv "github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/goexts/generic/settings"

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
		NewConfig(*configv1.SourceConfig, ...Option) (KConfig, error)
	}

	// Syncer is an interface that defines a method for synchronizing a config.
	Syncer interface {
		SyncConfig(*configv1.SourceConfig, any, ...Option) error
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
	cfg         *configv1.SourceConfig
	envVars     EnvVars
	sources     map[string]KConfig
	EnvPrefixes []string
	builder     Builder
	path        string
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
	if c.EnvPrefixes != nil {
		sources = append(sources, configenv.NewSource(c.EnvPrefixes...))
		option.SourceOptions = append(option.SourceOptions, WithSource(sources...))
	}
	config, err := c.builder.NewConfig(cfg, opts...)
	if err != nil {
		return err
	}
	c.sources[name] = config
	return config.Load()
}

func (c *Config) Scan(name string, v any) error {
	service, ok := c.sources[name]
	if !ok {
		return ErrNotFound
	}
	return service.Scan(v)
}

func (c *Config) Watch(key string, ob KObserver) error {
	service, ok := c.sources[key]
	if !ok {
		return ErrNotFound
	}
	return service.Watch(key, ob)
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

func NewBuilder() Builder {
	return &builder{
		factories: make(map[string]Factory),
	}
}

func New(path string, builder Builder) Config {
	return Config{
		path: path,
		envVars: EnvVars{
			KeyValues: make(map[string]string),
		},
		builder: builder,
	}
}
