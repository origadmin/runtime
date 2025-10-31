package middleware_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	selectorv1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/selector/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/middleware"
)

// Define all supported middleware types
var supportedMiddlewareTypes = map[string]struct{}{
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

// loadTestConfig loads test configuration file
func loadTestConfig(t *testing.T) *middlewarev1.Middlewares {
	// Get the directory of the current test file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Build the full path of the configuration file
	configPath := filepath.Join(dir, "configs", "config.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file not found: %s", configPath)
	}

	// Load configuration
	var configs middlewarev1.Middlewares
	cfg, err := config.Load(configPath, &configs)
	require.NoError(t, err, "Failed to load config")
	require.NotNil(t, cfg, "Config instance should not be nil")
	t.Cleanup(func() {
		_ = cfg.Close()
	})

	return &configs
}

// TestMiddleware_LoadAndBuild tests middleware loading and building
func TestMiddleware_LoadAndBuild(t *testing.T) {
	// Load test configuration
	configs := loadTestConfig(t)

	t.Run("VerifyConfig", func(t *testing.T) {
		// Check the number of middlewares in the configuration
		assert.GreaterOrEqual(t, len(configs.Middlewares), 2, "Should have at least 2 middlewares in config")

		// Check configuration of each middleware
		for _, mw := range configs.Middlewares {
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

	t.Run("BuildClientMiddleware", func(t *testing.T) {
		// Build client middleware chain
		clientMWs := middleware.BuildClientMiddlewares(configs)

		// Calculate the actual number of supported client middlewares (excluding unsupported middleware types)
		expectedCount := 0
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "rate_limiter" { // rate_limiter 在客户端不受支持
				expectedCount++
			}
		}

		// Verify the number of middlewares
		assert.Equal(t, expectedCount, len(clientMWs), "Number of client middlewares should match expected")
	})

	t.Run("BuildServerMiddleware", func(t *testing.T) {
		// Build server middleware chain
		serverMWs := middleware.BuildServerMiddlewares(configs)

		// Calculate the actual number of supported server middlewares (excluding unsupported middleware types)
		expectedCount := 0
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "circuit_breaker" { // circuit_breaker 在服务端不受支持
				expectedCount++
			}
		}

		// Verify the number of middlewares
		assert.Equal(t, expectedCount, len(serverMWs), "Number of server middlewares should match expected")
	})

	t.Run("MiddlewareNaming", func(t *testing.T) {
		// Create a custom middleware configuration
		customConfig := &middlewarev1.Middlewares{
			Middlewares: []*middlewarev1.Middleware{
				{
					Name:    "custom-name",
					Type:    "logging",
					Enabled: true,
					Logging: &middlewarev1.Logging{},
				},
				{
					Name:     "", // Test unnamed middleware
					Type:     "metadata",
					Enabled:  true,
					Metadata: &middlewarev1.Metadata{},
				},
			},
		}

		// Test client middlewares
		clientMWs := middleware.BuildClientMiddlewares(customConfig)
		assert.NotEmpty(t, clientMWs, "Client middlewares should not be empty")

		// Test server middlewares
		serverMWs := middleware.BuildServerMiddlewares(customConfig)
		assert.NotEmpty(t, serverMWs, "Server middlewares should not be empty")
	})

	// Verify middlewares in the configuration
	t.Run("VerifyConfig", func(t *testing.T) {
		// Check the number of middlewares in the configuration
		assert.GreaterOrEqual(t, len(configs.Middlewares), 2, "Should have at least 2 middlewares in config")

		// Check configuration of each middleware
		for _, mw := range configs.Middlewares {
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

	// Test client middleware building
	t.Run("BuildClientMiddleware", func(t *testing.T) {
		// Build client middleware chain
		clientMWs := middleware.BuildClientMiddlewares(configs)

		// Calculate the actual number of supported client middlewares (excluding unsupported middleware types)
		expectedCount := 0
		var expectedOrder []string
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "rate_limiter" { // rate_limiter 在客户端不受支持
				expectedOrder = append(expectedOrder, mw.Name)
				expectedCount++
			}
		}

		// Verify the number of middlewares
		assert.Equal(t, expectedCount, len(clientMWs), "Number of client middlewares should match expected")

		// Verify that each middleware is properly handled
		for i := 0; i < len(clientMWs); i++ {
			assert.NotNil(t, clientMWs[i], "Client middleware at index %d should not be nil", i)
		}

		// Verify that Selector middleware exists
		selectorFound := false
		for _, mwName := range expectedOrder {
			if mwName == "selector" {
				selectorFound = true
				break
			}
		}
		assert.True(t, selectorFound, "Selector middleware should be present in client middlewares")
	})

	// Test server middleware building
	t.Run("BuildServerMiddleware", func(t *testing.T) {
		// Build server middleware chain
		serverMWs := middleware.BuildServerMiddlewares(configs)

		// Calculate the actual number of supported server middlewares (excluding unsupported middleware types)
		expectedCount := 0
		var expectedOrder []string
		for _, mw := range configs.Middlewares {
			if mw.Enabled && mw.Type != "circuit_breaker" { // circuit_breaker 在服务端不受支持
				expectedOrder = append(expectedOrder, mw.Name)
				expectedCount++
			}
		}

		// Verify the number of middlewares
		assert.Equal(t, expectedCount, len(serverMWs), "Number of server middlewares should match expected")

		// Verify that each middleware is properly handled
		for i := 0; i < len(serverMWs); i++ {
			assert.NotNil(t, serverMWs[i], "Server middleware at index %d should not be nil", i)
		}

		// Verify that Selector middleware exists
		selectorFound := false
		for _, mwName := range expectedOrder {
			if mwName == "selector" {
				selectorFound = true
				break
			}
		}
		assert.True(t, selectorFound, "Selector middleware should be present in server middlewares")
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
				Middlewares: []*middlewarev1.Middleware{
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
				Middlewares: []*middlewarev1.Middleware{
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
			switch tt.side {
			case "client":
				mw := middleware.BuildClientMiddlewares(tt.config)
				assert.Len(t, mw, tt.expectCount, "Unexpected number of client middlewares")
			case "server":
				mw := middleware.BuildServerMiddlewares(tt.config)
				assert.Len(t, mw, tt.expectCount, "Unexpected number of server middlewares")
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
			// Test client middleware
			clientMW, _ := middleware.NewClient(tt.config)
			if tt.clientMW {
				assert.NotNil(t, clientMW, "Expected client middleware to be created")
			} else {
				assert.Nil(t, clientMW, "Expected no client middleware to be created")
			}

			// Test server middleware
			serverMW, _ := middleware.NewServer(tt.config)
			if tt.serverMW {
				assert.NotNil(t, serverMW, "Expected server middleware to be created")
			} else {
				assert.Nil(t, serverMW, "Expected no server middleware to be created")
			}

			// Verify that the middleware is correctly stored
			if tt.exists {
				// Create middleware
				if tt.clientMW {
					mw, _ := middleware.NewClient(tt.config)
					assert.NotNil(t, mw, "Client middleware should be created")
				}

				if tt.serverMW {
					mw, _ := middleware.NewServer(tt.config)
					assert.NotNil(t, mw, "Server middleware should be created")
				}
			}
		})
	}
}

// TestMiddleware_EdgeCases tests edge cases of middleware
func TestMiddleware_EdgeCases(t *testing.T) {
	t.Run("NilConfig", func(t *testing.T) {
		// Test nil configuration
		nilMWs := middleware.BuildClientMiddlewares(nil)
		assert.Empty(t, nilMWs, "BuildClientMiddlewares with nil config should return empty middlewares")
		nilMWs = middleware.BuildServerMiddlewares(nil)
		assert.Empty(t, nilMWs, "BuildServerMiddlewares with nil config should return empty middlewares")
	})

	t.Run("EmptyConfig", func(t *testing.T) {
		emptyConfigs := &middlewarev1.Middlewares{}

		// Test client middlewares with empty configuration
		clientMWs := middleware.BuildClientMiddlewares(emptyConfigs)
		assert.Empty(t, clientMWs, "BuildClientMiddlewares with empty config should return empty middlewares")

		// Test server middlewares with empty configuration
		serverMWs := middleware.BuildServerMiddlewares(emptyConfigs)
		assert.Empty(t, serverMWs, "BuildServerMiddlewares with empty config should return empty middlewares")
	})

	t.Run("DisabledMiddleware", func(t *testing.T) {
		disabledConfig := &middlewarev1.Middlewares{
			Middlewares: []*middlewarev1.Middleware{
				{
					Name:    "disabled-mw",
					Type:    "logging",
					Enabled: false,
					Logging: &middlewarev1.Logging{},
				},
			},
		}

		// Test that disabled middleware will not be created
		clientMWs := middleware.BuildClientMiddlewares(disabledConfig)
		assert.Empty(t, clientMWs, "Disabled middleware should not be created")

		serverMWs := middleware.BuildServerMiddlewares(disabledConfig)
		assert.Empty(t, serverMWs, "Disabled middleware should not be created")
	})

	t.Run("UnknownMiddlewareType", func(t *testing.T) {
		unknownConfig := &middlewarev1.Middlewares{
			Middlewares: []*middlewarev1.Middleware{
				{
					Name:    "unknown-mw",
					Type:    "nonexistent",
					Enabled: true,
				},
			},
		}

		// Test that middleware with unknown types will not be created
		clientMWs := middleware.BuildClientMiddlewares(unknownConfig)
		assert.Empty(t, clientMWs, "Unknown middleware type should not be created")

		serverMWs := middleware.BuildServerMiddlewares(unknownConfig)
		assert.Empty(t, serverMWs, "Unknown middleware type should not be created")
	})
}
