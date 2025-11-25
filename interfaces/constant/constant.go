/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package constant defines the constant keys for structured configuration components.
// These keys serve as a contract between the configuration interface and its implementations,
// ensuring consistent access to different parts of the configuration.
package constant

// ComponentKey defines the type for configuration component keys.
type ComponentKey string

const (
	// ConfigApp is the key for the application's core configuration.
	ConfigApp ComponentKey = "app"
	// ComponentLogger is the key for the logger configuration.
	ComponentLogger ComponentKey = "logger"
	// ComponentData is the key for the data sources (databases, caches) configuration.
	ComponentData ComponentKey = "data"
	// ComponentDatabases is the key for the database configurations.
	ComponentDatabases ComponentKey = "databases"
	// ComponentCaches is the key for the cache configurations.
	ComponentCaches ComponentKey = "caches"
	// ComponentObjectStores is the key for the object store configurations.
	ComponentObjectStores ComponentKey = "object_stores"
	// ComponentRegistries is the key for the service discovery/registry configuration.
	ComponentRegistries ComponentKey = "discoveries"
	// ComponentDefaultRegistry is the key for the default registry name.
	ComponentDefaultRegistry ComponentKey = "default_registry_name"
	// ComponentMiddlewares is the key for the middleware configuration.
	ComponentMiddlewares ComponentKey = "middlewares"
	// ComponentServers is the key for the server (e.g., HTTP, gRPC) configuration.
	ComponentServers ComponentKey = "servers"
	// ComponentClients is the key for the client/service consumer configuration.
	ComponentClients ComponentKey = "clients"
)