package config_test

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/stretchr/testify/assert"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
)

type CustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"endpoints"`
}

// TestConfigSourceMergingAndPriority verifies that the configuration manager correctly
// merges settings from multiple sources, respecting their defined priorities.
// It uses a bootstrap file that loads a base config and a higher-priority override config.
func (s *ConfigTestSuite) TestConfigSourceMergingAndPriority() {
	t := s.T()
	assert := assert.New(t)

	// Setup: Change CWD to the project root to resolve config paths correctly.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file info")
	}
	currentDir := filepath.Dir(filename)
	runtimeRoot := filepath.Join(currentDir, "../../..")

	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalCwd) }()

	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	// Path to the bootstrap file that defines the multi-source configuration.
	bootstrapPath := "test/integration/config/configs/bootstrap.yaml"

	// 1. Initialize Runtime from the bootstrap file.
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-source-merging",
			Name:    "TestSourceMerging",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer rtInstance.Cleanup()

	// 2. Get the configuration decoder from the runtime.
	configDecoder := rtInstance.Config()
	assert.NotNil(configDecoder)

	// 3. Decode and verify the logger configuration.
	// We expect the values from 'test_config.yaml' (priority 200) to override 'config.yaml' (priority 100).
	var loggerConfig struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	}
	err = configDecoder.Decode("logger", &loggerConfig)
	assert.NoError(err)
	assert.Equal("debug", loggerConfig.Level, "Logger level should be overridden by the higher-priority source")
	assert.Equal("text", loggerConfig.Format, "Logger format should be overridden by the higher-priority source")

	// 4. Decode and verify a section that only exists in the base config ('config.yaml').
	// This ensures that the base configuration is not completely discarded.
	var customSettings CustomSettings // Re-using struct from config_test.go
	err = configDecoder.Decode("components.my-custom-settings", &customSettings)
	assert.NoError(err)
	assert.True(customSettings.FeatureEnabled, "Custom settings from the base config should still be present")
	assert.Equal("super-secret-key-123", customSettings.APIKey, "Custom settings from the base config should still be present")
}
