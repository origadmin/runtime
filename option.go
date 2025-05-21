/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/context"
)

type Options struct {
	Context       context.Context
	Prefix        string
	ConfigOptions []config.Option
	Logger        log.Logger
	Signals       []os.Signal
	Resolver      config.Resolver
}

type Option func(*Options)

func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

func WithConfigOptions(options ...config.Option) Option {
	return func(o *Options) {
		o.ConfigOptions = append(o.ConfigOptions, options...)
	}
}

func WithLogger(logger log.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func WithSignals(signals ...os.Signal) Option {
	return func(o *Options) {
		o.Signals = signals
	}
}

func WithResolver(resolver config.Resolver) Option {
	return func(o *Options) {
		o.Resolver = resolver
	}
}

func WithDefaultOptions() Option {
	return func(o *Options) {
		WithLogger(log.DefaultLogger)(o)
		WithResolver(config.DefaultResolver)(o)
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
