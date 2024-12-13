/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package customize implements the functions, types, and interfaces for the module.
package customize

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// Config is a struct that holds a value.
type Config struct {
	Customize *configv1.Customize
}

// ConfigSetting is a function that sets a value on a Setting.
type ConfigSetting = func(config *Config)

// WithCustomizeConfig sets the customize config.
func WithCustomizeConfig(customize *configv1.Customize) ConfigSetting {
	return func(option *Config) {
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
