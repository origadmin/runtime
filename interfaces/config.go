package interfaces

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
)

type Config map[string]any

type ConfigDecoderFunc func(config kratosconfig.Config) (ConfigDecoder, error)

func (c ConfigDecoderFunc) GetConfigDecoder(config kratosconfig.Config) (ConfigDecoder, error) {
	return c(config)
}

// ConfigDecoder provides a generic way to decode a portion of the configuration
// into a Go struct.
type ConfigDecoder interface {
	Config() any
	// Decode unmarshals the configuration section identified by 'key' into the 'target'.
	// The 'key' can be a dot-separated path (e.g., "service.http").
	// 'target' must be a pointer to a Go struct.
	Decode(key string, target interface{}) error
}

type ConfigDecoderProvider interface {
	GetConfigDecoder(config kratosconfig.Config) (ConfigDecoder, error)
}

//// ServiceConfig provides access to service configurations.
//type ServiceConfig interface {
//	GetService(name string) *configv1.Service
//	GetServices() map[string]*configv1.Service
//}
type ServiceConfig interface{}

// DiscoveryConfig provides access to discovery/registry configurations.
type DiscoveryConfig interface {
	GetDiscovery(name string) *discoveryv1.Discovery
	GetDiscoveries() map[string]*discoveryv1.Discovery
}

// LoggerConfig provides access to the logger configuration.
type LoggerConfig interface {
	GetLogger() *loggerv1.Logger
}

// MiddlewareConfig provides access to middleware configurations.
type MiddlewareConfig interface {
	GetMiddleware(name string) *middlewarev1.MiddlewareConfig
	GetMiddlewares() map[string]*middlewarev1.MiddlewareConfig
}
