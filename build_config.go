/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// ConfigSyncFunc is a function type that takes a SourceConfig and a list of Options and returns an error.
type ConfigSyncFunc func(*configv1.SourceConfig, any, ...config.OptionSetting) error

// SyncConfig is a method that implements the ConfigSyncer interface for ConfigSyncFunc.
func (fn ConfigSyncFunc) SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.OptionSetting) error {
	// Call the function with the given SourceConfig and a list of Options.
	return fn(cfg, v, ss...)
}

// SyncConfig is a method that implements the ConfigSyncer interface for ConfigSyncFunc.
func (b *builder) SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.OptionSetting) error {
	b.syncMux.RLock()
	defer b.syncMux.RUnlock()
	configSyncer, ok := b.syncs[cfg.Type]
	if !ok {
		return ErrNotFound
	}
	return configSyncer.SyncConfig(cfg, v, ss...)
}

func (b *builder) RegisterConfigSyncer(name string, configSyncer config.Syncer) {
	b.syncMux.Lock()
	defer b.syncMux.Unlock()
	b.syncs[name] = configSyncer
}

// RegisterConfigSync registers a new ConfigSyncer with the given name.
func (b *builder) RegisterConfigSync(name string, configSyncer ConfigSyncFunc) {
	b.RegisterConfigSyncer(name, configSyncer)
}
