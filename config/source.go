/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

const Type = "config"

type SourceOption struct {
	Options []Option
	Decoder Decoder
	Encoder Encoder
}

// Encoder is a function that takes a value and returns a byte slice and an error.
type Encoder func(v any) ([]byte, error)

// SourceOptionSetting is a function that takes a pointer to a SourceOption struct and modifies it.
type SourceOptionSetting = func(s *SourceOption)

// WithOptions sets the options field of the SourceOption struct.
func WithOptions(options ...Option) SourceOptionSetting {
	return func(s *SourceOption) {
		s.Options = options
	}
}
