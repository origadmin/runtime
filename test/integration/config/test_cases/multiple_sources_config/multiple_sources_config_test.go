package multiple_sources_config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/test/helper"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// MultipleSourcesConfigTestSuite tests the coordination of multiple configuration sources.
type MultipleSourcesConfigTestSuite struct {
	suite.Suite
}

func TestMultipleSourcesConfigTestSuite(t *testing.T) {
	suite.Run(t, new(MultipleSourcesConfigTestSuite))
}

// TestMultipleSourcesLoading verifies that configurations from multiple sources are correctly
// loaded, merged, and prioritized by the runtime.
func (s *MultipleSourcesConfigTestSuite) TestMultipleSourcesLoading() {
	t := s.T()
	assert := assert.New(t)
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	bootstrapPath := "test/integration/config/test_cases/multiple_sources_config/bootstrap_multiple_sources.yaml"

	// Initialize Runtime from the bootstrap file that defines multiple sources.
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "multi-source-test-app",
			Name:    "MultiSourceTestApp",
			Version: "1.0.0",
		}),
	)
	assert.NoError(err, "Failed to initialize runtime from bootstrap: %v", err)
	defer rtInstance.Cleanup()

	configDecoder := rtInstance.Config()
	assert.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil")

	var cfg testconfigs.TestConfig
	err = configDecoder.Decode("", &cfg)
	assert.NoError(err, "Failed to decode config from runtime: %v", err)

	// Assertions for merged configuration:
	// app.id should be from source1 (not overridden)
	assert.NotNil(cfg.App)
	assert.Equal("source1-app-id", cfg.App.Id)

	// app.name should be overridden by source2
	assert.Equal("Source2App", cfg.App.Name)

	// app.env should be overridden by source2
	assert.Equal("prod", cfg.App.Env)

	// app.metadata should contain key3 from source2
	assert.Contains(cfg.App.Metadata, "key3")
	assert.Equal("value3", cfg.App.Metadata["key3"])

	// logger.level should be overridden by source2
	assert.NotNil(cfg.Logger)
	assert.Equal("info", cfg.Logger.Level)

	// client.timeout should be overridden by source2
	assert.NotNil(cfg.Client)
	assert.Equal("5s", cfg.Client.Timeout.AsDuration().String())

	// client.endpoint should be from source1 (not overridden)
	assert.Equal("discovery:///source1-service", cfg.Client.Endpoint)

	// client.selector.version should be from source2 (new field)
	assert.NotNil(cfg.Client.Selector)
	assert.Equal("v2.0.0", cfg.Client.Selector.Version)

	t.Log("Multiple sources config loaded and merged successfully!")
}
