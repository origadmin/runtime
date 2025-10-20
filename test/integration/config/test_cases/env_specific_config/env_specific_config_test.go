package env_specific_config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	_ "github.com/origadmin/runtime/config/envsource"
	_ "github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/test/helper"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// EnvSpecificConfigTestSuite tests environment-specific configuration loading.
type EnvSpecificConfigTestSuite struct {
	suite.Suite
}

func TestEnvSpecificConfigTestSuite(t *testing.T) {
	suite.Run(t, new(EnvSpecificConfigTestSuite))
}

// TestEnvSpecificLoading verifies that the runtime correctly loads environment-specific configurations.
func (s *EnvSpecificConfigTestSuite) TestEnvSpecificLoading() {
	t := s.T()
	assert := assert.New(t)
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	bootstrapPath := "config/test_cases/env_specific_config/bootstrap_env.yaml"

	// Test cases for different environments
	testCases := []struct {
		name                string
		envVar              string
		expectedAppID       string
		expectedAppName     string
		expectedAppEnv      string
		expectedLoggerLevel string
		expectedMetadataKey string
	}{
		{
			name:                "DevEnvironment",
			envVar:              "dev",
			expectedAppID:       "base-app-id", // Not overridden by config_dev.yaml
			expectedAppName:     "DevApp",      // Overridden by config_dev.yaml
			expectedAppEnv:      "dev",         // Overridden by config_dev.yaml
			expectedLoggerLevel: "info",        // Overridden by config_dev.yaml
			expectedMetadataKey: "dev_value",
		},
		{
			name:                "ProdEnvironment",
			envVar:              "prod",
			expectedAppID:       "base-app-id", // Not overridden by config_prod.yaml
			expectedAppName:     "ProdApp",     // Overridden by config_prod.yaml
			expectedAppEnv:      "prod",        // Overridden by config_prod.yaml
			expectedLoggerLevel: "error",       // Overridden by config_prod.yaml
			expectedMetadataKey: "prod_value",
		},
		{
			name:                "DefaultEnvironment", // No specific env config, should load base
			envVar:              "nonexistent",
			expectedAppID:       "base-app-id",
			expectedAppName:     "BaseApp",
			expectedAppEnv:      "default",
			expectedLoggerLevel: "debug",
			expectedMetadataKey: "", // Should not exist
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set APP_ENV environment variable
			os.Setenv("APP_ENV", tc.envVar)
			defer os.Unsetenv("APP_ENV") // Clean up env var after test

			rtInstance, rtCleanup, err := rt.NewFromBootstrap(
				bootstrapPath,
				bootstrap.WithAppInfo(&interfaces.AppInfo{
					ID:      "env-test-app",
					Name:    "EnvTestApp",
					Version: "1.0.0",
				}),
			)
			assert.NoError(err, "Failed to initialize runtime from bootstrap for %s: %v", tc.name, err)
			defer rtCleanup()

			configDecoder := rtInstance.Config()
			assert.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil for %s", tc.name)

			var cfg testconfigs.TestConfig
			err = configDecoder.Decode("app", &cfg.App) // Decode only the app section for specific assertions
			assert.NoError(err, "Failed to decode app config for %s: %v", tc.name, err)

			err = configDecoder.Decode("logger", &cfg.Logger) // Decode only the logger section
			assert.NoError(err, "Failed to decode logger config for %s: %v", tc.name, err)

			// Assertions for environment-specific overrides
			assert.Equal(tc.expectedAppID, cfg.App.Id, "App ID mismatch for %s", tc.name)
			assert.Equal(tc.expectedAppName, cfg.App.Name, "App Name mismatch for %s", tc.name)
			assert.Equal(tc.expectedAppEnv, cfg.App.Env, "App Env mismatch for %s", tc.name)
			assert.Equal(tc.expectedLoggerLevel, cfg.Logger.Level, "Logger Level mismatch for %s", tc.name)

			if tc.expectedMetadataKey != "" {
				assert.Contains(cfg.App.Metadata, "env_specific_key", "Metadata should contain env_specific_key for %s", tc.name)
				assert.Equal(tc.expectedMetadataKey, cfg.App.Metadata["env_specific_key"], "Metadata value mismatch for %s", tc.name)
			} else {
				assert.NotContains(cfg.App.Metadata, "env_specific_key", "Metadata should not contain env_specific_key for %s", tc.name)
			}

			t.Logf("Environment-specific config for %s loaded and verified successfully!", tc.name)
		})
	}
}
