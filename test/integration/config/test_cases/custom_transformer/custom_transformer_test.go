package custom_transformer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
	"github.com/origadmin/runtime/test/integration/config/test_cases/custom_transformer"
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
	assertions := assert.New(t)

	bootstrapPath := "bootstrap_transformer.yaml"

	// Initialize Runtime, which should apply the registered custom transformer.
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "transformer-test-app",
			Name:    "TransformerTestApp",
			Version: "1.0.0",
		}),
		bootstrap.WithConfigTransformer(&custom_transformer.TestTransformer{Suffix: "-transformed"}),
	)
	assertions.NoError(err, "Failed to initialize runtime from bootstrap with custom transformer: %v", err)
	defer rtInstance.Cleanup()

	// Use StructuredConfig() to get the transformed configuration, not Config().
	// Config() returns the raw, untransformed configuration.
	configDecoder := rtInstance.Config()
	assertions.NotNil(configDecoder, "Runtime ConfigDecoder should not be nil")

	var cfg testconfigs.TestConfig
	err = configDecoder.Decode("", &cfg) // Decode only the app section for specific assertions
	assertions.NoError(err, "Failed to decode app config from runtime: %v", err)
	// Assert that the App.Name has been transformed by our custom transformer.
	app, err := rtInstance.StructuredConfig().DecodeApp()
	assertions.NoError(err, "Failed to decode app config from runtime: %v", err)
	assertions.NotNil(app, "App section should be decoded")
	assertions.Equal("transformer-app-id", app.Id, "App ID should remain unchanged")
	assertions.Equal("OriginalApp-transformed", app.Name, "App Name should be transformed by the custom transformer")
	assertions.Equal("1.0.0", app.Version, "App Version should remain unchanged")
	assertions.Equal("transformer-test", app.Env, "App Env should remain unchanged")

	t.Logf("Custom transformer applied and verified successfully!")
}
