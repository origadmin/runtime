package custom_transformer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/bootstrap"
	parentconfig "github.com/origadmin/runtime/test/integration/config"
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

	// The path should be relative to the test's working directory.
	bootstrapPath := "bootstrap_transformer.yaml"

	// Create AppInfo using the new functional options pattern
	appInfo := rt.NewAppInfo(
		"TransformerTestApp",
		"1.0.0",
	).SetID("transformer-test-app")

	// Initialize App, which should apply the registered custom transformer.
	rtInstance := rt.New(
		appInfo.Name(),
		appInfo.Version(),
		rt.WithAppInfo(appInfo), // Pass the created AppInfo
	)

	// Removed defer rtInstance.Cleanup() as it's no longer available
	wd, _ := os.Getwd()
	fmt.Printf("working directory:%s\n", wd)
	// Load the configuration from the bootstrap file with all options.
	err := rtInstance.Load(bootstrapPath, bootstrap.WithConfigTransformer(&custom_transformer.TestTransformer{
		Suffix: "-transformed",
	}))
	require.NoError(t, err, "Failed to load configuration from bootstrap with custom transformer")

	defer rtInstance.Config().Close()

	// Define the expected configuration after transformation
	expectedApp := &appv1.App{
		Id:      "transformer-app-id",
		Name:    "OriginalApp-transformed",
		Version: "1.0.0",
		Env:     "transformer-test",
		// Metadata is not defined in the local config, so we expect the default nil or empty map.
		Metadata: nil, // Or map[string]string{} depending on desired behavior
	}
	// Use StructuredConfig() to get the transformed configuration
	ts, ok := rtInstance.StructuredConfig().DecodedConfig().(*testconfigs.TestConfig)
	if !ok {
		t.Fatalf("Failed to convert to TestTransformer")
	}
	require.NoError(t, err, "Failed to decode app config")
	// Perform a detailed, field-by-field assertion.
	parentconfig.AssertAppConfig(t, rt.ConvertToAppInfo(expectedApp), rt.ConvertToAppInfo(ts.App))

	t.Logf("Custom transformer applied and verified successfully!")
}
