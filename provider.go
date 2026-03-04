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

// --- Default Global Resolver (The Dispatcher) ---

// DefaultGlobalResolver is the primary config dispatcher for the framework.
var DefaultGlobalResolver component.Resolver = func(source any, cat component.Category) (*component.ModuleConfig, error) {
	// Centralized dispatching based on category.
	// This reduces redundant assertions across multiple extractors.
	switch cat {
	case CategoryLogger:
		return resolveLogger(source)
	case CategoryRegistry:
		return resolveRegistry(source)
	default:
		// Unknown categories are ignored by the default resolver.
		// They can still be handled by local resolvers (WithResolver option).
		return nil, nil
	}
}

// resolveLogger handles extraction for Logger components.
func resolveLogger(source any) (*component.ModuleConfig, error) {
	if c, ok := source.(component.LoggerConfig); ok {
		logger := c.GetLogger()
		if logger == nil {
			return nil, nil
		}
		// Treating the single logger segment as the default instance.
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "default", Value: logger}},
		}, nil
	}
	return nil, nil
}

// resolveRegistry handles extraction for Registry/Discovery components.
func resolveRegistry(source any) (*component.ModuleConfig, error) {
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
