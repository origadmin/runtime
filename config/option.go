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
	ServiceName   string
	ResolverName  string
	SourceOptions []KOption
	Decoder       KDecoder
	Encoder       Encoder
	Configure     Configure
	Prefixes      []string
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

func AppendOptions(options ...KOption) Option {
	return func(option *Options) {
		option.SourceOptions = append(option.SourceOptions, options...)
	}
}

// WithDecoderOption sets the decoder field of the KOption struct.
func WithDecoderOption(decoder KDecoder) Option {
	return func(option *Options) {
		option.Decoder = decoder
	}
}

// WithEncoderOption sets the encoder field of the KOption struct.
func WithEncoderOption(encoder Encoder) Option {
	return func(option *Options) {
		option.Encoder = encoder
	}
}

func WithConfigure(cfg Configure) Option {
	return func(option *Options) {
		option.Configure = cfg
	}
}

func WithServiceName(name string) Option {
	return func(option *Options) {
		option.ServiceName = name
	}
}

func WithPrefixes(prefixes ...string) Option {
	return func(option *Options) {
		option.Prefixes = prefixes
	}
}
