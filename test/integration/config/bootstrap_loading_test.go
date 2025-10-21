// Tests integration with Runtime
package config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// RuntimeIntegrationTestSuite tests configuration integration with Runtime
type RuntimeIntegrationTestSuite struct {
	suite.Suite
}

func TestRuntimeIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RuntimeIntegrationTestSuite))
}

// TestRuntimeLoadCompleteConfig tests loading complete configuration using Runtime
func (s *RuntimeIntegrationTestSuite) TestRuntimeLoadCompleteConfig() {
	t := s.T()
	assert := assert.New(t)

	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file info")
	}
	currentDir := filepath.Dir(filename)

	// Save original working directory and restore it after test
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	defer func() {
		err := os.Chdir(originalCwd)
		if err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	// Change the working directory to the root directory of the runtime module
	runtimeRoot := filepath.Join(currentDir, "../../..")
	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	// Run tests for each format
	formats := []string{"yaml", "json", "toml"}
	for _, format := range formats {
		t.Run("Runtime_"+format, func(t *testing.T) {
			// Create temporary bootstrap config file
			tempDir := t.TempDir()
			tempBootstrapPath := filepath.Join(tempDir, "bootstrap."+format)

			// Write temporary bootstrap config
			var bootstrapContent string
			switch format {
			case "yaml":
				bootstrapContent = "sources:\n  - type: \"file\"\n    name: \"complete-config\"\n    file:\n      path: \"test/integration/config/configs/complete_config." + format + "\"\n    priority: 100"
			case "json":
				bootstrapContent = `{\"sources\": [{\"type\": \"file\", \"name\": \"complete-config\", \"file\": {\"path\": \"test/integration/config/configs/complete_config.` + format + `\"}, \"priority\": 100}]}`
			case "toml":
				bootstrapContent = `[[sources]]\ntype = \"file\"\nname = \"complete-config\"\nfile.path = \"test/integration/config/configs/complete_config.` + format + `\"\npriority = 100`
			}

			if err := os.WriteFile(tempBootstrapPath, []byte(bootstrapContent), 0644); err != nil {
				t.Fatalf("Failed to write temp bootstrap file: %v", err)
			}

			// Initialize Runtime
			rt, cleanup, err := rt.NewFromBootstrap(
				tempBootstrapPath,
				bootstrap.WithAppInfo(&interfaces.AppInfo{
					ID:      "test-complete-config",
					Name:    "TestCompleteConfig",
					Version: "1.0.0",
				}),
			)
			if err != nil {
				t.Fatalf("Failed to initialize runtime with %s config: %v", format, err)
			}
			defer cleanup()

			// Get configuration decoder
			configDecoder := rt.Config()
			assert.NotNil(configDecoder)

			// Decode into complete config structure
			var config testconfigs.TestConfig
			err = configDecoder.Decode("", &config)
			assert.NoError(err)

			// Run assertions
			AssertTestConfig(t, &config)
			t.Logf("Runtime loaded and verified %s format complete config successfully!", format)
		})
	}
}

// TestConfigProtoIntegration tests integration between configuration and Protocol Buffers
func (s *RuntimeIntegrationTestSuite) TestConfigProtoIntegration() {
	t := s.T()
	assert := assert.New(t)

	// Get the directory of the current file to build the absolute path of the configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file info")
	}
	currentDir := filepath.Dir(filename)

	// Calculate the root of the runtime module
	// From .../runtime/test/integration/config, go up 3 levels to .../runtime
	runtimeRoot := filepath.Join(currentDir, "../../..")

	// Save original working directory and restore it after test
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	defer func() {
		err := os.Chdir(originalCwd)
		if err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	// Change the working directory to the root directory of the runtime module
	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	// bootstrapPath is now relative to the runtime module root
	bootstrapPath := "test/integration/config/configs/bootstrap.yaml"

	// 1. Initialize Runtime with default decoder provider
	boot, cleanup, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-proto-config",
			Name:    "TestProtoConfig",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer cleanup()

	// 2. Get ConfigDecoder from runtime
	configDecoder := boot.Config()
	assert.NotNil(configDecoder)

	// 3. Decode entire config into generated Bootstrap struct
	var bootstrapConfig testconfigs.TestConfig
	err = configDecoder.Decode("", &bootstrapConfig)
	assert.NoError(err)

	// 4. Assert decoded values
	// Verify logger (from test_config.yaml)
	logger := boot.Logger()
	assert.NotNil(logger)

	// Verify registration_discovery_name
	assert.Equal("test-discovery", bootstrapConfig.RegistrationDiscoveryName)

	// Verify servers
	assert.Len(bootstrapConfig.GetServers(), 1, "Should have 1 Servers message")
	serverConfigs := bootstrapConfig.GetServers()[0].GetServers()
	assert.Len(serverConfigs, 2, "Should have 2 Server configurations (gRPC and HTTP)")

	var grpcServer *transportv1.Server
	var httpServer *transportv1.Server

	for _, s := range serverConfigs {
		if s.GetGrpc() != nil {
			grpcServer = s
		}
		if s.GetHttp() != nil {
			httpServer = s
		}
	}

	assert.NotNil(grpcServer, "gRPC server configuration not found")
	assert.Equal("tcp", grpcServer.GetGrpc().GetNetwork())
	assert.Equal(":9000", grpcServer.GetGrpc().GetAddr())
	assert.Equal("1s", grpcServer.GetGrpc().GetTimeout().AsDuration().String())

	assert.NotNil(httpServer, "HTTP server configuration not found")
	assert.Equal("tcp", httpServer.GetHttp().GetNetwork())
	assert.Equal(":8000", httpServer.GetHttp().GetAddr())
	assert.Equal("2s", httpServer.GetHttp().GetTimeout().AsDuration().String())

	// Verify clients (should be empty as not defined in test_config.yaml)
	assert.Empty(bootstrapConfig.Clients)

	// Verify discoveries (should be empty as not defined in test_config.yaml)
	assert.Empty(bootstrapConfig.Discoveries)
}

// TestCustomSettings represents the structure for custom configuration section in tests
type TestCustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"endpoints"`
}

// TestRuntimeDecoder verifies that the configuration decoder is properly exposed by runtime
// and can be used to parse custom configuration sections
func (s *RuntimeIntegrationTestSuite) TestRuntimeDecoder() {
	t := s.T()
	assert := assert.New(t)

	// Get the directory of the current file to build the absolute path of the configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file info")
	}
	currentDir := filepath.Dir(filename)

	// Calculate the root of the runtime module
	// From .../runtime/test/integration/config, go up 3 levels to .../runtime
	runtimeRoot := filepath.Join(currentDir, "../../..")

	// Save original working directory and restore it after test
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get original working directory: %v", err)
	}
	defer func() {
		err := os.Chdir(originalCwd)
		if err != nil {
			t.Errorf("Failed to restore original working directory: %v", err)
		}
	}()

	// Change the working directory to the root directory of the runtime module
	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	bootstrapPath := "examples/configs/load_with_custom_parser/config/bootstrap.yaml"

	// 1. Initialize Runtime with correct AppInfo
	rt, cleanup, err := rt.NewFromBootstrap(
		bootstrapPath, // Use path relative to the new working directory
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "test-decoder",
			Name:    "TestDecoder",
			Version: "1.0.0",
			Env:     "dev", // Changed from "test" to "dev" to match the config
		}),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer cleanup()

	// 2. Verify core components are properly initialized
	assert.NotNil(rt.Logger())
	assert.Equal("test-decoder", rt.AppInfo().ID)
	assert.Equal("TestDecoder", rt.AppInfo().Name)
	assert.Equal("1.0.0", rt.AppInfo().Version)
	assert.Equal("dev", rt.AppInfo().Env) // Changed expected value

	// 3. Get ConfigDecoder from runtime
	decoder := rt.Config()
	assert.NotNil(decoder)

	// 4. Verify custom_settings is properly decoded using the exposed decoder
	var customSettings TestCustomSettings
	// Updated to use "components.my-custom-settings" path to match the config structure
	err = decoder.Decode("components.my-custom-settings", &customSettings)
	assert.NoError(err)
	assert.True(customSettings.FeatureEnabled)
	assert.Equal("super-secret-key-123", customSettings.APIKey)
	assert.Equal(100, customSettings.RateLimit)
	assert.Len(customSettings.Endpoints, 2)
	assert.Equal("users", customSettings.Endpoints[0].Name)
	assert.Equal("/api/v1/users", customSettings.Endpoints[0].Path)
	assert.Equal("products", customSettings.Endpoints[1].Name)
	assert.Equal("/api/v1/products", customSettings.Endpoints[1].Path)
}
