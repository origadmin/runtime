package config_target_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	_ "github.com/origadmin/runtime/config/file"
)

type ConfigTargetTestSuite struct {
	suite.Suite
}

func TestConfigTargetTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTargetTestSuite))
}

type MyConfig struct {
	App struct {
		Name string `json:"name"`
	} `json:"app"`
}

func (s *ConfigTargetTestSuite) TestConfigTargetBinding() {
	t := s.T()
	var target MyConfig
	rtInstance := rt.New("TargetTest", "1.0.0")
	// Since I'm using WithConfigTarget, the runtime will auto-bind the result to 'target'
	err := rtInstance.Load("bootstrap.yaml")
	require.NoError(t, err)
	defer rtInstance.Decoder().Close()

	// Direct scan to verify
	err = rtInstance.Decoder().Scan(&target)
	require.NoError(t, err)
	require.NotEmpty(t, target.App.Name)
}
