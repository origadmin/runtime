/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/data/storage/objectstore"
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
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "default", Value: logger}},
			Active:  "default",
		}, nil
	}
	return nil, nil
}

func resolveRegistry(source any, _ component.Category) (*component.ModuleConfig, error) {
	if c, ok := source.(component.RegistryConfig); ok {
		discoveries := c.GetDiscoveries()
		if discoveries == nil {
			return nil, nil
		}
		res := &component.ModuleConfig{Active: discoveries.GetActive()}
		if d := discoveries.GetDefault(); d != nil {
			name := extractName(d)
			res.Entries = append(res.Entries, component.ConfigEntry{Name: "default", Value: d})
			if name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: d})
			}
			if res.Active == "" {
				res.Active = "default"
			}
		}
		for _, entry := range discoveries.GetConfigs() {
			if name := extractName(entry); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: entry})
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
			if name := extractName(entry); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: entry})
			}
		}
		if len(res.Entries) == 1 {
			res.Entries = append(res.Entries, component.ConfigEntry{Name: "default", Value: res.Entries[0].Value})
			res.Active = "default"
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
		res := &component.ModuleConfig{Active: dbs.GetActive()}
		if d := dbs.GetDefault(); d != nil {
			name := extractName(d)
			res.Entries = append(res.Entries, component.ConfigEntry{Name: "default", Value: d})
			if name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: d})
			}
			if res.Active == "" {
				res.Active = "default"
			}
		}
		for _, entry := range dbs.GetConfigs() {
			if name := extractName(entry); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: entry})
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
		res := &component.ModuleConfig{Active: caches.GetActive()}
		if d := caches.GetDefault(); d != nil {
			name := extractName(d)
			res.Entries = append(res.Entries, component.ConfigEntry{Name: "default", Value: d})
			if name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: d})
			}
			if res.Active == "" {
				res.Active = "default"
			}
		}
		for _, entry := range caches.GetConfigs() {
			if name := extractName(entry); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: entry})
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
		res := &component.ModuleConfig{Active: oss.GetActive()}
		if d := oss.GetDefault(); d != nil {
			name := extractName(d)
			res.Entries = append(res.Entries, component.ConfigEntry{Name: "default", Value: d})
			if name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: d})
			}
			if res.Active == "" {
				res.Active = "default"
			}
		}
		for _, entry := range oss.GetConfigs() {
			if name := extractName(entry); name != "" {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: name, Value: entry})
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
	if n, ok := item.(interface{ GetName() string }); ok {
		if name := n.GetName(); name != "" {
			return name
		}
	}
	if d, ok := item.(interface{ GetDialect() string }); ok {
		if name := d.GetDialect(); name != "" {
			return name
		}
	}
	if t, ok := item.(interface{ GetType() string }); ok {
		if name := t.GetType(); name != "" {
			return name
		}
	}
	if d, ok := item.(interface{ GetDriver() string }); ok {
		if name := d.GetDriver(); name != "" {
			return name
		}
	}
	return ""
}

// --- Default Providers ---

var DefaultLoggerProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	return log.DefaultLogger, nil
}

var DefaultRegistryProvider component.Provider = func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
	return nil, nil
}
