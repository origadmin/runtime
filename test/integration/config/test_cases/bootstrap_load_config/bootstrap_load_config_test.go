package bootstrap_load_config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/bootstrap/constant"
	"github.com/origadmin/runtime/interfaces"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

const (
	testCaseDir = "test/integration/config/test_cases/bootstrap_load_config"
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

// BootstrapLoadConfigTestSuite tests the loading of configuration via the bootstrap process.
type BootstrapLoadConfigTestSuite struct {
	suite.Suite
}

func TestBootstrapLoadConfigTestSuite(t *testing.T) {
	suite.Run(t, new(BootstrapLoadConfigTestSuite))
}

// TestBootstrapLoading verifies that the runtime can be initialized from the bootstrap file,
// which then correctly loads and merges the main application configuration.
func (s *BootstrapLoadConfigTestSuite) TestBootstrapLoading() {
	t := s.T()
	asserts := assert.New(t)

	// Get the path to the current test file
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)

	// Define test cases for different formats of the bootstrap config.
	testCases := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: filepath.Join(testDir, "bootstrap.yaml")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Runtime from the bootstrap file.
			// The AppInfo here is minimal as the main config is loaded via bootstrap.
			rtInstance, err := rt.NewFromBootstrap(
				tc.filePath,
				bootstrap.WithAppInfo(&interfaces.AppInfo{
					ID:      "bootstrap-test-app",
					Name:    "BootstrapTestApp",
					Version: "1.0.0",
				}),
				bootstrap.WithDefaultPaths(map[string]string{
					constant.ComponentMiddlewares: "middlewares.middlewares",
				}),
			)
			asserts.NoError(err, "Failed to initialize runtime from bootstrap: %v", err)
			defer rtInstance.Cleanup()

			// Get the configuration decoder from the runtime.
			configDecoder := rtInstance.Config()
			asserts.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil")

			// Decode the entire configuration into our unified TestConfig struct.
			var cfg testconfigs.TestConfig
			err = configDecoder.Decode("", &cfg) // Decode the root of the config

			asserts.NoError(err, "Failed to decode config from runtime: %v", err)
			// Run assertions on the loaded configuration using the custom assertion.
			AssertTestConfig(t, &cfg)
			t.Logf("Runtime loaded and verified %s bootstrap config successfully!", tc.name)
		})
	}
}

// TestConfigProtoIntegration tests integration between configuration and Protocol Buffers
func (s *BootstrapLoadConfigTestSuite) TestConfigProtoIntegration() {
	t := s.T()
	asserts := assert.New(t)

	// All test files are in the same directory, use relative paths directly

	// 1. Initialize Runtime with default decoder provider
	rtInstance, err := rt.NewFromBootstrap(
		"bootstrap.yaml",
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-app-id",
			Name:    "TestApp",
			Version: "1.0.0",
		}),
	)
	require.NoError(t, err, "Failed to initialize runtime")
	defer rtInstance.Cleanup()

	// 2. Get ConfigDecoder from runtime
	configDecoder := rtInstance.Config()
	asserts.NotNil(configDecoder, "ConfigDecoder should not be nil")

	// 3. Decode entire config into generated Bootstrap struct
	var bootstrapConfig testconfigs.TestConfig
	err = configDecoder.Decode("", &bootstrapConfig)
	asserts.NoError(err, "Failed to decode config")

	// 4. Assert decoded values
	// Verify app info
	appInfo := rtInstance.AppInfo()
	asserts.Equal("test-app-id", appInfo.ID)
	asserts.Equal("TestApp", appInfo.Name)
	asserts.Equal("1.0.0", appInfo.Version)

	// Verify servers
	servers := bootstrapConfig.GetServers()
	asserts.Len(servers, 1, "Should have 1 Servers message")
	serverConfigs := servers[0].GetServers()
	asserts.Len(serverConfigs, 2, "Should have 2 Server configurations (gRPC and HTTP)")

	var grpcServer, httpServer *transportv1.Server
	for _, s := range serverConfigs {
		switch {
		case s.GetGrpc() != nil:
			grpcServer = s
		case s.GetHttp() != nil:
			httpServer = s
		}
	}

	// Verify gRPC server config
	asserts.NotNil(grpcServer, "gRPC server configuration not found")
	asserts.Equal("tcp", grpcServer.GetGrpc().GetNetwork())
	asserts.Equal(":9000", grpcServer.GetGrpc().GetAddr())

	// Verify HTTP server config
	asserts.NotNil(httpServer, "HTTP server configuration not found")
	asserts.Equal("tcp", httpServer.GetHttp().GetNetwork())
	asserts.Equal(":8000", httpServer.GetHttp().GetAddr())
}

// TestRuntimeDecoder verifies that the configuration decoder is properly exposed by runtime
// and can be used to parse custom configuration sections
func (s *BootstrapLoadConfigTestSuite) TestRuntimeDecoder() {
	t := s.T()
	asserts := assert.New(t)

	// All test files are in the same directory, use relative paths directly
	wd, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")
	t.Logf("Current working directory: %s", wd)
	// 1. Initialize Runtime with correct AppInfo
	bootstrapPath := filepath.Join(testCaseDir, "bootstrap.yaml")
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-app-id",
			Name:    "TestApp",
			Version: "1.0.0",
			Env:     "test",
		}),
	)
	require.NoError(t, err, "Failed to initialize runtime")
	defer rtInstance.Cleanup()

	// 2. Verify core components are properly initialized
	asserts.NotNil(rtInstance.Logger())
	asserts.Equal("test-app-id", rtInstance.AppInfo().ID)
	asserts.Equal("TestApp", rtInstance.AppInfo().Name)
	asserts.Equal("1.0.0", rtInstance.AppInfo().Version)

	// 3. Get ConfigDecoder from runtime
	decoder := rtInstance.Config()
	asserts.NotNil(decoder, "ConfigDecoder should not be nil")

	// 4. Verify we can decode the entire config into a map
	var configMap map[string]interface{}
	err = decoder.Decode("", &configMap)
	asserts.NoError(err, "Failed to decode config into map")
	asserts.NotEmpty(configMap, "Decoded config map should not be empty")
}
