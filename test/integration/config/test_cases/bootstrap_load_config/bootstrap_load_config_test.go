package bootstrap_load_config_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// AssertTestConfig contains the validation logic for the TestConfig struct
// specific to the bootstrap_load_config test case, omitting discovery assertions.
func AssertTestConfig(t *testing.T, cfg *testconfigs.TestConfig) {
	asserts := assert.New(t)

	// App configuration assertions
	asserts.NotNil(cfg.App)
	asserts.Equal("test-app-id", cfg.App.GetId())
	asserts.Equal("TestApp", cfg.App.GetName())
	asserts.Equal("1.0.0", cfg.App.GetVersion())
	asserts.Equal("test", cfg.App.GetEnv())
	asserts.Contains(cfg.App.GetMetadata(), "key1")
	asserts.Contains(cfg.App.GetMetadata(), "key2")
	asserts.Equal("value1", cfg.App.GetMetadata()["key1"])
	asserts.Equal("value2", cfg.App.GetMetadata()["key2"])

	// Server configuration assertions
	asserts.Len(cfg.GetServers(), 1, "Should have 1 Servers message")
	serverConfigs := cfg.GetServers()[0].GetServers()
	asserts.Len(serverConfigs, 2, "Should have 2 Server configurations (gRPC and HTTP)")

	var grpcServer *transportv1.Server
	var httpServer *transportv1.Server

	for _, s := range serverConfigs {
		if s.GetGrpc() != nil {
			grpcServer = s
		}
		if s.GetHttp() != nil {
			httpServer = s
		}
	}

	// Verify gRPC server config
	asserts.NotNil(grpcServer, "gRPC server configuration not found")
	asserts.Equal("tcp", grpcServer.GetGrpc().GetNetwork())
	asserts.Equal(":9000", grpcServer.GetGrpc().GetAddr())
	asserts.Equal("1s", grpcServer.GetGrpc().GetTimeout().AsDuration().String())

	// Verify HTTP server config
	asserts.NotNil(httpServer, "HTTP server configuration not found")
	asserts.Equal("tcp", httpServer.GetHttp().GetNetwork())
	asserts.Equal(":8000", httpServer.GetHttp().GetAddr())
	asserts.Equal("2s", httpServer.GetHttp().GetTimeout().AsDuration().String())

	// Client configuration assertions
	asserts.NotNil(cfg.Client)
	asserts.Equal("discovery:///user-service", cfg.Client.GetEndpoint())
	asserts.Equal("3s", cfg.Client.GetTimeout().AsDuration().String())
	asserts.NotNil(cfg.Client.GetSelector())
	asserts.Equal("v1.0.0", cfg.Client.GetSelector().GetVersion())

	// Discovery configuration assertions (omitted for this test case)
	asserts.Nil(cfg.Discoveries, "Discoveries should be empty for this test case")
	asserts.Empty(cfg.GetRegistrationDiscoveryName(), "RegistrationDiscoveryName should be empty for this test case")

	// Logger configuration assertions
	asserts.NotNil(cfg.Logger)
	asserts.Equal("info", cfg.Logger.GetLevel())
	asserts.Equal("json", cfg.Logger.GetFormat())
	asserts.True(cfg.Logger.GetStdout())

	// Tracer configuration assertions
	asserts.NotNil(cfg.Tracer)
	asserts.Equal("jaeger", cfg.Tracer.GetName())
	asserts.Equal("http://jaeger:14268/api/traces", cfg.Tracer.GetEndpoint())

	// Middleware configuration assertions
	asserts.NotNil(cfg.Middlewares)
	asserts.Len(cfg.Middlewares.GetMiddlewares(), 1, "Should have 1 middleware configured")

	corsMiddleware := cfg.Middlewares.GetMiddlewares()[0]
	asserts.Equal("cors", corsMiddleware.GetType())
	asserts.True(corsMiddleware.GetEnabled())
}

// TestCustomSettings represents the structure for custom configuration section in tests
type TestCustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"endpoints"`
}

// RuntimeIntegrationTestSuite tests configuration integration with Runtime
type RuntimeIntegrationTestSuite struct {
	suite.Suite
}

func TestRuntimeIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RuntimeIntegrationTestSuite))
}

// TestRuntimeLoadCompleteConfig tests loading complete configuration using Runtime
func (s *RuntimeIntegrationTestSuite) TestRuntimeLoadCompleteConfig() {
	t := s.T()
	asserts := assert.New(t)

	// Use a robust relative path to the dedicated bootstrap config for this test.
	bootstrapPath := filepath.Join("testdata", "complete_config", "bootstrap.yaml")

	// Initialize Runtime
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-complete-config",
			Name:    "TestCompleteConfig",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer rtInstance.Cleanup()

	// Get configuration decoder
	configDecoder := rtInstance.Config()
	asserts.NotNil(configDecoder)

	// Decode into complete config structure
	var config testconfigs.TestConfig
	err = configDecoder.Decode("", &config)
	asserts.NoError(err)

	// Run assertions
	AssertTestConfig(t, &config)
	t.Logf("Runtime loaded and verified complete config successfully!")
}

// TestConfigProtoIntegration tests integration between configuration and Protocol Buffers
func (s *RuntimeIntegrationTestSuite) TestConfigProtoIntegration() {
	t := s.T()
	assertions := assert.New(t)

	// Use a robust relative path to the dedicated bootstrap config for this test.
	bootstrapPath := filepath.Join("testdata", "proto_integration", "bootstrap.yaml")

	// 1. Initialize Runtime with default decoder provider
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-proto-config",
			Name:    "TestProtoConfig",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer rtInstance.Cleanup()

	// 2. Get ConfigDecoder from runtime
	configDecoder := rtInstance.Config()
	assertions.NotNil(configDecoder)

	// 3. Decode entire config into generated Bootstrap struct
	var bootstrapConfig testconfigs.TestConfig
	err = configDecoder.Decode("", &bootstrapConfig)
	assertions.NoError(err)

	// 4. Assert decoded values
	// Verify logger (from test_config.yaml)
	logger := rtInstance.Logger()
	assertions.NotNil(logger)

	// Verify registration_discovery_name
	assertions.Equal("test-discovery", bootstrapConfig.RegistrationDiscoveryName)

	// Verify servers
	assertions.Len(bootstrapConfig.GetServers(), 1, "Should have 1 Servers message")
	serverConfigs := bootstrapConfig.GetServers()[0].GetServers()
	assertions.Len(serverConfigs, 2, "Should have 2 Server configurations (gRPC and HTTP)")

	var grpcServer *transportv1.Server
	var httpServer *transportv1.Server

	for _, s := range serverConfigs {
		if s.GetGrpc() != nil {
			grpcServer = s
		}
		if s.GetHttp() != nil {
			httpServer = s
		}
	}

	assertions.NotNil(grpcServer, "gRPC server configuration not found")
	assertions.Equal("tcp", grpcServer.GetGrpc().GetNetwork())
	assertions.Equal(":9000", grpcServer.GetGrpc().GetAddr())
	assertions.Equal("1s", grpcServer.GetGrpc().GetTimeout().AsDuration().String())

	assertions.NotNil(httpServer, "HTTP server configuration not found")
	assertions.Equal("tcp", httpServer.GetHttp().GetNetwork())
	assertions.Equal(":8000", httpServer.GetHttp().GetAddr())
	assertions.Equal("2s", httpServer.GetHttp().GetTimeout().AsDuration().String())

	// Verify clients (should be empty as not defined in test_config.yaml)
	assertions.Empty(bootstrapConfig.Clients)

	// Verify discoveries (should be empty as not defined in test_config.yaml)
	assertions.Empty(bootstrapConfig.Discoveries)
}

// TestRuntimeDecoder verifies that the configuration decoder is properly exposed by runtime
// and can be used to parse custom configuration sections
func (s *RuntimeIntegrationTestSuite) TestRuntimeDecoder() {
	t := s.T()
	asserts := assert.New(t)

	// Use a robust relative path to the dedicated bootstrap config for this test.
	bootstrapPath := filepath.Join("testdata", "decoder_test", "bootstrap.yaml")

	// 1. Initialize Runtime with correct AppInfo
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-decoder",
			Name:    "TestDecoder",
			Version: "1.0.0",
			Env:     "test",
		}),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer rtInstance.Cleanup()

	// 2. Verify core components are properly initialized
	asserts.NotNil(rtInstance.Logger())
	asserts.Equal("test-decoder", rtInstance.AppInfo().ID)
	asserts.Equal("TestDecoder", rtInstance.AppInfo().Name)
	asserts.Equal("1.0.0", rtInstance.AppInfo().Version)

	// 3. Get ConfigDecoder from runtime
	decoder := rtInstance.Config()
	asserts.NotNil(decoder, "ConfigDecoder should not be nil")

	// 4. Verify we can decode the entire config into a map
	var configMap map[string]interface{}
	err = decoder.Decode("", &configMap)
	asserts.NoError(err, "Failed to decode config into map")
	asserts.NotEmpty(configMap, "Decoded config map should not be empty")

	// 5. Verify custom_settings is properly decoded using the exposed decoder
	var customSettings TestCustomSettings
	// Updated to use "components.my-custom-settings" path to match the config structure
	err = decoder.Decode("components.my-custom-settings", &customSettings)
	asserts.NoError(err, "Failed to decode custom settings")
	asserts.True(customSettings.FeatureEnabled, "Feature should be enabled")
	asserts.Equal(100, customSettings.RateLimit, "Rate limit should be 100")
	asserts.Len(customSettings.Endpoints, 2, "Should have 2 endpoints")
	asserts.Equal("users", customSettings.Endpoints[0].Name, "First endpoint should be 'users'")
	asserts.Equal("/api/v1/users", customSettings.Endpoints[0].Path, "Users endpoint path should be correct")
	asserts.Equal("products", customSettings.Endpoints[1].Name, "Second endpoint should be 'products'")
	asserts.Equal("/api/v1/products", customSettings.Endpoints[1].Path, "Products endpoint path should be correct")
}
