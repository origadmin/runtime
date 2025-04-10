/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Configure interface {
	FromConfig(cfg *configv1.SourceConfig) error
}

type Options struct {
	SourceOptions []KOption
	Decoder       KDecoder
	Encoder       Encoder
	Configure     Configure
}

// Encoder is a function that takes a value and returns a byte slice and an error.
type Encoder func(v any) ([]byte, error)

// Option is a function that takes a pointer to a KOption struct and modifies it.
type Option = func(s *Options)

// WithOptions sets the options field of the KOption struct.
func WithOptions(options ...KOption) Option {
	return func(option *Options) {
		option.SourceOptions = options
	}
}

func WithConfigure(cfg Configure) Option {
	return func(option *Options) {
		option.Configure = cfg
	}
}
