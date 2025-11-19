package custom_transformer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	rt "github.com/origadmin/runtime"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	parentconfig "github.com/origadmin/runtime/test/integration/config"
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

	// The path should be relative to the test's working directory.
	bootstrapPath := "bootstrap_transformer.yaml"

	// Initialize App, which should apply the registered custom transformer.
	rtInstance, err := rt.NewFromBootstrap(
		bootstrapPath,
		bootstrap.WithAppInfo(&interfaces.AppInfo{
			ID:      "transformer-test-app",
			Name:    "TransformerTestApp",
			Version: "1.0.0",
		}),
		bootstrap.WithConfigTransformer(&custom_transformer.TestTransformer{Suffix: "-transformed"}),
	)
	// Use require.NoError to fail fast if the runtime fails to initialize.
	// This prevents panics from defer calls on a nil rtInstance.
	require.NoError(t, err, "Failed to initialize runtime from bootstrap with custom transformer")
	defer rtInstance.Cleanup()

	// Use StructuredConfig() to get the transformed configuration
	actualApp, err := rtInstance.StructuredConfig().DecodeApp()
	require.NoError(t, err, "Failed to decode app config from runtime")

	// Define the expected configuration after transformation
	expectedApp := &appv1.App{
		Id:      "transformer-app-id",
		Name:    "OriginalApp-transformed",
		Version: "1.0.0",
		Env:     "transformer-test",
		// Metadata is not defined in the local config, so we expect the default nil or empty map.
		Metadata: nil, // Or map[string]string{} depending on desired behavior
	}

	// Perform a detailed, field-by-field assertion.
	parentconfig.AssertAppConfig(t, expectedApp, actualApp)

	t.Logf("Custom transformer applied and verified successfully!")
}
