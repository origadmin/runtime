/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

type BootstrapOption struct {
	EnvPrefix string
}

// BootstrapOptionSetting is a function that takes a pointer to a BootstrapOption struct and modifies it.
type BootstrapOptionSetting = func(s *BootstrapOption)

func WithBootstrapEnvPrefix(prefix string) BootstrapOptionSetting {
	return func(s *BootstrapOption) {
		s.EnvPrefix = prefix
	}
}
