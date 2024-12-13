/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package customize implements the functions, types, and interfaces for the module.
package customize

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// ConfigFromName returns the config with the given name.
func ConfigFromName(cc *configv1.Customize, name string) *configv1.Customize_Config {
	configs := cc.GetConfigs()
	if configs != nil {
		if ret, ok := configs[name]; ok {
			return ret
		}
	}
	return nil
}

// ConfigsFromType returns all configs with the given type.
func ConfigsFromType(cc *configv1.Customize, typo string) map[string]*configv1.Customize_Config {
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
