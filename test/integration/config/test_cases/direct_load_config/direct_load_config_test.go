package direct_load_config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/origadmin/runtime/test/helper"
	parentconfig "github.com/origadmin/runtime/test/integration/config"
	"github.com/origadmin/runtime/test/integration/config/builders"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

var defaultConfig *testconfigs.TestConfig

func init() {
	defaultConfig = &testconfigs.TestConfig{
		App:                       builders.NewDefaultApp(),
		Servers:                   builders.NewDefaultServers(),
		Client:                    builders.NewDefaultClient(),
		Logger:                    builders.NewDefaultLogger(),
		Discoveries:               builders.NewDefaultDiscoveries(),
		RegistrationDiscoveryName: "internal-consul",
		Trace:                     builders.NewDefaultTrace(),
		Middlewares:               builders.NewDefaultMiddlewares(),
	}
}

// DirectLoadConfigTestSuite tests the direct loading of the unified config.yaml.
type DirectLoadConfigTestSuite struct {
	suite.Suite
}

func TestDirectLoadConfigTestSuite(t *testing.T) {
	suite.Run(t, new(DirectLoadConfigTestSuite))
}

// TestDirectConfigLoading verifies that the raw config.yaml file is well-formed and parsable
// into the unified TestConfig struct.
func (s *DirectLoadConfigTestSuite) TestDirectConfigLoading() {
	// Define test cases for different formats of the unified config.
	testCases := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: "config.yaml"},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			// If the config file does not exist, create it from a default TestConfig struct.
			// This ensures the test has a valid config and provides a "live" template for developers.
			if _, err := os.Stat(tc.filePath); os.IsNotExist(err) {
				helper.SaveConfigToFileWithViper(t, defaultConfig, tc.filePath, tc.name)
			}

			var cfg testconfigs.TestConfig
			helper.LoadConfigFromFile(t, tc.filePath, &cfg)

			// Assert that the loaded config matches the default config by performing a detailed, field-by-field assertion.
			parentconfig.AssertTestConfig(t, defaultConfig, &cfg)

			t.Logf("%s unified config loaded and verified successfully!", tc.name)
		})
	}
}
