package bootstrap_load_config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/test/helper"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// assertBootstrapTestConfig contains the validation logic for the TestConfig struct
// specific to the bootstrap_load_config test case, omitting discovery assertions.
func assertBootstrapTestConfig(t *testing.T, cfg *testconfigs.TestConfig) {
	assert := assert.New(t)

	// App configuration assertions
	assert.NotNil(cfg.App)
	assert.Equal("test-app-id", cfg.App.GetId())
	assert.Equal("TestApp", cfg.App.GetName())
	assert.Equal("1.0.0", cfg.App.GetVersion())
	assert.Equal("test", cfg.App.GetEnv())
	assert.Contains(cfg.App.GetMetadata(), "key1")
	assert.Contains(cfg.App.GetMetadata(), "key2")
	assert.Equal("value1", cfg.App.GetMetadata()["key1"])
	assert.Equal("value2", cfg.App.GetMetadata()["key2"])

	// Server configuration assertions
	assert.Len(cfg.GrpcServers, 1)
	assert.Equal("tcp", cfg.GrpcServers[0].GetNetwork())
	assert.Equal(":9000", cfg.GrpcServers[0].GetAddr())
	assert.Equal("1s", cfg.GrpcServers[0].GetTimeout().AsDuration().String())

	assert.Len(cfg.HttpServers, 1)
	assert.Equal("tcp", cfg.HttpServers[0].GetNetwork())
	assert.Equal(":8000", cfg.HttpServers[0].GetAddr())
	assert.Equal("2s", cfg.HttpServers[0].GetTimeout().AsDuration().String())

	// Client configuration assertions
	assert.NotNil(cfg.Client)
	assert.Equal("discovery:///user-service", cfg.Client.GetEndpoint())
	assert.Equal("3s", cfg.Client.GetTimeout().AsDuration().String())
	assert.NotNil(cfg.Client.GetSelector())
	assert.Equal("v1.0.0", cfg.Client.GetSelector().GetVersion())

	// Discovery configuration assertions (omitted for this test case)
	assert.Len(cfg.Discoveries, 0, "Discoveries should be empty for this test case")
	assert.Empty(cfg.GetRegistrationDiscoveryName(), "RegistrationDiscoveryName should be empty for this test case")

	// Logger configuration assertions
	assert.NotNil(cfg.Logger)
	assert.Equal("info", cfg.Logger.GetLevel())
	assert.Equal("json", cfg.Logger.GetFormat())
	assert.True(cfg.Logger.GetStdout())

	// Tracer configuration assertions
	assert.NotNil(cfg.Tracer)
	assert.Equal("jaeger", cfg.Tracer.GetName())
	assert.Equal("http://jaeger:14268/api/traces", cfg.Tracer.GetEndpoint())

	// Middleware configuration assertions
	assert.NotNil(cfg.Middlewares)
	assert.Len(cfg.Middlewares.GetMiddlewares(), 1, "Should have 1 middleware configured")

	corsMiddleware := cfg.Middlewares.GetMiddlewares()[0]
	assert.Equal("cors", corsMiddleware.GetType())
	assert.True(corsMiddleware.GetEnabled())
	assert.NotNil(corsMiddleware.GetCors(), "CORS config should not be nil for middleware of type cors")
	assert.Len(corsMiddleware.GetCors().GetAllowedOrigins(), 1)
	assert.Equal("*", corsMiddleware.GetCors().GetAllowedOrigins()[0])
	assert.Len(corsMiddleware.GetCors().GetAllowedMethods(), 2)
	assert.Equal("GET", corsMiddleware.GetCors().GetAllowedMethods()[0])
	assert.Equal("POST", corsMiddleware.GetCors().GetAllowedMethods()[1])
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
	assert := assert.New(t)
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	// Define test cases for different formats of the bootstrap config.
	// For now, we only have YAML, but this structure allows easy expansion.
	testCases := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: "config/test_cases/bootstrap_load_config/bootstrap.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Runtime from the bootstrap file.
			// The AppInfo here is minimal as the main config is loaded via bootstrap.
			rtInstance, rtCleanup, err := rt.NewFromBootstrap(
				tc.filePath,
				bootstrap.WithAppInfo(&interfaces.AppInfo{
					ID:      "bootstrap-test-app",
					Name:    "BootstrapTestApp",
					Version: "1.0.0",
				}),
			)
			assert.NoError(err, "Failed to initialize runtime from bootstrap: %v", err)
			defer rtCleanup()

			// Get the configuration decoder from the runtime.
			configDecoder := rtInstance.Config()
			assert.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil")

			// Decode the entire configuration into our unified TestConfig struct.
			var cfg testconfigs.TestConfig
			err = configDecoder.Decode("", &cfg) // Decode the root of the config
			assert.NoError(err, "Failed to decode config from runtime: %v", err)

			// Run assertions on the loaded configuration using the custom assertion.
			assertBootstrapTestConfig(t, &cfg)
			t.Logf("Runtime loaded and verified %s bootstrap config successfully!", tc.name)
		})
	}
}
