package multiple_sources_config_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
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
	assertions := assert.New(t)

	// Use a robust relative path from the test file to the dedicated test data.
	bootstrapPath := filepath.Join("testdata", "bootstrap_multiple_sources.yaml")

	// Initialize Runtime from the bootstrap file that defines multiple sources.
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "multi-source-test-app",
			Name:    "MultiSourceTestApp",
			Version: "1.0.0",
		}),
	)
	assertions.NoError(err, "Failed to initialize runtime from bootstrap: %v", err)
	defer rtInstance.Cleanup()

	configDecoder := rtInstance.Config()
	assertions.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil")

	var cfg testconfigs.TestConfig
	err = configDecoder.Decode("", &cfg)
	assertions.NoError(err, "Failed to decode config from runtime: %v", err)

	// Assertions for merged configuration:
	// app.id should be from source1 (not overridden)
	assertions.NotNil(cfg.App)
	assertions.Equal("source1-app-id", cfg.App.Id)

	// app.name should be overridden by source2
	assertions.Equal("Source2App", cfg.App.Name)

	// app.env should be overridden by source2
	assertions.Equal("prod", cfg.App.Env)

	// app.metadata should contain key3 from source2
	assertions.Contains(cfg.App.Metadata, "key3")
	assertions.Equal("value3", cfg.App.Metadata["key3"])

	// logger.level should be overridden by source2
	assertions.NotNil(cfg.Logger)
	assertions.Equal("info", cfg.Logger.Level)

	// client.timeout should be overridden by source2
	assertions.NotNil(cfg.Client)
	assertions.Equal("5s", cfg.Client.Timeout.AsDuration().String())

	// client.endpoint should be from source1 (not overridden)
	assertions.Equal("discovery:///source1-service", cfg.Client.Endpoint)

	// client.selector.version should be from source2 (new field)
	assertions.NotNil(cfg.Client.Selector)
	assertions.Equal("v2.0.0", cfg.Client.Selector.Version)

	t.Log("Multiple sources config loaded and merged successfully!")
}
