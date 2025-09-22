package interfaces

import (
	"errors"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
)

// ErrNotImplemented is returned when a specific decoder method is not implemented
// by a custom decoder. This signals the runtime to fall back to generic decoding.
var ErrNotImplemented = errors.New("method not implemented by this decoder")

// LoggerConfigDecoder defines the interface for decoding logger configuration.
type LoggerConfigDecoder interface {
	DecodeLogger() (*loggerv1.Logger, error)
}

// DiscoveriesConfigDecoder defines the interface for decoding service discovery configurations.
type DiscoveriesConfigDecoder interface {
	DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error)
}

// Config defines the interface for the application's configuration, providing
// both generic decoding and specialized "fast path" decoding for common components.
// It embeds smaller, more specific decoder interfaces for better organization.
type Config interface {
	// Decode provides generic decoding of a configuration key into a target struct.
	// This is the fallback mechanism if a specialized method is not implemented.
	Decode(key string, value any) error

	// Raw returns the underlying Kratos config.Config instance.
	// This allows access to the raw configuration source for advanced scenarios.
	Raw() kratosconfig.Config

	LoggerConfigDecoder      // Embed LoggerConfigDecoder
	DiscoveriesConfigDecoder // Embed DiscoveriesConfigDecoder
}

// ConfigProvider defines the interface for creating a Config from a Kratos config.
type ConfigProvider interface {
	New(kratosconfig.Config) (Config, error) // Renamed from NewDecoder
}

// ConfigProviderFunc is an adapter to allow the use of ordinary functions as ConfigProviders.
type ConfigProviderFunc func(kratosconfig.Config) (Config, error)

// New calls f(config).
func (f ConfigProviderFunc) New(config kratosconfig.Config) (Config, error) {
	return f(config)
}
