package multiple_sources_config_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/durationpb"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	selectorv1 "github.com/origadmin/runtime/api/gen/go/config/selector/v1"
	grpcv1 "github.com/origadmin/runtime/api/gen/go/config/transport/grpc/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	parentconfig "github.com/origadmin/runtime/test/integration/config"
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

	bootstrapPath := filepath.Join("testdata", "bootstrap_multiple_sources.yaml")

	// Create AppInfo using the new functional options pattern
	appInfo := rt.NewAppInfo(
		"MultiSourceTestApp",
		"1.0.0",
		rt.WithAppInfoID("multi-source-test-app"),
	)

	rtInstance, err := rt.New(
		appInfo.Name(),
		appInfo.Version(),
		rt.WithAppInfo(appInfo), // Pass the created AppInfo
	)
	require.NoError(t, err, "Failed to initialize runtime from bootstrap")
	// Removed defer rtInstance.Cleanup() as it's no longer available
	// Load the configuration from the bootstrap file with all options.
	err = rtInstance.Load(bootstrapPath)
	require.NoError(t, err, "Failed to load configuration from bootstrap")
	defer rtInstance.Config().Close()

	var actualConfig testconfigs.TestConfig
	err = rtInstance.Config().Decode("", &actualConfig)
	require.NoError(t, err, "Failed to decode config from runtime")

	// Define the expected configuration after merging sources with priority.
	expectedApp := &appv1.App{
		Id:       "source1-app-id",
		Name:     "Source2App",
		Version:  "1.0.0",
		Env:      "prod",
		Metadata: map[string]string{"key3": "value3"},
	}

	expectedLogger := &loggerv1.Logger{
		Level:  "info",
		Format: "text",
	}

	expectedClient := &transportv1.Client{
		Grpc: &grpcv1.Client{
			Endpoint: "discovery:///source1-service",
			Timeout:  durationpb.New(5 * 1000 * 1000 * 1000), // 5s
			Selector: &selectorv1.SelectorConfig{
				Version: "v2.0.0",
			},
		},
	}

	// Perform assertions using the modular assertion toolkit.
	parentconfig.AssertAppConfig(t, rt.ConvertToAppInfo(expectedApp), rt.ConvertToAppInfo(actualConfig.App))
	parentconfig.AssertLoggerConfig(t, expectedLogger, actualConfig.Logger)
	parentconfig.AssertClientConfig(t, expectedClient, actualConfig.Client)

	t.Log("Multiple sources config loaded and merged successfully!")
}
