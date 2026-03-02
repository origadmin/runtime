package custom_transformer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	runtimeconfig "github.com/origadmin/runtime/config"
	_ "github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/engine/bootstrap"
)

type CustomTransformerTestSuite struct {
	suite.Suite
}

func TestCustomTransformerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomTransformerTestSuite))
}

type TransformedConfig struct {
	AppName string
}

func (s *CustomTransformerTestSuite) TestCustomTransformerApplication() {
	t := s.T()

	transformer := bootstrap.ConfigTransformFunc(func(cfg runtimeconfig.KConfig) (any, error) {
		var raw struct {
			App struct {
				Name string `json:"name"`
			} `json:"app"`
		}
		// Directly cast to runtimeconfig.KConfig
		if err := cfg.Scan(&raw); err != nil {
			return nil, err
		}
		return &TransformedConfig{AppName: raw.App.Name + "-transformed"}, nil
	})

	rtInstance := rt.New("TransformerTest", "1.0.0")
	err := rtInstance.Load("bootstrap_transformer.yaml", bootstrap.WithConfigTransformer(transformer))
	require.NoError(t, err)
	defer rtInstance.Config().Close()

	res := rtInstance.Result().Config()
	transformed, ok := res.(*TransformedConfig)
	require.True(t, ok)
	require.Contains(t, transformed.AppName, "-transformed")
}
