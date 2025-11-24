package middleware

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	rt "github.com/origadmin/runtime"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/config/middleware/selector/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
)

// Define all supported middleware types
var supportedMiddlewareTypes = map[string]struct{}{
	"recovery":        {},
	"metadata":        {},
	"logging":         {},
	"selector":        {},
	"rate_limiter":    {},
	"circuit_breaker": {},
	"jwt":             {},
	"cors":            {},
	"metrics":         {},
	"validator":       {},
}

func TestMiddleware(t *testing.T) {
	t.Run("LoadAndBuild", TestMiddleware_LoadAndBuild)
	t.Run("Selector", TestSelectorMiddleware)
	t.Run("Creation", TestMiddleware_Creation)
	t.Run("EdgeCases", TestMiddleware_EdgeCases)
}

// getMiddlewareField gets middleware configuration field
func getMiddlewareField(mw *middlewarev1.Middleware, fieldName string) interface{} {
	switch fieldName {
	case "metadata":
		return mw.Metadata
	case "logging":
		return mw.Logging
	case "selector":
		return mw.Selector
	case "rate_limiter":
		return mw.RateLimiter
	case "circuit_breaker":
		return mw.CircuitBreaker
	default:
		return nil
	}
}

// setupRuntime initializes the runtime for middleware tests
func setupRuntime(t *testing.T, configFilePath string) (*rt.App, *middlewarev1.Middlewares) {
	bootstrapPath := filepath.Join(filepath.Dir(configFilePath), "bootstrap.yaml")

	appInfo := rt.NewAppInfo("middleware-test-app", "1.0.0", rt.WithAppInfoID("middleware-test-app"))
	rtInstance, err := rt.NewFromBootstrap(bootstrapPath, rt.WithAppInfo(appInfo))
	require.NoError(t, err, "Failed to initialize runtime")

	var configs middlewarev1.Middlewares
	err = rtInstance.Config().Decode("middlewares", &configs)
	require.NoError(t, err, "Failed to decode middlewares config")

	return rtInstance, &configs
}

// TestMiddleware_LoadAndBuild tests middleware loading and building
func TestMiddleware_LoadAndBuild(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	configFilePath := filepath.Join(filepath.Dir(filename), "configs", "config.yaml")

	rtInstance, configs := setupRuntime(t, configFilePath)

	t.Run("VerifyConfig", func(t *testing.T) {
		// Check the number of middlewares in the configuration
		assert.GreaterOrEqual(t, len(configs.Configs), 2, "Should have at least 2 middlewares in config")

		// Check configuration of each middleware
		for _, mw := range configs.Configs {
			t.Run(mw.Name, func(t *testing.T) {
				// Verify required fields
				assert.NotEmpty(t, mw.Name, "Middleware name should not be empty")
				assert.True(t, mw.Enabled, "Middleware should be enabled by default")

				// Verify if the type is supported
				_, exists := supportedMiddlewareTypes[mw.Type]
				assert.True(t, exists, "Unsupported middleware type: %s", mw.Type)

				// Verify if the configuration field exists
				configValue := getMiddlewareField(mw, mw.Type)
				assert.NotNil(t, configValue, "Config for %s should not be nil", mw.Type)
			})
		}
	})

	t.Run("ClientMiddlewares", func(t *testing.T) {
		// Get client middleware map
		middlewareProvider, err := rtInstance.Container().Middleware()
		require.NoError(t, err)
		clientMWsMap, err := middlewareProvider.ClientMiddlewares()
		require.NoError(t, err)

		// Calculate the actual number of supported client middlewares (excluding unsupported middleware types)
		expectedCount := 0
		for _, mw := range configs.Configs {
			if mw.Enabled && mw.Type != "rate_limiter" { // rate_limiter 在客户端不受支持
				t.Logf("Including %s middleware for client", mw.Type)
				expectedCount++
			} else {
				t.Logf("Skipping rate_limiter middleware for client")
			}
		}

		// Verify the number of middlewares
		assert.Equal(t, expectedCount, len(clientMWsMap), "Number of client middlewares should match expected")
	})

	t.Run("ServerMiddlewares", func(t *testing.T) {
		// Get server middleware map
		middlewareProvider, err := rtInstance.Container().Middleware()
		require.NoError(t, err)
		serverMWsMap, err := middlewareProvider.ServerMiddlewares()
		require.NoError(t, err)

		// Calculate the actual number of supported server middlewares (excluding unsupported middleware types)
		expectedCount := 0
		for _, mw := range configs.Configs {
			if mw.Enabled && mw.Type != "circuit_breaker" { // circuit_breaker 在服务端不受支持
				t.Logf("Including %s middleware for server", mw.Type)
				expectedCount++
			} else {
				t.Logf("Skipping circuit_breaker middleware for server")
			}
		}

		// Verify the number of middlewares
		assert.Equal(t, expectedCount, len(serverMWsMap), "Number of server middlewares should match expected")
	})
}

// TestSelectorMiddleware tests the includes/excludes functionality of Selector middleware
func TestSelectorMiddleware(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	configFilePath := filepath.Join(filepath.Dir(filename), "configs", "config.yaml")

	rtInstance, _ := setupRuntime(t, configFilePath)
	_, err := rtInstance.Container().Middleware()
	require.NoError(t, err, "Failed to get middleware provider")

	// Prepare test data
	tests := []struct {
		name        string
		config      *middlewarev1.Middlewares
		expectCount int
		side        string // "client" or "server"
	}{
		{
			name: "selector_with_includes",
			config: &middlewarev1.Middlewares{
				Configs: []*middlewarev1.Middleware{
					{
						Name:    "metadata",
						Type:    "metadata",
						Enabled: true,
					},
					{
						Name:    "logging",
						Type:    "logging",
						Enabled: true,
					},
					{
						Name:    "selector",
						Type:    "selector",
						Enabled: true,
						Selector: &selectorv1.Selector{
							Includes: []string{"logging"}, // Only include logging middleware
						},
					},
				},
			},
			expectCount: 2, // selector + logging
			side:        "client",
		},
		{
			name: "selector_with_excludes",
			config: &middlewarev1.Middlewares{
				Configs: []*middlewarev1.Middleware{
					{
						Name:    "metadata",
						Type:    "metadata",
						Enabled: true,
					},
					{
						Name:    "logging",
						Type:    "logging",
						Enabled: true,
					},
					{
						Name:    "selector",
						Type:    "selector",
						Enabled: true,
						Selector: &selectorv1.Selector{
							Excludes: []string{"metadata"}, // Exclude metadata middleware
						},
					},
				},
			},
			expectCount: 2, // selector + logging
			side:        "server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempBootstrapPath := filepath.Join(tempDir, "bootstrap.yaml")
			tempConfigPath := filepath.Join(tempDir, "config.yaml")

			tempConfigs := &middlewarev1.Middlewares{Configs: []*middlewarev1.Middleware{tt.config.Configs[0]}}
			configBytes, err := yaml.Marshal(tempConfigs)
			require.NoError(t, err, "Failed to marshal middleware config to YAML")
			err = os.WriteFile(tempConfigPath, configBytes, 0644)
			require.NoError(t, err, "Failed to write config file")

			bootstrapContent := []byte(`sources:
  - kind: file
    path: ` + filepath.Base(tempConfigPath) + `
    format: yaml`)
			err = os.WriteFile(tempBootstrapPath, bootstrapContent, 0644)
			require.NoError(t, err, "Failed to write bootstrap file")

			tempAppInfo := rt.NewAppInfo("temp-middleware-test", "1.0.0", rt.WithAppInfoID("temp-middleware-test"))
			tempRtInstance, err := rt.NewFromBootstrap(tempBootstrapPath, rt.WithAppInfo(tempAppInfo))
			require.NoError(t, err)

			switch tt.side {
			case "client":
				middlewareProvider, err := tempRtInstance.Container().Middleware()
				require.NoError(t, err)
				mwMap, err := middlewareProvider.ClientMiddlewares()
				require.NoError(t, err)
				assert.Len(t, mwMap, tt.expectCount, "Unexpected number of client middlewares")
			case "server":
				middlewareProvider, err := tempRtInstance.Container().Middleware()
				require.NoError(t, err)
				mwMap, err := middlewareProvider.ServerMiddlewares()
				require.NoError(t, err)
				assert.Len(t, mwMap, tt.expectCount, "Unexpected number of server middlewares")
			}
		})
	}
}

// TestMiddleware_Creation tests middleware creation
func TestMiddleware_Creation(t *testing.T) {
	//_, filename, _, _ := runtime.Caller(0)
	//configFilePath := filepath.Join(filepath.Dir(filename), "configs", "config.yaml")

	//rtInstance, _ := setupRuntime(t, configFilePath)

	// Test creation of different types of middleware
	tests := []struct {
		name     string
		config   *middlewarev1.Middleware
		exists   bool
		clientMW bool
		serverMW bool
	}{
		{
			name: "valid logging middleware",
			config: &middlewarev1.Middleware{
				Name:    "test-logging",
				Type:    "logging",
				Enabled: true,
				Logging: &middlewarev1.Logging{},
			},
			exists:   true,
			clientMW: true,
			serverMW: true,
		},
		{
			name: "disabled middleware",
			config: &middlewarev1.Middleware{
				Name:    "disabled-mw",
				Type:    "logging",
				Enabled: false,
			},
			exists:   false,
			clientMW: false,
			serverMW: false,
		},
		{
			name: "unknown middleware type",
			config: &middlewarev1.Middleware{
				Name:    "unknown",
				Type:    "nonexistent",
				Enabled: true,
			},
			exists:   false,
			clientMW: false,
			serverMW: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tempBootstrapPath := filepath.Join(tempDir, "bootstrap.yaml")
			tempConfigPath := filepath.Join(tempDir, "config.yaml")

			tempConfigs := &middlewarev1.Middlewares{Configs: []*middlewarev1.Middleware{tt.config}}
			configBytes, err := yaml.Marshal(tempConfigs)
			require.NoError(t, err, "Failed to marshal middleware config to YAML")
			err = os.WriteFile(tempConfigPath, configBytes, 0644)
			require.NoError(t, err, "Failed to write config file")

			bootstrapContent := []byte(`sources:
  - kind: file
    path: ` + filepath.Base(tempConfigPath) + `
    format: yaml`)
			err = os.WriteFile(tempBootstrapPath, bootstrapContent, 0644)
			require.NoError(t, err, "Failed to write bootstrap file")

			tempAppInfo := rt.NewAppInfo("temp-single-mw-test", "1.0.0", rt.WithAppInfoID("temp-single-mw-test"))
			tempRtInstance, err := rt.NewFromBootstrap(tempBootstrapPath, rt.WithAppInfo(tempAppInfo))
			require.NoError(t, err)

			middlewareProvider, err := tempRtInstance.Container().Middleware()
			require.NoError(t, err)

			// Test client middleware
			clientMW, err := middlewareProvider.ClientMiddleware(tt.config.Name)
			if tt.clientMW {
				assert.NoError(t, err)
				assert.NotNil(t, clientMW, "Expected client middleware to be created")
			} else {
				assert.Error(t, err) // Expect an error if middleware not found/created
				assert.Nil(t, clientMW, "Expected no client middleware to be created")
			}

			// Test server middleware
			serverMW, err := middlewareProvider.ServerMiddleware(tt.config.Name)
			if tt.serverMW {
				assert.NoError(t, err)
				assert.NotNil(t, serverMW, "Expected server middleware to be created")
			} else {
				assert.Error(t, err) // Expect an error if middleware not found/created
				assert.Nil(t, serverMW, "Expected no server middleware to be created")
			}
		})
	}
}

// TestMiddleware_EdgeCases tests edge cases of middleware
func TestMiddleware_EdgeCases(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	configFilePath := filepath.Join(filepath.Dir(filename), "configs", "config.yaml")

	rtInstance, _ := setupRuntime(t, configFilePath)

	t.Run("NilConfig", func(t *testing.T) {
		tempBootstrapPath := filepath.Join(t.TempDir(), "bootstrap.yaml")
		err := os.WriteFile(tempBootstrapPath, []byte("sources:\n  - kind: file\n    path: nonexistent.yaml"), 0644)
		require.NoError(t, err)

		tempAppInfo := rt.NewAppInfo("temp-nil-config-test", "1.0.0", rt.WithAppInfoID("temp-nil-config-test"))
		tempRtInstance, err := rt.NewFromBootstrap(tempBootstrapPath, rt.WithAppInfo(tempAppInfo))
		require.NoError(t, err)
		middlewareProvider, err := tempRtInstance.Container().Middleware()
		require.NoError(t, err)

		nilMWsMap, err := middlewareProvider.ClientMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, nilMWsMap, "ClientMiddlewares with nil config should return empty middlewares")

		nilMWsMap, err = middlewareProvider.ServerMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, nilMWsMap, "ServerMiddlewares with nil config should return empty middlewares")
	})

	t.Run("EmptyConfig", func(t *testing.T) {
		tempBootstrapPath := filepath.Join(t.TempDir(), "bootstrap.yaml")
		tempConfigPath := filepath.Join(t.TempDir(), "config.yaml")

		emptyConfigs := &middlewarev1.Middlewares{}
		configBytes, err := yaml.Marshal(emptyConfigs)
		require.NoError(t, err, "Failed to marshal empty middleware config to YAML")
		err = os.WriteFile(tempConfigPath, configBytes, 0644)
		require.NoError(t, err, "Failed to write empty config file")

		err = os.WriteFile(tempBootstrapPath, []byte("sources:\n  - kind: file\n    path: "+filepath.Base(tempConfigPath)), 0644)
		require.NoError(t, err)

		tempAppInfo := rt.NewAppInfo("temp-empty-config-test", "1.0.0", rt.WithAppInfoID("temp-empty-config-test"))
		tempRtInstance, err := rt.NewFromBootstrap(tempBootstrapPath, rt.WithAppInfo(tempAppInfo))
		require.NoError(t, err)
		tempMiddlewareProvider, err := tempRtInstance.Container().Middleware()
		require.NoError(t, err)

		// Test client middlewares with empty configuration
		clientMWsMap, err := tempMiddlewareProvider.ClientMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, clientMWsMap, "ClientMiddlewares with empty config should return empty middlewares")

		// Test server middlewares with empty configuration
		serverMWsMap, err := tempMiddlewareProvider.ServerMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, serverMWsMap, "ServerMiddlewares with empty config should return empty middlewares")
	})

	t.Run("DisabledMiddleware", func(t *testing.T) {
		tempBootstrapPath := filepath.Join(t.TempDir(), "bootstrap.yaml")
		tempConfigPath := filepath.Join(t.TempDir(), "config.yaml")

		disabledConfig := &middlewarev1.Middlewares{
			Configs: []*middlewarev1.Middleware{
				{
					Name:    "disabled-mw",
					Type:    "logging",
					Enabled: false,
					Logging: &middlewarev1.Logging{},
				},
			},
		}
		err := rtInstance.Config().Decode("", disabledConfig)
		require.NoError(t, err)

		err = os.WriteFile(tempBootstrapPath, []byte("sources:\n  - kind: file\n    path: "+filepath.Base(tempConfigPath)), 0644)
		require.NoError(t, err)

		tempAppInfo := rt.NewAppInfo("temp-disabled-mw-test", "1.0.0", rt.WithAppInfoID("temp-disabled-mw-test"))
		tempRtInstance, err := rt.NewFromBootstrap(tempBootstrapPath, rt.WithAppInfo(tempAppInfo))
		require.NoError(t, err)
		tempMiddlewareProvider, err := tempRtInstance.Container().Middleware()
		require.NoError(t, err)

		// Test that disabled middleware will not be created
		clientMWsMap, err := tempMiddlewareProvider.ClientMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, clientMWsMap, "Disabled middleware should not be created")

		serverMWsMap, err := tempMiddlewareProvider.ServerMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, serverMWsMap, "Disabled middleware should not be created")
	})

	t.Run("UnknownMiddlewareType", func(t *testing.T) {
		tempBootstrapPath := filepath.Join(t.TempDir(), "bootstrap.yaml")
		tempConfigPath := filepath.Join(t.TempDir(), "config.yaml")

		unknownConfig := &middlewarev1.Middlewares{
			Configs: []*middlewarev1.Middleware{
				{
					Name:    "unknown-mw",
					Type:    "nonexistent",
					Enabled: true,
				},
			},
		}
		err := rtInstance.Config().Decode("", unknownConfig)
		require.NoError(t, err)

		err = os.WriteFile(tempBootstrapPath, []byte("sources:\n  - kind: file\n    path: "+filepath.Base(tempConfigPath)), 0644)
		require.NoError(t, err)

		tempAppInfo := rt.NewAppInfo("temp-unknown-mw-test", "1.0.0", rt.WithAppInfoID("temp-unknown-mw-test"))
		tempRtInstance, err := rt.NewFromBootstrap(tempBootstrapPath, rt.WithAppInfo(tempAppInfo))
		require.NoError(t, err)
		tempMiddlewareProvider, err := tempRtInstance.Container().Middleware()
		require.NoError(t, err)

		// Test that middleware with unknown types will not be created
		clientMWsMap, err := tempMiddlewareProvider.ClientMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, clientMWsMap, "Unknown middleware type should not be created")

		serverMWsMap, err := tempMiddlewareProvider.ServerMiddlewares()
		require.NoError(t, err)
		assert.Empty(t, serverMWsMap, "Unknown middleware type should not be created")
	})
}
