/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/config"
)

type Options struct {
	Prefix        string
	ConfigOptions []config.Option
	Logger        log.Logger
	Signals       []os.Signal
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
