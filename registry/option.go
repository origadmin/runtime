/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	"time"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Configure interface {
	Type() string
	FromConfig(cfg *configv1.Registry) error
}
type Options struct {
	Configure Configure
	Timeout   time.Duration
	Retries   int
}

type Option = func(o *Options)

func WithConfigure(cfg Configure) Option {
	return func(o *Options) {
		o.Configure = cfg
	}
}

func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

func WithRetries(n int) Option {
	return func(o *Options) {
		o.Retries = n
	}
}
