package env_specific_config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/bootstrap"
	_ "github.com/origadmin/runtime/config/file"
)

type EnvSpecificConfigTestSuite struct {
	suite.Suite
}

func TestEnvSpecificConfigTestSuite(t *testing.T) {
	suite.Run(t, new(EnvSpecificConfigTestSuite))
}

var meta struct {
	App *appv1.App `json:"app" yaml:"app"`
}

func (s *EnvSpecificConfigTestSuite) TestEnvSpecificLoading() {
	t := s.T()
	bootstrapPath := "bootstrap_env.yaml"

	testCases := []struct {
		env          string
		expectedName string
	}{
		{"dev", "DevApp"},
		{"prod", "ProdApp"},
	}

	for _, tc := range testCases {
		t.Run(tc.env, func(t *testing.T) {
			os.Setenv("APP_ENV", tc.env)
			defer os.Unsetenv("APP_ENV")

			rtInstance := rt.New("EnvTest", "1.0.0")
			err := rtInstance.Load(bootstrapPath, bootstrap.WithConfigTarget(&meta))
			require.NoError(t, err)
			defer rtInstance.Config().Close()

			require.Equal(t, tc.expectedName, rtInstance.AppInfo().Name)
		})
	}
}
