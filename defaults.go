/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"context"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
)

var (
	// DefaultConfigExtractor extracts the configuration from the config sources.
	DefaultConfigExtractor component.Extractor = func(root any) (*component.ModuleConfig, error) {
		return nil, nil
	}

	// DefaultLoggerExtractor extracts the logger configuration from the config sources.
	DefaultLoggerExtractor component.Extractor = func(root any) (*component.ModuleConfig, error) {
		return nil, nil
	}

	// DefaultRegistryExtractor extracts the registry configuration from the config sources.
	DefaultRegistryExtractor component.Extractor = func(root any) (*component.ModuleConfig, error) {
		if c, ok := root.(component.RegistryConfig); ok {
			discoveries := c.GetDiscoveries()
			if discoveries == nil {
				return nil, nil
			}
			res := &component.ModuleConfig{Active: discoveries.GetActive()}
			for _, entry := range discoveries.GetConfigs() {
				name := entry.GetName()
				if name == "" {
					name = entry.GetType()
				}
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: entry})
			}
			return res, nil
		}
		return nil, nil
	}

	DefaultLoggerProvider   component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) { return nil, nil }
	DefaultRegistryProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) { return nil, nil }
)
