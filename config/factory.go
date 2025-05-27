/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"fmt"

	configenv "github.com/go-kratos/kratos/v2/config/env"
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/log"
)

var (
	DefaultBuilder = NewBuilder()
)

type configBuilder struct {
	factory.Registry[Factory]
}

// RegisterConfigFunc registers a new ConfigBuilder with the given name and function.
func (b *configBuilder) RegisterConfigFunc(name string, buildFunc BuildFunc) {
	b.Register(name, buildFunc)
}

// BuildFunc is a function type that takes a KConfig and a list of Options and returns a Selector and an error.
type BuildFunc func(*configv1.SourceConfig, *Options) (KSource, error)

// NewSource is a method that implements the ConfigBuilder interface for ConfigBuildFunc.
func (fn BuildFunc) NewSource(cfg *configv1.SourceConfig, opts *Options) (KSource, error) {
	// Call the function with the given KConfig and a list of Options.
	return fn(cfg, opts)
}

// NewConfig creates a new Selector object based on the given KConfig and options.
func (b *configBuilder) NewConfig(cfg *configv1.SourceConfig, opts ...Option) (KConfig, error) {
	options := settings.ApplyZero(opts)
	sources := options.Sources
	if sources == nil {
		sources = make([]KSource, 0)
	}

	for _, t := range cfg.Types {
		b, ok := b.Get(t)
		if !ok {
			return nil, fmt.Errorf("unknown type: %s", t)
		}
		log.Infof("registering type: %s", t)
		source, err := b.NewSource(cfg, options)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	if options.EnvPrefixes != nil {
		sources = append(sources, configenv.NewSource(options.EnvPrefixes...))
	}

	options.ConfigOptions = append(options.ConfigOptions, WithSource(sources...))
	return NewSourceConfig(options.ConfigOptions...), nil
}

func NewBuilder() Builder {
	return &configBuilder{
		Registry: factory.New[Factory](),
	}
}
