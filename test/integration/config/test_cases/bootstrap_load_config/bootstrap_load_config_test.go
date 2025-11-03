package bootstrap_load_config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	parentconfig "github.com/origadmin/runtime/test/integration/config"
	"github.com/origadmin/runtime/test/integration/config/builders"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

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

	bootstrapPath := "testdata/complete_config/bootstrap.yaml"

	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-complete-config",
			Name:    "TestCompleteConfig",
			Version: "1.0.0",
		}),
	)
	require.NoError(t, err, "Failed to initialize runtime")
	defer rtInstance.Cleanup()

	var actualConfig testconfigs.TestConfig
	err = rtInstance.Config().Decode("", &actualConfig)
	require.NoError(t, err)

	// Build the expected config to EXACTLY match the content of complete_config/config.yaml
	expectedServers := builders.NewDefaultServers()
	expectedServers.Configs[0].Name = "grpc_servers"
	expectedServers.Configs[0].Protocol = "" // Not present in YAML, so it should be the zero value
	expectedServers.Configs[1].Name = "http_servers"
	expectedServers.Configs[1].Protocol = "" // Not present in YAML, so it should be the zero value

	expectedConfig := &testconfigs.TestConfig{
		App:         builders.NewDefaultApp(),
		Servers:     expectedServers,
		Client:      builders.NewDefaultClient(),
		Logger:      builders.NewDefaultLogger(),
		Trace:       builders.NewDefaultTrace(),
		Middlewares: builders.NewDefaultMiddlewares(),
		// These fields are not in the config file, so they should be nil/zero
		Discoveries:               nil,
		RegistrationDiscoveryName: "",
	}

	// Use the shared, robust assertion logic
	parentconfig.AssertTestConfig(t, expectedConfig, &actualConfig)
	t.Logf("Runtime loaded and verified complete config successfully!")
}

// TestConfigProtoIntegration tests integration between configuration and Protocol Buffers
func (s *RuntimeIntegrationTestSuite) TestConfigProtoIntegration() {
	t := s.T()

	bootstrapPath := "testdata/proto_integration/bootstrap.yaml"

	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-proto-config",
			Name:    "TestProtoConfig",
			Version: "1.0.0",
		}),
	)
	require.NoError(t, err, "Failed to initialize runtime")
	defer rtInstance.Cleanup()

	var actualConfig testconfigs.TestConfig
	err = rtInstance.Config().Decode("", &actualConfig)
	require.NoError(t, err)

	// Build the expected config based on proto_integration/config.yaml
	expectedApp := &appv1.App{
		Id:      "test-app-id",
		Name:    "TestApp",
		Version: "1.0.0",
		// Env and Metadata are not in the YAML, so they should be zero values.
		Env:      "",
		Metadata: nil,
	}

	expectedServers := builders.NewDefaultServers()
	expectedServers.Configs[0].Name = "grpc_servers"
	expectedServers.Configs[0].Protocol = "" // Not present in YAML
	expectedServers.Configs[1].Name = "http_servers"
	expectedServers.Configs[1].Protocol = "" // Not present in YAML

	// Assertions for the loaded sections
	parentconfig.AssertAppConfig(t, expectedApp, actualConfig.App)
	parentconfig.AssertServersConfig(t, expectedServers, actualConfig.Servers)
	require.Equal(t, "test-discovery", actualConfig.RegistrationDiscoveryName)

	// Assertions for sections that should NOT be loaded
	require.Nil(t, actualConfig.Client, "Client config should be nil")
	require.Nil(t, actualConfig.Logger, "Logger config should be nil")
	require.Nil(t, actualConfig.Discoveries, "Discoveries config should be nil")
	require.Nil(t, actualConfig.Trace, "Trace config should be nil")
	require.Nil(t, actualConfig.Middlewares, "Middlewares config should be nil")
}

// TestRuntimeDecoder verifies that the configuration decoder is properly exposed by runtime
// and can be used to parse custom configuration sections
func (s *RuntimeIntegrationTestSuite) TestRuntimeDecoder() {
	t := s.T()

	bootstrapPath := "testdata/decoder_test/bootstrap.yaml"

	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-decoder",
			Name:    "TestDecoder",
			Version: "1.0.0",
			Env:     "test",
		}),
	)
	require.NoError(t, err, "Failed to initialize runtime")
	defer rtInstance.Cleanup()

	decoder := rtInstance.Config()
	require.NotNil(t, decoder, "ConfigDecoder should not be nil")

	var customSettings TestCustomSettings
	err = decoder.Decode("components.my-custom-settings", &customSettings)
	require.NoError(t, err, "Failed to decode custom settings")

	require.True(t, customSettings.FeatureEnabled, "Feature should be enabled")
	require.Equal(t, 100, customSettings.RateLimit, "Rate limit should be 100")
	require.Len(t, customSettings.Endpoints, 2, "Should have 2 endpoints")
	require.Equal(t, "users", customSettings.Endpoints[0].Name, "First endpoint should be 'users'")
	require.Equal(t, "/api/v1/users", customSettings.Endpoints[0].Path, "Users endpoint path should be correct")
	require.Equal(t, "products", customSettings.Endpoints[1].Name, "Second endpoint should be 'products'")
	require.Equal(t, "/api/v1/products", customSettings.Endpoints[1].Path, "Products endpoint path should be correct")
}
