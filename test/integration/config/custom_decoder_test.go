package config_test

import (
	"os"
	"path/filepath"
	stdruntime "runtime" // Alias the standard library runtime package
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
)

// TestCustomSettings represents the structure of our custom configuration section for testing.
type TestCustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
	RateLimit      int    `json:"rate_limit"`
	Endpoints      []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"endpoints"`
}

// TestRuntimeDecoder verifies that the configuration decoder is correctly exposed
// by the runtime and can be used to parse custom configuration sections.
func TestRuntimeDecoder(t *testing.T) {
	assert := assert.New(t)

	// Get the current file's directory to construct an absolute path for the config.
	_, filename, _, ok := stdruntime.Caller(0) // Use stdruntime.Caller
	if !ok {
		t.Fatalf("Failed to get current file info")
	}
	currentDir := filepath.Dir(filename)

	// Calculate the runtime module root directory.
	// From .../runtime/test/integration/config, go up 3 levels to .../runtime
	runtimeRoot := filepath.Join(currentDir, "../../..")

	// Store original CWD and defer its restoration.
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

	// Change CWD to the runtime module root.
	if err := os.Chdir(runtimeRoot); err != nil {
		t.Fatalf("Failed to change working directory to runtime root: %v", err)
	}

	// The bootstrapPath is now relative to the runtime module root.
	bootstrapPath := "examples/configs/load_with_custom_parser/config/bootstrap.yaml"

	// 1. Initialize Runtime with the correct AppInfo.
	rt, cleanup, err := runtime.NewFromBootstrap(
		bootstrapPath, // Use the path relative to the new CWD
		bootstrap.WithAppInfo("test-decoder", "1.0.0", "test"),
	)
	if err != nil {
		t.Fatalf("Failed to initialize runtime: %v", err)
	}
	defer cleanup()

	// 2. Verify core components are still initialized correctly.
	assert.NotNil(rt.Logger())
	assert.Equal("test-decoder", rt.AppInfo().Name())
	assert.Equal("1.0.0", rt.AppInfo().Version())
	assert.Equal("test", rt.AppInfo().Env())

	// 3. Get the ConfigDecoder from the runtime.
	decoder := rt.Decoder()
	assert.NotNil(decoder)

	// 4. Verify custom_settings are decoded correctly using the exposed decoder.
	var customSettings TestCustomSettings
	err = decoder.Decode("custom_settings", &customSettings)
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
