package source_manager

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	_ "github.com/origadmin/runtime/config/envsource"
	_ "github.com/origadmin/runtime/config/file"
)

// SourceManagerTestSuite defines the test suite for configuration source manager
type SourceManagerTestSuite struct {
	suite.Suite
}

// TestSourceManagerSuite runs the test suite
func TestSourceManagerSuite(t *testing.T) {
	suite.Run(t, new(SourceManagerTestSuite))
}

type CustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []struct {
		Path string `json:"path"`
	} `json:"endpoints"`
}

// TestConfigSourceMergingAndPriority verifies that the configuration manager correctly
// merges settings from multiple sources, respecting their defined priorities.
// It uses a bootstrap file that loads a base config and a higher-priority override config.
func (s *SourceManagerTestSuite) TestConfigSourceMergingAndPriority() {
	t := s.T()

	// Create AppInfo using the new functional options pattern
	appInfo := rt.NewAppInfo(
		"test-app",
		"1.0.0",
		rt.WithAppInfoID("test-app"),
	)

	// Use a path relative to the test file itself. This is the robust way to handle test data
	// and is independent of the current working directory.
	bootstrapPath := filepath.Join("testdata", "merging_and_priority", "bootstrap.yaml")
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		rt.WithAppInfo(appInfo), // Pass the created AppInfo
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	// Removed defer rtInstance.Cleanup() as it's no longer available

	// 2. Get the configuration decoder from the runtime
	configDecoder := rtInstance.Config()
	s.NotNil(configDecoder)

	// 3. Decode and verify the logger configuration
	// The values should be properly merged from all sources based on their priorities
	var loggerConfig struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	}
	err = configDecoder.Decode("logger", &loggerConfig)
	s.NoError(err)
	s.Equal("debug", loggerConfig.Level, "Logger level should be overridden by the higher-priority source")
	s.Equal("text", loggerConfig.Format, "Logger format should be overridden by the higher-priority source")

	// 4. Decode and verify a section from the configuration
	// This ensures that the configuration is properly loaded
	var customSettings CustomSettings // Re-using struct from config_test.go
	err = configDecoder.Decode("components.my-custom-settings", &customSettings)
	s.NoError(err)
	s.True(customSettings.FeatureEnabled, "Custom settings from the base config should still be present")
	s.Equal("super-secret-key-123", customSettings.APIKey, "Custom settings from the base config should still be present")
}
