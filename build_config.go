/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/config"
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

// SyncConfig is a method that implements the ConfigSyncer interface for ConfigSyncFunc.
func (b *builder) SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.Option) error {
	b.syncMux.RLock()
	defer b.syncMux.RUnlock()
	configSyncer, ok := b.syncs[cfg.Type]
	if !ok {
		return ErrNotFound
	}
	return configSyncer.SyncConfig(cfg, "", v, ss...)
}

func (b *builder) RegisterConfigSyncer(name string, configSyncer config.Syncer) {
	b.syncMux.Lock()
	defer b.syncMux.Unlock()
	b.syncs[name] = configSyncer
}

// RegisterConfigSync registers a new ConfigSyncer with the given name.
func (b *builder) RegisterConfigSync(name string, configSyncer config.Syncer) {
	b.RegisterConfigSyncer(name, configSyncer)
}

// LoadConfig loads the config file from the given path
func LoadConfig(bs *bootstrap.Bootstrap, v any, ss ...config.Option) error {
	sourceConfig, err := bootstrap.LoadSourceConfig(bs)
	if err != nil {
		return err
	}
	runtimeConfig, err := NewConfig(sourceConfig, ss...)
	if err != nil {
		return err
	}
	if err := runtimeConfig.Load(); err != nil {
		return err
	}
	if err := runtimeConfig.Scan(v); err != nil {
		return err
	}
	return nil
}
