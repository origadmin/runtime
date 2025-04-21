/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/customize"
	"github.com/origadmin/runtime/service/selector"
)

type Options struct {
	EnvPrefix string
	Selector  *selector.Options
	Customize *customize.Options
	Source    config.KOption
}

// Option is a function that takes a pointer to a Options struct and modifies it.
type Option = func(option *Options)

func WithEnvPrefix(prefix string) Option {
	return func(s *Options) {
		s.EnvPrefix = prefix
	}
}
