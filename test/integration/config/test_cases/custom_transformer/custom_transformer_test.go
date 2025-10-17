package custom_transformer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/test/helper"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
	_ "github.com/origadmin/runtime/test/integration/config/test_cases/custom_transformer" // Import for transformer registration
)

// CustomTransformerTestSuite tests the integration of custom configuration transformers.
type CustomTransformerTestSuite struct {
	suite.Suite
}

func TestCustomTransformerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomTransformerTestSuite))
}

// TestCustomTransformerApplication verifies that a custom transformer is correctly applied
// during the bootstrap process and modifies the configuration as expected.
func (s *CustomTransformerTestSuite) TestCustomTransformerApplication() {
	t := s.T()
	assert := assert.New(t)
	cleanup := helper.SetupIntegrationTest(t)
	defer cleanup()

	bootstrapPath := "test/integration/config/test_cases/custom_transformer/bootstrap_transformer.yaml"

	// Initialize Runtime, which should apply the registered custom transformer.
	rtInstance, rtCleanup, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "transformer-test-app",
			Name:    "TransformerTestApp",
			Version: "1.0.0",
		}),
	)
	assert.NoError(err, "Failed to initialize runtime from bootstrap with custom transformer: %v", err)
	defer rtCleanup()

	configDecoder := rtInstance.Config()
	assert.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil")

	var cfg testconfigs.TestConfig
	err = configDecoder.Decode("app", &cfg.App) // Decode only the app section for specific assertions
	assert.NoError(err, "Failed to decode app config from runtime: %v", err)

	// Assert that the App.Name has been transformed by our custom transformer.
	assert.NotNil(cfg.App)
	assert.Equal("transformer-app-id", cfg.App.Id, "App ID should remain unchanged")
	assert.Equal("OriginalApp-transformed", cfg.App.Name, "App Name should be transformed by the custom transformer")
	assert.Equal("1.0.0", cfg.App.Version, "App Version should remain unchanged")
	assert.Equal("transformer-test", cfg.App.Env, "App Env should remain unchanged")

	t.Logf("Custom transformer applied and verified successfully!")
}
