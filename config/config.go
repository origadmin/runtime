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
	// Source is an option for configuring the source.
	source *SourceOption
	// Service is an option for configuring the service.
	service *ServiceOption
	// Selector is an option for selecting a service instance.
	selector *SelectorOption
	// Customize is an option for customizing the runtime.
	customize *CustomizeOption
}

var DefaultRuntimeConfig = NewRuntimeConfig()

func (r RuntimeConfig) Source() *SourceOption {
	return r.source
}

func (r RuntimeConfig) Service() *ServiceOption {
	return r.service
}

func (r RuntimeConfig) Selector() *SelectorOption {
	return r.selector
}

func (r RuntimeConfig) Customize() *CustomizeOption {
	return r.customize
}

// RuntimeConfigSetting is a type alias for a function that takes a pointer to a RuntimeConfig and modifies it.
type RuntimeConfigSetting = func(config *RuntimeConfig)

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
		source:    new(SourceOption),
		service:   new(ServiceOption),
		selector:  new(SelectorOption),
		customize: new(CustomizeOption),
	}, ss)
	return config
}
