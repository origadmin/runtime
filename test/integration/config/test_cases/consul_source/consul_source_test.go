package consul_source_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/test/helper"
	parentconfig "github.com/origadmin/runtime/test/integration/config" // Import the parent package for AssertTestConfig
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// ConsulSourceTestSuite tests the Consul configuration source integration.
type ConsulSourceTestSuite struct {
	suite.Suite
}

func TestConsulSourceTestSuite(t *testing.T) {
	suite.Run(t, new(ConsulSourceTestSuite))
}

// TestConsulSourceLoading verifies that configuration can be loaded correctly from a mock Consul source.
func (s *ConsulSourceTestSuite) TestConsulSourceLoading() {
	t := s.T()
	assert := assert.New(t)
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	// Content that would typically come from Consul
	consulConfigContent := `
app:
  id: "consul-app-id"
  name: "ConsulApp"
  version: "1.0.0"
  env: "consul-test"

logger:
  level: "warn"
  format: "json"
`
	// Create a mock Consul source with the predefined content
	mockConsulData := map[string]string{
		"config/test/app_config": consulConfigContent, // Key matches the path in bootstrap_consul.yaml
	}
	mockSource := helper.NewMockConsulSource(mockConsulData)

	bootstrapPath := "test/integration/config/test_cases/consul_source/bootstrap_consul.yaml"

	// Initialize Runtime, injecting the mock Consul source
	rtInstance, rtCleanup, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "consul-test-app",
			Name:    "ConsulTestApp",
			Version: "1.0.0",
		}),
		// This is a placeholder. The actual mechanism to inject a custom source
		// might vary (e.g., a custom bootstrap.Source implementation or a specific
		// option in the runtime package to register mock sources).
		// For now, we assume bootstrap.WithSource can take a custom source.
		// If not, the runtime.NewFromBootstrap or bootstrap.New function needs to be adapted.
		bootstrap.WithSource("consul", mockSource), // Assuming "consul" is the type name used in bootstrap.yaml
	)
	assert.NoError(err, "Failed to initialize runtime from bootstrap with mock Consul source: %v", err)
	defer rtCleanup()

	configDecoder := rtInstance.Config()
	assert.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil")

	var cfg testconfigs.TestConfig
	err = configDecoder.Decode("app", &cfg.App)
	assert.NoError(err, "Failed to decode app config from runtime: %v", err)

	err = configDecoder.Decode("logger", &cfg.Logger)
	assert.NoError(err, "Failed to decode logger config from runtime: %v", err)

	// Assertions based on config-in-consul.yaml content
	assert.NotNil(cfg.App)
	assert.Equal("consul-app-id", cfg.App.Id)
	assert.Equal("ConsulApp", cfg.App.Name)
	assert.Equal("1.0.0", cfg.App.Version)
	assert.Equal("consul-test", cfg.App.Env)

	assert.NotNil(cfg.Logger)
	assert.Equal("warn", cfg.Logger.Level)
	assert.Equal("json", cfg.Logger.Format)

	t.Logf("Consul source config loaded and verified successfully!")
}
