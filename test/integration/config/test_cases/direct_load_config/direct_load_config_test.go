package direct_load_config_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/origadmin/runtime/test/helper"
	parentconfig "github.com/origadmin/runtime/test/integration/config" // Import the parent package for AssertTestConfig
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

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
	t := s.T()
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	// Define test cases for different formats of the unified config.
	// For now, we only have YAML, but this structure allows easy expansion.
	testCases := []struct {
		name     string
		filePath string
	}{
		{name: "YAML", filePath: "test/integration/config/test_cases/direct_load_config/config.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cfg testconfigs.TestConfig
			helper.LoadConfigFromFile(t, tc.filePath, &cfg)
			parentconfig.AssertTestConfig(t, &cfg)
			t.Logf("%s unified config loaded and verified successfully!", tc.name)
		})
	}
}
