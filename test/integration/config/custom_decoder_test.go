package config_test

import (
	"testing"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/stretchr/testify/assert"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/logger/v1"
	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/config/decoder" // Import the public decoder package
	"github.com/origadmin/runtime/interfaces"
)

// TestCustomSettings represents the structure of our custom configuration section for testing.
type TestCustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"endpoints"`
}

// customTestConfigDecoder implements the interfaces.ConfigDecoder interface for testing.
// It embeds decoder.BaseDecoder and overrides specific methods to return ErrNotImplemented.
type customTestConfigDecoder struct {
	*decoder.BaseDecoder
}

// DecodeLogger overrides the BaseDecoder's DecodeLogger to return ErrNotImplemented.
// This forces the runtime to fall back to the generic Decode method for logger config.
func (d *customTestConfigDecoder) DecodeLogger() (*loggerv1.Logger, error) {
	return nil, interfaces.ErrNotImplemented
}

// DecodeDiscoveries overrides the BaseDecoder's DecodeDiscoveries to return ErrNotImplemented.
// This forces the runtime to fall back to the generic Decode method for discovery configs.
func (d *customTestConfigDecoder) DecodeDiscoveries() (map[string]*discoveryv1.Discovery, error) {
	return nil, interfaces.ErrNotImplemented
}

// customTestDecoderProvider implements the interfaces.ConfigDecoderProvider interface for testing.
type customTestDecoderProvider struct{}

// GetConfigDecoder returns a new customTestConfigDecoder.
func (p *customTestDecoderProvider) GetConfigDecoder(kratosConfig kratosconfig.Config) (interfaces.ConfigDecoder, error) {
	return &customTestConfigDecoder{
		BaseDecoder: decoder.NewBaseDecoder(kratosConfig),
	}, nil
}

// TestCustomConfigDecoderIntegration tests the integration of a custom ConfigDecoder
// that relies on BaseDecoder and returns ErrNotImplemented for specific fast paths.
func TestCustomConfigDecoderIntegration(t *testing.T) {
	assert := assert.New(t)

	// 1. Initialize Runtime with the custom decoder provider.
	rt, cleanup, err := runtime.NewFromBootstrap(
		"./config/full_config.yaml",
		runtime.WithAppInfo(runtime.AppInfo{
			ID:      "test-custom-decoder",
			Name:    "TestCustomDecoder",
			Version: "1.0.0",
			Env:     "test",
		}),
		runtime.WithDecoderProvider(&customTestDecoderProvider{}),
	)
	assert.NoError(err)
	defer cleanup()

	// 2. Get the ConfigDecoder from the runtime.
	configDecoder := rt.Config()
	assert.NotNil(configDecoder)

	// 3. Verify Logger configuration (should use generic Decode due to ErrNotImplemented).
	logger := rt.Logger()
	assert.NotNil(logger)

	// We expect the logger level to be "info" as defined in config.yaml
	// Note: The actual slog.Level value might not be directly comparable, but we can check its string representation or behavior.
	// For simplicity, we'll just assert that the logger was created without error.

	// 4. Verify Registries configuration (should use generic Decode due to ErrNotImplemented).
	// The getRegistriesConfig function in runtime.go will attempt to decode into registriesConfig struct.
	// Since our test config.yaml doesn't define registries, we expect it to be empty.
	// We can't directly access the internal registriesCfg, but we can check if default registrar is nil.
	assert.Nil(rt.DefaultRegistrar(), "Default registrar should be nil if no registries are configured")

	// 5. Verify custom_settings are decoded correctly using the generic Decode method.
	var customSettings TestCustomSettings
	err = configDecoder.Decode("custom_settings", &customSettings)
	assert.NoError(err)
	assert.True(customSettings.FeatureEnabled)
	assert.Equal("super-secret-key-123", customSettings.APIKey)
	assert.Equal(100, customSettings.RateLimit)
	assert.Len(customSettings.Endpoints, 2)
	assert.Equal("users", customSettings.Endpoints[0].Name)
	assert.Equal("/api/v1/users", customSettings.Endpoints[0].Path)
	assert.Equal("products", customSettings.Endpoints[1].Name)
	assert.Equal("/api/v1/products", customSettings.Endpoints[1].Path)

	// 6. Verify a standard config section (e.g., servers) is decoded correctly.
	var servers []struct {
		Http *struct {
			Network string `json:"network"`
			Addr    string `json:"addr"`
			Timeout string `json:"timeout"`
		}
		Grpc *struct {
			Network string `json:"network"`
			Addr    string `json:"addr"`
			Timeout string `json:"timeout"`
		}
	}
	err = configDecoder.Decode("servers", &servers)
	assert.NoError(err)
	assert.Len(servers, 2)
	assert.NotNil(servers[0].Http)
	assert.Equal("tcp", servers[0].Http.Network)
	assert.Equal(":8080", servers[0].Http.Addr)
	assert.Equal("1s", servers[0].Http.Timeout)
	assert.NotNil(servers[1].Grpc)
	assert.Equal("tcp", servers[1].Grpc.Network)
	assert.Equal(":9090", servers[1].Grpc.Addr)
	assert.Equal("1s", servers[1].Grpc.Timeout)
}
