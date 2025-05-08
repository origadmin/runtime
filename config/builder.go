/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/builder"
)

var (
	DefaultBuilder = NewBuilder()
)

type buildImpl struct {
	builder.Builder[Factory]
}

// RegisterConfigFunc registers a new ConfigBuilder with the given name and function.
func (b *buildImpl) RegisterConfigFunc(name string, buildFunc BuildFunc) {
	b.Register(name, buildFunc)
}

// BuildFunc is a function type that takes a KConfig and a list of Options and returns a Selector and an error.
type BuildFunc func(*configv1.SourceConfig, ...Option) (KConfig, error)

// NewConfig is a method that implements the ConfigBuilder interface for ConfigBuildFunc.
func (fn BuildFunc) NewConfig(cfg *configv1.SourceConfig, ss ...Option) (KConfig, error) {
	// Call the function with the given KConfig and a list of Options.
	return fn(cfg, ss...)
}

// NewConfig creates a new Selector object based on the given KConfig and options.
func (b *buildImpl) NewConfig(cfg *configv1.SourceConfig, ss ...Option) (KConfig, error) {
	configBuilder, ok := b.Get(cfg.Type)
	if !ok {
		return nil, ErrNotFound
	}

	return configBuilder.NewConfig(cfg, ss...)
}

func NewBuilder() Builder {
	return &buildImpl{
		Builder: builder.New[Factory](),
	}
}
