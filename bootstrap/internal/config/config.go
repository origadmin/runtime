package config

import (
	"errors"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	"github.com/origadmin/runtime/bootstrap/constant"
	"github.com/origadmin/runtime/interfaces"
)

// configImpl is the default implementation of the interfaces.Config interface.
// It also implements all optional decoder interfaces (LoggerConfigDecoder, etc.)
// to provide an optimized "fast path" for the bootstrap process.
type configImpl struct {
	kratosConfig kratosconfig.Config
	paths        map[string]string
}

// Statically assert that configImpl implements all required interfaces.
var (
	_ interfaces.Config                   = (*configImpl)(nil)
	_ interfaces.LoggerConfigDecoder      = (*configImpl)(nil)
	_ interfaces.DiscoveriesConfigDecoder = (*configImpl)(nil)
)

// NewConfigImpl creates a new instance of the default config implementation.
func NewConfigImpl(kc kratosconfig.Config, paths map[string]string) interfaces.Config {
	// Ensure paths map is not nil to prevent panics.
	if paths == nil {
		paths = make(map[string]string)
	}
	return &configImpl{
		kratosConfig: kc,
		paths:        paths,
	}
}

// Decode implements the interfaces.Config interface.
func (c *configImpl) Decode(key string, value any) error {
	return c.kratosConfig.Value(key).Scan(value)
}

// Raw implements the interfaces.Config interface.
func (c *configImpl) Raw() kratosconfig.Config {
	return c.kratosConfig
}

// Close implements the interfaces.Config interface.
func (c *configImpl) Close() error {
	return c.kratosConfig.Close()
}

// DecodeLogger implements the interfaces.LoggerConfigDecoder interface.
func (c *configImpl) DecodeLogger() (*loggerv1.Logger, error) {
	path, ok := c.paths[constant.ComponentLogger]
	if !ok {
		return nil, errors.New("logger component path not configured")
	}

	loggerConfig := new(loggerv1.Logger)
	if err := c.kratosConfig.Value(path).Scan(loggerConfig); err != nil {
		return nil, err
	}
	return loggerConfig, nil
}

// DecodeDiscoveries implements the interfaces.DiscoveriesConfigDecoder interface.
func (c *configImpl) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	path, ok := c.paths[constant.ComponentRegistries]
	if !ok {
		return nil, errors.New("registries component path not configured")
	}

	var discoveries map[string]*discoveryv1.Discovery
	if err := c.kratosConfig.Value(path).Scan(&discoveries); err != nil {
		return nil, err
	}

	return discoveries, nil
}
