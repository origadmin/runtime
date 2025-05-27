/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package runtime implements the functions, types, and interfaces for the module.
package runtime

import (
	"github.com/goexts/generic/settings"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/config"
)

// SyncConfig is a method that implements the ConfigSyncer interface for ConfigSyncFunc.
func (b *builder) SyncConfig(cfg *configv1.SourceConfig, v any, ss ...config.Option) error {
	b.syncMux.RLock()
	defer b.syncMux.RUnlock()
	syncers := make([]config.Syncer, 0)
	for _, tp := range cfg.Types {
		configSyncer, ok := b.syncs[tp]
		if !ok {
			continue
		}
		syncers = append(syncers, configSyncer)
	}
	options := settings.ApplyZero(ss)

	for _, s := range syncers {
		if err := s.SyncConfig(cfg, "", v, options); err != nil {
			return err
		}
	}
	return nil
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
func LoadConfig(path string, v any, ss ...config.Option) error {
	sourceConfig, err := bootstrap.LoadSourceConfig(path)
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

func LoadConfigFromBootstrap(bs *bootstrap.Bootstrap, v any, ss ...config.Option) error {
	sourceConfig, err := bootstrap.LoadSourceConfigFromBootstrap(bs)
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
