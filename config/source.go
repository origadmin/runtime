/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

type Option struct {
	SourceOptions []SourceOption
	Decoder       Decoder
	Encoder       Encoder
}

// Encoder is a function that takes a value and returns a byte slice and an error.
type Encoder func(v any) ([]byte, error)

// OptionSetting is a function that takes a pointer to a SourceOption struct and modifies it.
type OptionSetting = func(s *Option)

// WithOptions sets the options field of the SourceOption struct.
func WithOptions(options ...SourceOption) OptionSetting {
	return func(option *Option) {
		option.SourceOptions = options
	}
}
