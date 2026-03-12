package contracts

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
)

// AppConfig defines the contract for a configuration that provides app information.
type AppConfig interface {
	GetApp() *appv1.App
}

// LoggerConfig defines the contract for a configuration that provides logger information.
type LoggerConfig interface {
	GetLogger() *loggerv1.Logger
}

// MiddlewareConfig defines the contract for a configuration that provides middleware information.
type MiddlewareConfig interface {
	GetMiddlewares() *middlewarev1.Middlewares
}

// DataConfig defines the contract for a configuration that provides data information.
type DataConfig interface {
	GetData() *datav1.Data
}

// DiscoveryConfig defines the contract for a configuration that provides discovery information.
type DiscoveryConfig interface {
	GetDiscoveries() *discoveryv1.Discoveries
}

// ServerConfig defines the contract for a configuration that provides server information.
type ServerConfig interface {
	GetServers() *transportv1.Servers
}

// ClientConfig defines the contract for a configuration that provides client information.
type ClientConfig interface {
	GetClients() *transportv1.Clients
}

// ConfigObject aggregates common configuration accessors.
type ConfigObject interface {
	AppConfig
	LoggerConfig
	DiscoveryConfig
	MiddlewareConfig
	DataConfig
	ServerConfig
	ClientConfig
}
