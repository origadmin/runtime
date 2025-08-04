/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/origadmin/runtime/interfaces"
)

type Options struct {
	ConfigName    string
	ServiceName   string
	ResolverName  string
	EnvPrefixes   []string
	Sources       []kratosconfig.Source
	ConfigOptions []kratosconfig.Option
	Decoder       kratosconfig.Decoder
	Encoder       interfaces.Encoder
	ForceReload   bool
}

// Encoder is a function that takes a value and returns a byte slice and an error.
type Encoder func(v any) ([]byte, error)

// Option is a function that takes a pointer to a KOption struct and modifies it.
type Option = func(s *Options)

// WithOptions sets the options field of the KOption struct.
func WithOptions(options ...kratosconfig.Option) Option {
	return func(option *Options) {
		option.ConfigOptions = options
	}
}

func AppendOptions(options ...kratosconfig.Option) Option {
	return func(option *Options) {
		option.ConfigOptions = append(option.ConfigOptions, options...)
	}
}

// WithDecoderOption sets the decoder field of the KOption struct.
func WithDecoderOption(decoder kratosconfig.Decoder) Option {
	return func(option *Options) {
		option.Decoder = decoder
	}
}

// WithEncoderOption sets the encoder field of the KOption struct.
func WithEncoderOption(encoder interfaces.Encoder) Option {
	return func(option *Options) {
		option.Encoder = encoder
	}
}

func WithServiceName(name string) Option {
	return func(option *Options) {
		option.ServiceName = name
	}
}

func WithEnvPrefixes(prefixes ...string) Option {
	return func(option *Options) {
		option.EnvPrefixes = prefixes
	}
}

func WithForceReload() Option {
	return func(o *Options) {
		o.ForceReload = true
	}
}

func WithConfigName(name string) Option {
	return func(option *Options) {
		option.ConfigName = name
	}
}
