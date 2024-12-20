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

type Option struct {
	SourceOptions []KOption
	Decoder       KDecoder
	Encoder       Encoder
	Configure     Configure
}

// Encoder is a function that takes a value and returns a byte slice and an error.
type Encoder func(v any) ([]byte, error)

// OptionSetting is a function that takes a pointer to a KOption struct and modifies it.
type OptionSetting = func(s *Option)

// WithOptions sets the options field of the KOption struct.
func WithOptions(options ...KOption) OptionSetting {
	return func(option *Option) {
		option.SourceOptions = options
	}
}

func WithConfigure(cfg Configure) OptionSetting {
	return func(option *Option) {
		option.Configure = cfg
	}
}
