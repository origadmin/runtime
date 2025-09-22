package decoder

import (
	"errors"

	"github.com/origadmin/runtime/api/gen/go/discovery/v1"
	"github.com/origadmin/runtime/api/gen/go/logger/v1"
	"github.com/origadmin/runtime/interfaces"
)

// pathAwareDecoder is a wrapper that implements the "fast path first, fallback to generic"
// decoding strategy. It uses a ConfigPaths map to know which key to use for the fallback.
type pathAwareDecoder struct {
	originalDecoder interfaces.ConfigDecoder
	paths           map[string]string
}

// NewPathAwareDecoder creates a new smart wrapper decoder.
func NewPathAwareDecoder(originalDecoder interfaces.ConfigDecoder, paths map[string]string) interfaces.ConfigDecoder {
	return &pathAwareDecoder{
		originalDecoder: originalDecoder,
		paths:           paths,
	}
}

// DecodeLogger first attempts to call the wrapped decoder's DecodeLogger.
// If that is not implemented, it falls back to calling the wrapped decoder's
// generic Decode method, using the path from its ConfigPaths map.
func (d *pathAwareDecoder) DecodeLogger() (*loggerv1.Logger, error) {
	// 1. Try the fast path first.
	loggerConfig, err := d.originalDecoder.DecodeLogger()
	if !errors.Is(err, interfaces.ErrNotImplemented) {
		// This means it was either a success (err == nil) or a real error.
		// In both cases, we return directly.
		return loggerConfig, err
	}

	// 2. Fallback to generic decoding using the configured path.
	path, ok := d.paths["logger"]
	if !ok || path == "" {
		// If no path is configured, we can't fall back.
		// Return the original ErrNotImplemented so the caller knows nothing was found.
		return nil, interfaces.ErrNotImplemented
	}

	var cfg loggerv1.Logger
	if decodeErr := d.originalDecoder.Decode(path, &cfg); decodeErr != nil {
		return nil, decodeErr
	}
	return &cfg, nil
}

// DecodeDiscoveries implements the same "fast path first, fallback to generic" strategy.
func (d *pathAwareDecoder) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	// 1. Try the fast path first.
	discoveries, err := d.originalDecoder.DecodeDiscoveries()
	if !errors.Is(err, interfaces.ErrNotImplemented) {
		return discoveries, err
	}

	// 2. Fallback to generic decoding.
	path, ok := d.paths["registries"] // Note: the component name is "registries"
	if !ok || path == "" {
		return nil, interfaces.ErrNotImplemented
	}

	var cfg map[string]*discoveryv1.Discovery
	if decodeErr := d.originalDecoder.Decode(path, &cfg); decodeErr != nil {
		return nil, decodeErr
	}
	return cfg, nil
}

// Decode simply passes the call through to the original decoder.
func (d *pathAwareDecoder) Decode(key string, value any) error {
	return d.originalDecoder.Decode(key, value)
}
