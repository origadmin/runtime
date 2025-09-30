package interfaces

import (
	"errors"

	appv1 "github.com/origadmin/runtime/api/gen/go/app/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
)

// ErrNotImplemented is returned when a specific decoder method is not implemented
// by a custom decoder. This signals the runtime to fall back to generic decoding.
var ErrNotImplemented = errors.New("method not implemented by this decoder")

// Config is the minimal contract for providing a custom configuration source.
// Developers wishing to extend the framework with a new config system should implement this interface.
type Config interface {
	// Load loads the configuration from its source.
	Load() error

	// Decode provides generic decoding of a configuration key into a target struct.
	// This is the fundamental method that MUST be implemented by any Config instance.
	Decode(key string, value any) error

	// Raw provides an "escape hatch" to the underlying Kratos config.Config instance.
	// Custom implementations can return nil if not applicable.
	Raw() any

	// Close releases any resources held by the configuration.
	// MUST be implemented; can be a no-op if no resources are held.
	Close() error
}

// StructuredConfig defines a set of type-safe, recommended methods for decoding configuration.
// This is the interface that should be exposed to the end user.
// It embeds the generic Config interface to allow for decoding arbitrary values.
type StructuredConfig interface {
	Config // Embed the generic config interface
	AppConfigDecoder
	LoggerConfigDecoder
	DiscoveriesConfigDecoder
	MiddlewareConfigDecoder
	// Add other specific decoders here
}

type AppConfigDecoder interface {
	DecodeApp() (*appv1.App, error)
}

// LoggerConfigDecoder defines an OPTIONAL interface for providing a "fast path"
// to decode logger configuration. Custom Config implementations can implement this
// interface to provide an optimized decoding path.
type LoggerConfigDecoder interface {
	DecodeLogger() (*loggerv1.Logger, error)
}

// DiscoveriesConfigDecoder defines an OPTIONAL interface for providing a "fast path"
// to decode service discovery configurations. Custom Config implementations can implement this
// interface to provide an optimized decoding path.
type DiscoveriesConfigDecoder interface {
	DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error)
}

// MiddlewareConfigDecoder defines an OPTIONAL interface for providing a "fast path"
// to decode middleware configurations. Custom Config implementations can implement this
// interface to provide an optimized decoding path.
type MiddlewareConfigDecoder interface {
	DecodeMiddleware() (*middlewarev1.Middlewares, error)
}
