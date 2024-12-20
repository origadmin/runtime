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
		// NewConfig creates a new config using the given SourceConfig and a list of Options.
		NewConfig(*configv1.SourceConfig, ...OptionSetting) (SourceConfig, error)
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
	EnvPrefixes []string
	envVars     EnvVars
	source      SourceConfig
}

func (c *Config) LoadFromFile(path string, opts ...Option) error {
	var sources = []Source{file.NewSource(path)}
	if c.EnvPrefixes != nil {
		sources = append(sources, configenv.NewSource(c.EnvPrefixes...))
		opts = append(opts, WithSource(sources...))
	}
	c.source = New(opts...)
	return c.source.Load()
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
