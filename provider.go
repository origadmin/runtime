/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/data/storage/cache"
	"github.com/origadmin/runtime/data/storage/database"
	"github.com/origadmin/runtime/data/storage/objectstore"
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
	CategoryLogger:      log.Resolve,
	CategoryRegistrar:   registry.Resolve,
	CategoryDiscovery:   registry.Resolve,
	CategoryMiddleware:  middleware.Resolve,
	CategoryDatabase:    database.Resolve,
	CategoryCache:       cache.Resolve,
	CategoryObjectStore: objectstore.Resolve,
}
