package bootstrap_load_config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	_ "github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/engine/bootstrap"
	testconfigs "github.com/origadmin/runtime/tests/integration/config/proto"
)

type RuntimeIntegrationTestSuite struct {
	suite.Suite
}

func TestRuntimeIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RuntimeIntegrationTestSuite))
}

func (s *RuntimeIntegrationTestSuite) TestConfigProtoIntegration() {
	t := s.T()
	// This path contains config.yaml with app.name: ProtoConfigApp
	bootstrapPath := "testdata/proto_integration/bootstrap.yaml"
	rtInstance := rt.New("IntegrationTest", "1.0.0")

	// Explicitly bind a structure that implements AppConfig to enable app.name override from config.yaml
	var actualConfig testconfigs.TestConfig
	err := rtInstance.Load(bootstrapPath, bootstrap.WithConfigTarget(&actualConfig))
	require.NoError(t, err)
	defer rtInstance.Config().Close()

	// Verification: Metadata should be overridden by the business configuration
	require.Equal(t, "ProtoConfigApp", rtInstance.AppInfo().Name)
}

func (s *RuntimeIntegrationTestSuite) TestRuntimeDecoder() {
	t := s.T()
	// This path contains config.yaml with app.name: DecoderApp
	bootstrapPath := "testdata/decoder_test/bootstrap.yaml"
	rtInstance := rt.New("DecoderTest", "1.0.0")
	err := rtInstance.Load(bootstrapPath)
	require.NoError(t, err)

	cfg := rtInstance.Config()
	require.NotNil(t, cfg)
	defer cfg.Close()

	// Correct way: Provide an explicit binding structure for Scan to extract data
	var result struct {
		App *appv1.App `json:"app"`
	}
	err = cfg.Scan(&result)
	require.NoError(t, err)
	require.NotNil(t, result.App)
	require.Equal(t, "DecoderApp", result.App.Name)
}

func (s *RuntimeIntegrationTestSuite) TestRuntimeLoadCompleteConfig() {
	t := s.T()
	// This path contains config.yaml with app.name: CompleteApp
	bootstrapPath := "testdata/complete_config/bootstrap.yaml"
	rtInstance := rt.New("TestCompleteConfig", "1.0.0")

	// Note: Without WithConfigTarget, content in config.yaml is treated as "unbound config" and ignored.
	// With WithConfigTarget, metadata will be merged according to the priority.
	var actualConfig testconfigs.TestConfig
	err := rtInstance.Load(bootstrapPath, bootstrap.WithConfigTarget(&actualConfig))
	require.NoError(t, err)
	defer rtInstance.Config().Close()

	// Final verification: In a complete loading flow, the app name is overridden by CompleteApp
	require.Equal(t, "CompleteApp", rtInstance.AppInfo().Name)

	var actualConfig2 testconfigs.TestConfig
	// Scan into the actual configuration structure
	err = rtInstance.Config().Scan(&actualConfig2)
	require.NoError(t, err)
	require.NotNil(t, actualConfig2.App)
	require.Equal(t, "CompleteApp", actualConfig2.App.Name)
}
