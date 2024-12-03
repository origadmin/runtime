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
func WithSourceOption(source *SourceOption, ss ...SourceOptionSetting) RuntimeConfigSetting {
	if source == nil {
		source = new(SourceOption)
	}
	source = settings.Apply(source, ss)
	return func(config *RuntimeConfig) {
		config.source = source
	}
}

// WithServiceOption is a function that returns a RuntimeConfigSetting.
// This function sets the Service field of the RuntimeConfig.
func WithServiceOption(service *ServiceOption, ss ...ServiceOptionSetting) RuntimeConfigSetting {
	if service == nil {
		service = new(ServiceOption)
	}
	service = settings.Apply(service, ss)
	return func(config *RuntimeConfig) {
		config.service = service
	}
}

// WithSelectorOption is a function that returns a RuntimeConfigSetting.
// This function sets the Selector field of the RuntimeConfig.
func WithSelectorOption(selector *SelectorOption, ss ...SelectorOptionSetting) RuntimeConfigSetting {
	if selector == nil {
		selector = new(SelectorOption)
	}
	selector = settings.Apply(selector, ss)
	return func(config *RuntimeConfig) {
		config.selector = selector
	}
}

func WithCustomizeOption(customize *CustomizeOption, ss ...CustomizeOptionSetting) RuntimeConfigSetting {
	if customize == nil {
		customize = new(CustomizeOption)
	}
	customize = settings.Apply(customize, ss)
	return func(config *RuntimeConfig) {
		config.customize = customize
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
