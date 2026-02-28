package consul_source_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	_ "github.com/origadmin/runtime/config/file"
)

type ConsulSourceTestSuite struct {
	suite.Suite
}

func TestConsulSourceTestSuite(t *testing.T) {
	suite.Run(t, new(ConsulSourceTestSuite))
}

func (s *ConsulSourceTestSuite) TestConsulSourceLoading() {
	t := s.T()
	bootstrapPath := "bootstrap_consul.yaml"
	rtInstance := rt.New("ConsulSourceTest", "1.0.0")
	err := rtInstance.Load(bootstrapPath)
	if err != nil {
		t.Logf("Expected error (no Consul server): %v", err)
		return
	}
	defer rtInstance.Config().Close()

	configDecoder := rtInstance.Config()
	s.NotNil(configDecoder)
	var appConfig struct {
		Name string `json:"name"`
	}
	err = configDecoder.Value("app").Scan(&appConfig)
	s.NoError(err)
}
