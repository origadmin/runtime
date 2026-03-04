/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
)

// --- Default Resolvers (Config Mapping) ---

// DefaultConfigResolver is a no-op resolver for the root config.
var DefaultConfigResolver component.Resolver = func(source any, _ component.Category) (*component.ModuleConfig, error) {
	return nil, nil
}

// DefaultLoggerResolver resolves logger configurations.
var DefaultLoggerResolver component.Resolver = func(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.LoggerConfig); ok {
		logger := c.GetLogger()
		if logger == nil {
			return nil, nil
		}
		// Based on loggerv1 definition, if no dynamic entries exist, we treat it as a single default instance.
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "default", Value: logger}},
		}, nil
	}
	return nil, nil
}

// DefaultRegistryResolver resolves registry discoveries.
var DefaultRegistryResolver component.Resolver = func(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.RegistryConfig); ok {
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

// --- Default Providers (Component Factory) ---

// DefaultLoggerProvider creates a default logger instance.
var DefaultLoggerProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	return log.DefaultLogger, nil
}

// DefaultRegistryProvider creates a default registry instance.
var DefaultRegistryProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	return nil, nil
}

// --- Wire Providers ---

// ProvideLogger is a Wire provider function that extracts the logger from the App.
func ProvideLogger(rt *App) log.Logger {
	return rt.Logger()
}

// ProvideDefaultRegistrar is a Wire provider function that extracts the registrar from the App.
func ProvideDefaultRegistrar(rt *App) (registry.Registrar, error) {
	return rt.DefaultRegistrar()
}
