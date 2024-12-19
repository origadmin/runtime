/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

type Config struct {
	envs        map[string]string
	EnvPrefixes []string
	Source      SourceConfig
}

func (c *Config) LoadFromFile(path string, opts ...Option) error {
	var sources = []Source{file.NewSource(path)}
	if c.EnvPrefixes != nil {
		sources = append(sources, env.NewSource(c.EnvPrefixes...))
		opts = append(opts, WithSource(sources...))
	}
	c.Source = New(opts...)
	return c.Source.Load()
}

func (c *Config) SetEnv(key, value string) {
	c.envs[key] = value
}

func (c *Config) GetEnv(key string) (string, bool) {
	if v, ok := c.envs[key]; ok {
		return v, true
	}
	return "", false
}
