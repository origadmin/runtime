/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// CustomizeOption is a struct that holds a value.
type CustomizeOption struct {
	Customize *configv1.Customize
}

// CustomizeOptionSetting is a function that sets a value on a Setting.
type CustomizeOptionSetting = func(option *CustomizeOption)

// WithCustomizeConfig sets the customize config.
func WithCustomizeConfig(customize *configv1.Customize) CustomizeOptionSetting {
	return func(option *CustomizeOption) {
		option.Customize = customize
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
