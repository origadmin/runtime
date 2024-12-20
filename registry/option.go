/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package registry implements the functions, types, and interfaces for the module.
package registry

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type Configure interface {
	FromConfig(cfg **configv1.Registry) error
}
type Option struct {
	Configure Configure
}

type OptionSetting = func(o *Option)

func WithConfigure(cfg Configure) OptionSetting {
	return func(o *Option) {
		o.Configure = cfg
	}
}
