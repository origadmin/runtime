/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package config implements the functions, types, and interfaces for the module.
package config

import (
	"github.com/goexts/generic/settings"
)

// RuntimeConfig is a struct that holds the configuration for the runtime.
type RuntimeConfig struct {
	// Bootstrap is an option for configuring the bootstrap.
	bootstrap *BootstrapOption
	// Source is an option for configuring the source.
	source *SourceOption
	// Service is an option for configuring the service.
	service *ServiceOption
	// Selector is an option for selecting a service instance.
	selector *SelectorOption
	// Customize is an option for customizing the runtime.
	customize *CustomizeOption
}

// RuntimeConfigSetting is a type alias for a function that takes a pointer to a RuntimeConfig and modifies it.
type RuntimeConfigSetting = func(config *RuntimeConfig)

// DefaultRuntimeConfig is a pre-initialized RuntimeConfig instance.
var DefaultRuntimeConfig = NewRuntimeConfig()

// Bootstrap returns the BootstrapOption associated with the RuntimeConfig.
func (r RuntimeConfig) Bootstrap() *BootstrapOption {
	// Return the bootstrap option.
	return r.bootstrap
}

// Source returns the SourceOption associated with the RuntimeConfig.
func (r RuntimeConfig) Source() *SourceOption {
	// Return the source option.
	return r.source
}

// Service returns the ServiceOption associated with the RuntimeConfig.
func (r RuntimeConfig) Service() *ServiceOption {
	// Return the service option.
	return r.service
}

// Selector returns the SelectorOption associated with the RuntimeConfig.
func (r RuntimeConfig) Selector() *SelectorOption {
	// Return the selector option.
	return r.selector
}

// Customize returns the CustomizeOption associated with the RuntimeConfig.
func (r RuntimeConfig) Customize() *CustomizeOption {
	// Return the customize option.
	return r.customize
}

// WithSourceOption is a function that returns a RuntimeConfigSetting.
// This function sets the Source field of the RuntimeConfig.
func WithSourceOption(ss ...SourceOptionSetting) RuntimeConfigSetting {
	return func(config *RuntimeConfig) {
		if config.source == nil {
			config.source = new(SourceOption)
		}
		config.source = settings.Apply(config.source, ss)
	}
}

// WithServiceOption is a function that returns a RuntimeConfigSetting.
// This function sets the Service field of the RuntimeConfig.
func WithServiceOption(ss ...ServiceOptionSetting) RuntimeConfigSetting {
	return func(config *RuntimeConfig) {
		if config.service == nil {
			config.service = new(ServiceOption)
		}
		config.service = settings.Apply(config.service, ss)
	}
}

// WithSelectorOption is a function that returns a RuntimeConfigSetting.
// This function sets the Selector field of the RuntimeConfig.
func WithSelectorOption(ss ...SelectorOptionSetting) RuntimeConfigSetting {
	return func(config *RuntimeConfig) {
		if config.selector == nil {
			config.selector = new(SelectorOption)
		}
		config.selector = settings.Apply(config.selector, ss)
	}
}

func WithCustomizeOption(ss ...CustomizeOptionSetting) RuntimeConfigSetting {
	return func(config *RuntimeConfig) {
		if config.customize == nil {
			config.customize = new(CustomizeOption)
		}
		config.customize = settings.Apply(config.customize, ss)
	}
}

// NewRuntimeConfig is a function that creates a new RuntimeConfig.
// It takes a variadic list of RuntimeConfigSettings and applies them to a new RuntimeConfig.
func NewRuntimeConfig(ss ...RuntimeConfigSetting) *RuntimeConfig {
	config := settings.Apply(&RuntimeConfig{
		bootstrap: &BootstrapOption{
			EnvPrefix: EnvPrefix,
		},
		source:    new(SourceOption),
		service:   new(ServiceOption),
		selector:  new(SelectorOption),
		customize: new(CustomizeOption),
	}, ss)
	return config
}
