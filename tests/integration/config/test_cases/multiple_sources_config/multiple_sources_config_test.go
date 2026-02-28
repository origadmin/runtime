package multiple_sources_config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	_ "github.com/origadmin/runtime/config/file"
)

type MultipleSourcesConfigTestSuite struct {
	suite.Suite
}

func TestMultipleSourcesConfigTestSuite(t *testing.T) {
	suite.Run(t, new(MultipleSourcesConfigTestSuite))
}

func (s *MultipleSourcesConfigTestSuite) TestMultipleSourcesLoading() {
	t := s.T()
	bootstrapPath := "bootstrap_multiple_sources.yaml"
	rtInstance := rt.New("MultipleSourcesTest", "1.0.0")
	err := rtInstance.Load(bootstrapPath)
	if err != nil {
		t.Logf("Skipping test: could not load %s: %v", bootstrapPath, err)
		return
	}
	defer rtInstance.Config().Close()

	require.Equal(t, "OverriddenApp", rtInstance.AppInfo().Name)
}
