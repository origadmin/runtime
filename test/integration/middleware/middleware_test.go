package middleware

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	rt "github.com/origadmin/runtime"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/config/middleware/selector/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/bootstrap"
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
	if mw == nil {
		return nil
	}
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

// setupRuntimeFromFile initializes a runtime instance from a given configuration file path.
// This function is dedicated to file-based test scenarios.
func setupRuntimeFromFile(t *testing.T, appID, configFilePath string) *rt.App {
	t.Helper()
	require.NotEmpty(t, configFilePath, "configFilePath cannot be empty for file-based setup")

	rtInstance := rt.NewWithOptions(
		rt.WithAppInfo(
			rt.NewAppInfo(appID, "1.0.0").SetEnv("development"),
		),
		rt.WithID(appID),
	)

	err := rtInstance.Load(configFilePath, bootstrap.WithDirectly())
	require.NoError(t, err, "Failed to load configuration from file: %s", configFilePath)
	defer rtInstance.Config().Close()

	return rtInstance
}

// setupRuntimeFromObject initializes a runtime instance from an in-memory middleware configuration object.
// It creates a temporary configuration file for the test.
func setupRuntimeFromObject(t *testing.T, appID string, mws *middlewarev1.Middlewares) *rt.App {
	t.Helper()

	tempDir := t.TempDir()
	finalConfigPath := filepath.Join(tempDir, "config.yaml")

	if mws == nil {
		mws = &middlewarev1.Middlewares{}
	}

	// The config structure expects a top-level 'middlewares' key.
	configWrapper := map[string]interface{}{
		"middlewares": mws,
	}
	yamlBytes, err := yaml.Marshal(configWrapper)
	require.NoError(t, err, "Failed to marshal in-memory config to YAML")
	err = os.WriteFile(finalConfigPath, yamlBytes, 0644)
	require.NoError(t, err, "Failed to write temp config file")

	// Since we created the file, we now delegate to the file-based setup function.
	// This avoids duplicating the NewFromBootstrap call.
	return setupRuntimeFromFile(t, appID, finalConfigPath)
}

// TestMiddleware_LoadAndBuild tests middleware loading and building
func TestMiddleware_LoadAndBuild(t *testing.T) {
	configFilePath := "configs/config.yaml"

	// Let the runtime handle the loading from the specified config file.
	rtInstance := setupRuntimeFromFile(t, "load-and-build-test", configFilePath)

	// Get the configuration directly from the initialized runtime.
	// This is the correct way, as it uses the framework's own decoding logic.
	configs, err := rtInstance.StructuredConfig().DecodeMiddlewares()
	var source map[string]any
	_ = rtInstance.Config().Decode("", &source)
	t.Logf("Loaded configs: %+v", source)
	require.NoError(t, err, "Failed to decode middlewares from structured config")
	t.Logf("Loaded middlewares: %+v", configs)
	t.Run("VerifyConfig", func(t *testing.T) {
		assert.GreaterOrEqual(t, len(configs.Configs), 2, "Should have at least 2 middlewares in config")
		t.Logf("Loaded middlewares: %+v", configs.Configs)
		// Check configuration of each middleware
		for _, mw := range configs.Configs {
			if mw == nil {
				t.Error("Found nil middleware configuration in configs.Configs")
				continue
			}
			t.Logf("Checking middleware: %s (type: %s, enabled: %v)", mw.Name, mw.Type, mw.Enabled)
			assert.NotEmpty(t, mw.Name, "Middleware name should not be empty")
			assert.NotEmpty(t, mw.Type, "Middleware type should not be empty")
			assert.True(t, mw.Enabled, "Middleware should be enabled by default")

			// Verify the configuration based on the middleware type
			field := getMiddlewareField(mw, mw.Type)
			assert.NotNil(t, field, "Middleware configuration for type %s should not be nil", mw.Type)
		}
	})

	t.Run("ClientMiddlewares", func(t *testing.T) {
		// Get middleware provider
		middlewareProvider, err := rtInstance.Container().Middleware()
		require.NoError(t, err)

		// Log the middleware provider's configuration
		middlewareConfig, err := rtInstance.StructuredConfig().DecodeMiddlewares()
		require.NoError(t, err)
		t.Logf("Middleware config from container: %+v", middlewareConfig)

		// Get client middlewares
		clientMWsMap, err := middlewareProvider.ClientMiddlewares()
		if err != nil {
			t.Logf("Error getting client middlewares: %v", err)
		}
		require.NoError(t, err)

		t.Logf("Client middlewares created: %d", len(clientMWsMap))
		for name := range clientMWsMap {
			t.Logf("  - %s", name)
		}

		// Calculate the expected number of client middlewares
		expectedCount := 0
		// From config.yaml: metadata, logging, selector, rate_limiter, circuit_breaker
		// All are enabled.
		// 'selector' middleware itself is a middleware.
		// 'rate_limiter' is not supported on the client side.
		// So, we expect: metadata, logging, selector, circuit_breaker
		expectedCount = 4

		t.Logf("Loaded middlewares from runtime: %+v", configs)
		// Log detailed information about the middlewares
		t.Logf("Test Configuration:")
		t.Logf("  - Expected client middlewares: %d", expectedCount)
		t.Logf("  - Found client middlewares: %d", len(clientMWsMap))
		t.Logf("  - Configured middlewares:")
		for i, mw := range configs.Configs {
			if mw == nil {
				t.Logf("    %d. <nil>", i+1)
				continue
			}
			t.Logf("    %d. %s (type: %s, enabled: %v)", i+1, mw.Name, mw.Type, mw.Enabled)
		}

		if len(clientMWsMap) == 0 {
			t.Log("No client middlewares were found. This might be due to:")
			t.Log("1. Middleware provider not properly initialized")
			t.Log("2. No middleware factories registered")
			t.Log("3. Configuration not properly loaded")
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
		// From config.yaml: metadata, logging, selector, rate_limiter, circuit_breaker
		// All are enabled.
		// 'circuit_breaker' is not supported on the server side.
		// So, we expect: metadata, logging, selector, rate_limiter
		expectedCount = 4

		// Verify the number of middlewares
		assert.Equal(t, expectedCount, len(serverMWsMap), "Number of server middlewares should match expected")
	})
}

// TestSelectorMiddleware tests the includes/excludes functionality of Selector middleware
func TestSelectorMiddleware(t *testing.T) {
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
			// Use the centralized setup function with the test-specific config.
			rtInstance := setupRuntimeFromObject(t, "selector-test-app", tt.config)

			switch tt.side {
			case "client":
				middlewareProvider, err := rtInstance.Container().Middleware()
				require.NoError(t, err)
				mwMap, err := middlewareProvider.ClientMiddlewares()
				require.NoError(t, err)
				assert.Len(t, mwMap, tt.expectCount, "Unexpected number of client middlewares")
			case "server":
				middlewareProvider, err := rtInstance.Container().Middleware()
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
			// Create a minimal middleware config for this specific test case.
			tempConfigs := &middlewarev1.Middlewares{Configs: []*middlewarev1.Middleware{tt.config}}
			// Use the centralized setup function.
			tempRtInstance := setupRuntimeFromObject(t, "creation-test-app", tempConfigs)
			middlewareProvider, err := tempRtInstance.Container().Middleware()
			require.NoError(t, err)

			// Test client middleware
			clientMW, err := middlewareProvider.ClientMiddleware(tt.config.Name)
			if tt.clientMW {
				assert.NoError(t, err)
				assert.NotNil(t, clientMW, "Expected client middleware to be created")
			} else {
				assert.Error(t, err, "Expected an error for non-existent client middleware")
				assert.Nil(t, clientMW, "Expected no client middleware to be returned")
			}

			// Test server middleware
			serverMW, err := middlewareProvider.ServerMiddleware(tt.config.Name)
			if tt.serverMW {
				assert.NoError(t, err)
				assert.NotNil(t, serverMW, "Expected server middleware to be created")
			} else {
				assert.Error(t, err, "Expected an error for non-existent server middleware")
				assert.Nil(t, serverMW, "Expected no server middleware to be returned")
			}
		})
	}
}

// TestMiddleware_EdgeCases tests edge cases of middleware
func TestMiddleware_EdgeCases(t *testing.T) {
	t.Run("NilConfig", func(t *testing.T) {
		// setupRuntime handles nil config gracefully.
		tempRtInstance := setupRuntimeFromObject(t, "nil-config-test", nil)
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
		emptyConfigs := &middlewarev1.Middlewares{}
		tempRtInstance := setupRuntimeFromObject(t, "empty-config-test", emptyConfigs)
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
		tempRtInstance := setupRuntimeFromObject(t, "disabled-mw-test", disabledConfig)
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
		unknownConfig := &middlewarev1.Middlewares{
			Configs: []*middlewarev1.Middleware{
				{
					Name:    "unknown-mw",
					Type:    "nonexistent",
					Enabled: true,
				},
			},
		}
		tempRtInstance := setupRuntimeFromObject(t, "unknown-mw-test", unknownConfig)
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
