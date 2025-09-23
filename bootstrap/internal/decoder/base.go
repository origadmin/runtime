package decoder

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	"github.com/origadmin/runtime/interfaces"
)

// baseDecoder provides a default implementation for the ConfigDecoder interface.
// Custom decoders can embed this struct and only override the methods they need.
// This avoids the need for runtime type assertions and provides a stable interface.
type baseDecoder struct {
	KratosConfig kratosconfig.Config
}

// NewDecoder creates a new baseDecoder.
func newBaseDecoder(config kratosconfig.Config) *baseDecoder {
	return &baseDecoder{KratosConfig: config}
}

// Decode provides a generic fallback for decoding using the underlying Kratos config.
func (b *baseDecoder) Decode(key string, value any) error {
	// If key is empty, scan the entire config. Otherwise, scan the specific key.
	if key == "" {
		return b.KratosConfig.Scan(value)
	}
	return b.KratosConfig.Value(key).Scan(value)
}

// DecodeLogger returns ErrNotImplemented by default.
// Custom decoders should override this method to provide a fast path for logger config.
func (b *baseDecoder) DecodeLogger() (*loggerv1.Logger, error) {
	return nil, interfaces.ErrNotImplemented
}

// DecodeDiscoveries returns ErrNotImplemented by default.
// Custom decoders should override this to provide a fast path for discovery configs.
func (b *baseDecoder) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	return nil, interfaces.ErrNotImplemented
}

func (b *baseDecoder) Config() kratosconfig.Config {
	return b.KratosConfig
}

func (b *baseDecoder) Close() error {
	return b.KratosConfig.Close()
}
