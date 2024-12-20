/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/customize"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/selector"
)

type Option struct {
	EnvPrefix string
	Service   *service.Option
	Selector  *selector.Option
	Customize *customize.Option
	Source    config.KOption
}

// OptionSetting is a function that takes a pointer to a Option struct and modifies it.
type OptionSetting = func(option *Option)

func WithEnvPrefix(prefix string) OptionSetting {
	return func(s *Option) {
		s.EnvPrefix = prefix
	}
}
