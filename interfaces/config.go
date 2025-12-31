package interfaces

import (
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	"github.com/origadmin/toolkits/errors"
)

// ErrNotImplemented is returned when a specific decoder method is not implemented
// by a custom decoder. This signals the runtime to fall back to generic decoding.
var ErrNotImplemented = errors.New("method not implemented")

// ConfigLoader is the minimal contract for providing a custom configuration source.
// Developers wishing to extend the framework with a new config system should implement this interface.
type ConfigLoader interface {
	// Load loads the configuration from its source.
	Load() error

	// Decode provides generic decoding of a configuration key into a target struct.
	// This is the fundamental method that MUST be implemented by any ConfigLoader instance.
	Decode(key string, value any) error

	// Raw provides an "escape hatch" to the underlying Kratos config.ConfigLoader instance.
	// Custom implementations can return nil if not applicable.
	Raw() any

	// Close releases any resources held by the configuration.
	// MUST be implemented; can be a no-op if no resources are held.
	Close() error
}

// StructuredConfig defines a set of type-safe, recommended methods for decoding configuration.
// It embeds the generic ConfigLoader interface to allow for decoding arbitrary values.
type StructuredConfig interface {
	DataConfigDecoder
	CacheConfigDecoder
	DatabaseConfigDecoder
	ObjectStoreConfigDecoder
	DiscoveriesConfigDecoder
	LoggerConfigDecoder
	MiddlewareConfigDecoder
	ServiceConfigDecoder

	// DecodedConfig returns the config original decoded value.
	DecodedConfig() any
}

type ConfigObject interface {
	GetLogger() *loggerv1.Logger
	GetDiscoveries() *discoveryv1.Discoveries
	GetMiddlewares() *middlewarev1.Middlewares
	GetServers() *transportv1.Servers
	GetClients() *transportv1.Clients
	GetData() *datav1.Data
}

//type AppConfigDecoder interface {
//	DecodeApp() (*appv1.App, error)
//}

// LoggerConfigDecoder defines an OPTIONAL interface for providing a "fast path"
// to decode logger configuration. Custom ConfigLoader implementations can implement this
// interface to provide an optimized decoding path.
type LoggerConfigDecoder interface {
	DecodeLogger() (*loggerv1.Logger, error)
}

// DiscoveriesConfigDecoder defines an OPTIONAL interface for providing a "fast path"
// to decode service discovery configurations. Custom ConfigLoader implementations can implement this
// interface to provide an optimized decoding path.
type DiscoveriesConfigDecoder interface {
	DecodeDefaultDiscovery() (string, error)
	DecodeDiscoveries() (*discoveryv1.Discoveries, error)
}

// MiddlewareConfigDecoder defines an OPTIONAL interface for providing a "fast path"
// to decode middleware configurations. Custom ConfigLoader implementations can implement this
// interface to provide an optimized decoding path.
type MiddlewareConfigDecoder interface {
	DecodeMiddlewares() (*middlewarev1.Middlewares, error)
}

// ServiceConfigDecoder defines an OPTIONAL interface for providing a "fast path"
// to decode service configurations. Custom ConfigLoader implementations can implement this
// interface to provide an optimized decoding path.
type ServiceConfigDecoder interface {
	DecodeServers() (*transportv1.Servers, error)
	DecodeClients() (*transportv1.Clients, error)
}

type DataConfigDecoder interface {
	DecodeData() (*datav1.Data, error)
}

type CacheConfigDecoder interface {
	DecodeCaches() (*datav1.Caches, error)
}

type DatabaseConfigDecoder interface {
	DecodeDatabases() (*datav1.Databases, error)
}

type ObjectStoreConfigDecoder interface {
	DecodeObjectStores() (*datav1.ObjectStores, error)
}
