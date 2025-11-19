package consul_source_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/test/helper"
	parentconfig "github.com/origadmin/runtime/test/integration/config"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"

	_ "github.com/origadmin/runtime/test/helper" // Import helper to ensure init() registers MockConsulSource
)

// ConsulSourceTestSuite tests the Consul configuration source integration.
type ConsulSourceTestSuite struct {
	suite.Suite
}

func TestConsulSourceTestSuite(t *testing.T) {
	suite.Run(t, new(ConsulSourceTestSuite))
}

// TestConsulSourceLoading verifies that configuration can be loaded correctly from a mock Consul source.
func (s *ConsulSourceTestSuite) TestConsulSourceLoading() {
	t := s.T()
	assert := assert.New(t)
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	bootstrapPath := "config/test_cases/consul_source/bootstrap_consul.yaml"

	// Get the absolute path to the mock JSON config file.
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file path")
	}
	currentDir := filepath.Dir(currentFile)
	mockConfigFilePath := filepath.Join(currentDir, "mock_config.json")

	// Read the mock JSON config file.
	data, err := os.ReadFile(mockConfigFilePath)
	assert.NoError(err, "Failed to read mock config file at %s: %v", mockConfigFilePath, err)

	// Unmarshal the JSON data into a map[string]interface{}.
	var mockData map[string]interface{}
	err = json.Unmarshal(data, &mockData)
	assert.NoError(err, "Failed to unmarshal mock config JSON: %v", err)

	// Initialize App. The framework should automatically use the registered MockConsulSource
	// based on the 'type: consul' in bootstrap_consul.yaml.
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:   "consul-test-app",
			Name: "ConsulApp", Version: "1.0.0",
		}),
		helper.WithMockDataJSON(mockData),
	)
	assert.NoError(err, "Failed to initialize runtime from bootstrap: %v", err)
	defer rtInstance.Cleanup()

	configDecoder := rtInstance.Config()
	assert.NotNil(configDecoder, "App ConfigDecoder should not be nil")

	var cfg testconfigs.TestConfig
	err = configDecoder.Decode("app", &cfg.App)
	assert.NoError(err, "Failed to decode app config from runtime: %v", err)

	err = configDecoder.Decode("logger", &cfg.Logger)
	assert.NoError(err, "Failed to decode logger config from runtime: %v", err)

	// Define the expected configurations based on the mock data.
	expectedApp := &appv1.App{
		Id:      "consul-app-id",
		Name:    "ConsulApp",
		Version: "1.0.0",
		Env:     "consul-test",
	}
	expectedLogger := &loggerv1.Logger{
		Level:  "warn",
		Format: "json",
	}

	// Perform assertions using the modular assertion toolkit.
	parentconfig.AssertAppConfig(t, expectedApp, cfg.App)
	parentconfig.AssertLoggerConfig(t, expectedLogger, cfg.Logger)

	t.Logf("Consul source config loaded and verified successfully!")
}
