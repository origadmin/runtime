/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package customize implements the functions, types, and interfaces for the module.
package customize

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// Setting is a struct that holds a value.
type Setting struct {
	Customize *configv1.Customize
}

// Option is a function that sets a value on a Setting.
type Option = func(s *Setting)

// WithCustomize sets the custom field of a Setting.
func WithCustomize(c *configv1.Customize) Option {
	return func(s *Setting) {
		s.Customize = c
	}
}

// GetNameConfig returns the config with the given name.
func GetNameConfig(cc *configv1.Customize, name string) *configv1.Customize_Config {
	configs := cc.GetConfigs()
	if configs != nil {
		if ret, ok := configs[name]; ok {
			return ret
		}
	}
	return nil
}

// GetTypeConfigs returns all configs with the given type.
func GetTypeConfigs(cc *configv1.Customize, typo string) map[string]*configv1.Customize_Config {
	configs := cc.GetConfigs()
	if configs == nil {
		return nil
	}
	r := make(map[string]*configv1.Customize_Config)
	for name, config := range configs {
		if config.GetType() == typo {
			r[name] = config
		}
	}
	return r
}
