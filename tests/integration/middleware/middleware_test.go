package middleware_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	_ "github.com/origadmin/runtime/config/file"
	"github.com/origadmin/runtime/engine/bootstrap"
)

type MiddlewareTestSuite struct {
	suite.Suite
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) TestLoadAndBuild() {
	t := s.T()

	// 1. Initialize and load the runtime App
	rtInstance := rt.New("MiddlewareTest", "1.0.0")
	err := rtInstance.Load("configs/config.yaml", bootstrap.WithDirectly(true))
	require.NoError(t, err)
	defer rtInstance.Config().Close()

	// 2. Decode and verify using native KConfig Scan
	var middlewareConfig middlewarev1.Middlewares
	err = rtInstance.Config().Value("middlewares").Scan(&middlewareConfig)
	require.NoError(t, err)
	require.NotEmpty(t, middlewareConfig.Configs)

	t.Logf("Loaded middlewares: %v", middlewareConfig.Configs)
}
