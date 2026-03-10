/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"context"

	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/data/storage/objectstore"
	"github.com/origadmin/runtime/helpers/comp"
	"github.com/origadmin/runtime/helpers/configutil"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/registry"
)

func init() {
	// 1. Storage Registrations
	Register(CategoryDatabase, database.DefaultProvider)
	Register(CategoryCache, cache.DefaultProvider)
	Register(CategoryObjectStore, objectstore.DefaultProvider)

	// 2. Registry Registrations (Strictly Split)
	Register(CategoryRegistrar, registry.DefaultRegistrarProvider)
	Register(CategoryDiscovery, registry.DefaultDiscoveryProvider)

	// 3. Middleware Registrations (Scoped)
	Register(CategoryMiddleware, middleware.ServerProvider, WithScopes(ServerScope))
	Register(CategoryMiddleware, middleware.ClientProvider, WithScopes(ClientScope))
}

// --- Default Resolvers Map ---

var DefaultResolvers = map[component.Category]component.Resolver{
	CategoryLogger:      resolveLogger,
	CategoryRegistrar:   resolveRegistry,
	CategoryDiscovery:   resolveRegistry,
	CategoryMiddleware:  resolveMiddleware,
	CategoryDatabase:    resolveDatabase,
	CategoryCache:       resolveCache,
	CategoryObjectStore: resolveObjectStore,
}

// --- Specific Category Resolvers ---

func resolveLogger(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.LoggerConfig); ok {
		logger := c.GetLogger()
		if logger == nil {
			return nil, nil
		}
		// Priority: Name -> Type
		name := extractName(logger)
		if name == "" {
			name = "logger"
		}
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: name, Value: logger}},
			Active:  name,
		}, nil
	}
	return nil, nil
}

func resolveRegistry(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.DiscoveryConfig); ok {
		discoveries := c.GetDiscoveries()
		if discoveries == nil {
			return nil, nil
		}

		// Authorization flow: Default -> Active -> First
		def, configs, err := configutil.Normalize(discoveries.GetActive(), discoveries.GetDefault(), discoveries.GetConfigs())
		if err != nil {
			return nil, err
		}

		res := &component.ModuleConfig{Active: extractName(def)}
		for _, cfg := range configs {
			if name := extractName(cfg); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: cfg})
			}
		}
		return res, nil
	}
	return nil, nil
}

func resolveMiddleware(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.MiddlewareConfig); ok {
		mws := c.GetMiddlewares()
		if mws == nil {
			return nil, nil
		}
		res := &component.ModuleConfig{}
		for _, entry := range mws.GetConfigs() {
			// Priority: Name -> Type
			name := extractName(entry)
			if name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: entry})
			}
		}
		// Fallback to first if only one exists
		if len(res.Entries) == 1 {
			res.Active = res.Entries[0].Name
		}
		return res, nil
	}
	return nil, nil
}

func resolveDatabase(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.DataConfig); ok {
		data := c.GetData()
		if data == nil || data.GetDatabases() == nil {
			return nil, nil
		}
		dbs := data.GetDatabases()

		def, configs, err := configutil.Normalize(dbs.GetActive(), dbs.GetDefault(), dbs.GetConfigs())
		if err != nil {
			return nil, err
		}

		res := &component.ModuleConfig{Active: extractName(def)}
		for _, cfg := range configs {
			if name := extractName(cfg); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: cfg})
			}
		}
		return res, nil
	}
	return nil, nil
}

func resolveCache(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.DataConfig); ok {
		data := c.GetData()
		if data == nil || data.GetCaches() == nil {
			return nil, nil
		}
		caches := data.GetCaches()

		def, configs, err := configutil.Normalize(caches.GetActive(), caches.GetDefault(), caches.GetConfigs())
		if err != nil {
			return nil, err
		}

		res := &component.ModuleConfig{Active: extractName(def)}
		for _, cfg := range configs {
			if name := extractName(cfg); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: cfg})
			}
		}
		return res, nil
	}
	return nil, nil
}

func resolveObjectStore(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.DataConfig); ok {
		data := c.GetData()
		if data == nil || data.GetObjectStores() == nil {
			return nil, nil
		}
		oss := data.GetObjectStores()

		def, configs, err := configutil.Normalize(oss.GetActive(), oss.GetDefault(), oss.GetConfigs())
		if err != nil {
			return nil, err
		}

		res := &component.ModuleConfig{Active: extractName(def)}
		for _, cfg := range configs {
			if name := extractName(cfg); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: cfg})
			}
		}
		return res, nil
	}
	return nil, nil
}

func extractName(item any) string {
	if item == nil {
		return ""
	}
	// Use formalized interfaces for identification
	if n, ok := item.(component.Named); ok {
		if name := n.GetName(); name != "" {
			return name
		}
	}
	if t, ok := item.(component.Typed); ok {
		if name := t.GetType(); name != "" {
			return name
		}
	}
	if d, ok := item.(component.Dialectal); ok {
		if name := d.GetDialect(); name != "" {
			return name
		}
	}
	if d, ok := item.(component.Driver); ok {
		if name := d.GetDriver(); name != "" {
			return name
		}
	}
	return ""
}

// --- Default Providers ---

var DefaultLoggerProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	cfg, err := comp.AsConfig[loggerv1.Logger](h)
	if err != nil || cfg == nil {
		return log.DefaultLogger, nil
	}
	return log.NewLogger(cfg), nil
}

var DefaultRegistryProvider component.Provider = func(ctx context.Context, h component.Handle) (any, error) {
	return nil, nil
}
