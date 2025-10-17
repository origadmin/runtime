package bootstrap_load_config_test

import (
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
		{name: "YAML", filePath: "test/integration/config/test_cases/bootstrap_load_config/bootstrap.yaml"},
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

			// Run assertions on the loaded configuration.
			parentconfig.AssertTestConfig(t, &cfg)
			t.Logf("Runtime loaded and verified %s bootstrap config successfully!", tc.name)
		})
	}
}
