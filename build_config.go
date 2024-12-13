/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

type (
	// configBuildRegistry is an interface that defines a method for registering a config builder.
	configBuildRegistry interface {
		// RegisterConfigBuilder registers a config builder with the given name.
		RegisterConfigBuilder(string, ConfigBuilder)
	}
	// configSyncRegistry is an interface that defines a method for synchronizing a config.
	configSyncRegistry interface {
		SyncConfig(*configv1.SourceConfig, any) error
	}

	// ConfigBuilder is an interface that defines a method for creating a new config.
	ConfigBuilder interface {
		// NewConfig creates a new config using the given SourceConfig and a list of Options.
		NewConfig(*configv1.SourceConfig, ...config.SourceOptionSetting) (config.Config, error)
	}

	// ConfigSyncer is an interface that defines a method for synchronizing a config.
	ConfigSyncer interface {
		SyncConfig(*configv1.SourceConfig, any, ...config.SourceOptionSetting) error
	}
)

// ConfigBuildFunc is a function type that takes a SourceConfig and a list of Options and returns a Selector and an error.
type ConfigBuildFunc func(*configv1.SourceConfig, ...config.SourceOptionSetting) (config.Config, error)

// NewConfig is a method that implements the ConfigBuilder interface for ConfigBuildFunc.
func (fn ConfigBuildFunc) NewConfig(cfg *configv1.SourceConfig, ss ...config.SourceOptionSetting) (config.Config, error) {
	// Call the function with the given SourceConfig and a list of Options.
	return fn(cfg, ss...)
}

// NewConfig creates a new Selector object based on the given SourceConfig and options.
func (b *builder) NewConfig(cfg *configv1.SourceConfig, ss ...config.SourceOptionSetting) (config.Config, error) {
	b.configMux.RLock()
	defer b.configMux.RUnlock()
	configBuilder, ok := b.configs[cfg.Type]
	if !ok {
		return nil, ErrNotFound
	}

	return configBuilder.NewConfig(cfg, ss...)
}

// ConfigSyncFunc is a function type that takes a SourceConfig and a list of Options and returns an error.
type ConfigSyncFunc func(*configv1.SourceConfig, any, ...config.SourceOptionSetting) error

// SyncConfig is a method that implements the ConfigSyncer interface for ConfigSyncFunc.
func (fn ConfigSyncFunc) SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.SourceOptionSetting) error {
	// Call the function with the given SourceConfig and a list of Options.
	return fn(cfg, v, ss...)
}

// SyncConfig is a method that implements the ConfigSyncer interface for ConfigSyncFunc.
func (b *builder) SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.SourceOptionSetting) error {
	b.syncMux.RLock()
	defer b.syncMux.RUnlock()
	configSyncer, ok := b.syncs[cfg.Type]
	if !ok {
		return ErrNotFound
	}
	return configSyncer.SyncConfig(cfg, v, ss...)
}

// RegisterConfigBuilder registers a new ConfigBuilder with the given name.
func (b *builder) RegisterConfigBuilder(name string, configBuilder ConfigBuilder) {
	b.configMux.Lock()
	defer b.configMux.Unlock()
	b.configs[name] = configBuilder
}

// RegisterConfigFunc registers a new ConfigBuilder with the given name and function.
func (b *builder) RegisterConfigFunc(name string, configBuilder ConfigBuildFunc) {
	b.RegisterConfigBuilder(name, configBuilder)
}

func (b *builder) RegisterConfigSyncer(name string, configSyncer ConfigSyncer) {
	b.configMux.Lock()
	defer b.configMux.Unlock()
	b.syncs[name] = configSyncer
}

// RegisterConfigSync registers a new ConfigSyncer with the given name.
func (b *builder) RegisterConfigSync(name string, configSyncer ConfigSyncFunc) {
	b.RegisterConfigSyncer(name, configSyncer)
}
