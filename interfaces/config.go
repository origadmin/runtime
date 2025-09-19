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

// Config defines the interface for the application's configuration, providing
// access to both the raw configuration source and the decoding capabilities.
type Config interface {
	// Raw returns the underlying Kratos config.Config instance.
	// It is primarily used for advanced scenarios where direct access to the source is needed.
	Raw() kratosconfig.Config

	// Decoder returns the ConfigDecoder for accessing and decoding configuration values.
	Decoder() ConfigDecoder
}

// ConfigDecoder defines the interface for decoding configuration values. It supports
// both generic decoding and specialized "fast path" decoding for common components.
//
// This design allows for custom decoders to be implemented by embedding BaseDecoder
// and overriding only the necessary methods, avoiding the need for runtime type assertions.
type ConfigDecoder interface {
	// Decode provides generic decoding of a configuration key into a target struct.
	// This is the fallback mechanism if a specialized method is not implemented.
	Decode(key string, value any) error

	// DecodeLogger provides a fast path for decoding the logger configuration.
	// If not implemented, it should return ErrNotImplemented.
	DecodeLogger() (*loggerv1.Logger, error)

	// DecodeDiscoveries provides a fast path for decoding the service discovery configurations.
	// If not implemented, it should return ErrNotImplemented.
	DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error)
}

// ConfigDecoderProvider defines the interface for creating a ConfigDecoder from a Kratos config.
type ConfigDecoderProvider interface {
	GetConfigDecoder(kratosconfig.Config) (ConfigDecoder, error)
}

// ConfigDecoderFunc is an adapter to allow the use of ordinary functions as ConfigDecoderProviders.

type ConfigDecoderFunc func(kratosconfig.Config) (ConfigDecoder, error)

// GetConfigDecoder calls f(config).
func (f ConfigDecoderFunc) GetConfigDecoder(config kratosconfig.Config) (ConfigDecoder, error) {
	return f(config)
}
