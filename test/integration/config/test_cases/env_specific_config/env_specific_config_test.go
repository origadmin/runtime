package env_specific_config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/runtime/bootstrap"
	_ "github.com/origadmin/runtime/config/envsource"
	_ "github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/interfaces"
	parentconfig "github.com/origadmin/runtime/test/integration/config"
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

	bootstrapPath := "bootstrap_env.yaml"

	// Test cases for different environments
	testCases := []struct {
		name           string
		envVar         string
		expectedApp    *appv1.App
		expectedLogger *loggerv1.Logger
	}{
		{
			name:   "DevEnvironment",
			envVar: "dev",
			expectedApp: &appv1.App{
				Id:       "base-app-id",
				Name:     "DevApp",
				Version:  "1.0.0",
				Env:      "dev",
				Metadata: map[string]string{"env_specific_key": "dev_value"},
			},
			expectedLogger: &loggerv1.Logger{
				Level:  "info",
				Format: "text",
			},
		},
		{
			name:   "ProdEnvironment",
			envVar: "prod",
			expectedApp: &appv1.App{
				Id:       "base-app-id",
				Name:     "ProdApp",
				Version:  "1.0.0",
				Env:      "prod",
				Metadata: map[string]string{"env_specific_key": "prod_value"},
			},
			expectedLogger: &loggerv1.Logger{
				Level:  "error",
				Format: "text",
			},
		},
		{
			name:   "DefaultEnvironment",
			envVar: "nonexistent",
			expectedApp: &appv1.App{
				Id:      "base-app-id",
				Name:    "BaseApp",
				Version: "1.0.0",
				Env:     "default",
			},
			expectedLogger: &loggerv1.Logger{
				Level:  "debug",
				Format: "text",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set APP_ENV environment variable
			os.Setenv("APP_ENV", tc.envVar)
			defer os.Unsetenv("APP_ENV") // Clean up env var after test
			t.Logf("Setting APP_ENV to %s for %s", tc.envVar, tc.name)

			rtInstance, err := rt.NewFromBootstrap(
				bootstrapPath,
				bootstrap.WithAppInfo(&interfaces.AppInfo{
					ID:      "env-test-app",
					Name:    "EnvTestApp",
					Version: "1.0.0",
				}),
			)
			require.NoError(t, err, "Failed to initialize runtime from bootstrap for %s", tc.name)
			defer rtInstance.Cleanup()

			// Decode the app and logger sections
			actualApp, err := rtInstance.StructuredConfig().DecodeApp()
			require.NoError(t, err, "Failed to decode app config for %s", tc.name)

			actualLogger, err := rtInstance.StructuredConfig().DecodeLogger()
			require.NoError(t, err, "Failed to decode logger config for %s", tc.name)

			// Perform assertions using the modular assertion toolkit.
			parentconfig.AssertAppConfig(t, tc.expectedApp, actualApp)
			parentconfig.AssertLoggerConfig(t, tc.expectedLogger, actualLogger)

			t.Logf("Environment-specific config for %s loaded and verified successfully!", tc.name)
		})
	}
}
