/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

const Type = "config"

type SourceOption struct {
	Options   []Option
	Customize *configv1.Customize
	Decoder   Decoder
	Encoder   Encoder
}

// Encoder is a function that takes a value and returns a byte slice and an error.
type Encoder func(v any) ([]byte, error)

// SourceSetting is a function that takes a pointer to a SourceOption struct and modifies it.
type SourceSetting = func(s *SourceOption)

// WithOptions sets the options field of the SourceOption struct.
func WithOptions(options ...Option) SourceSetting {
	return func(s *SourceOption) {
		s.Options = options
	}
}

// WithCustomize sets the customize field of the SourceOption struct.
func WithCustomize(customize *configv1.Customize) SourceSetting {
	return func(s *SourceOption) {
		s.Customize = customize
	}
}
