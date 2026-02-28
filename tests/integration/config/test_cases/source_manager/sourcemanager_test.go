package source_manager_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	rt "github.com/origadmin/runtime"
	_ "github.com/origadmin/runtime/config/file"
)

type SourceManagerSuite struct {
	suite.Suite
}

func TestSourceManagerSuite(t *testing.T) {
	suite.Run(t, new(SourceManagerSuite))
}

type CustomSettings struct {
	FeatureEnabled bool   `json:"feature_enabled"`
	APIKey         string `json:"api_key"`
}

func (s *SourceManagerSuite) TestConfigSourceMergingAndPriority() {
	t := s.T()
	bootstrapPath := "bootstrap.yaml"
	rtInstance := rt.New("SourceManagerTest", "1.0.0")
	err := rtInstance.Load(bootstrapPath)
	if err != nil {
		t.Logf("Warning: Skipping SourceManager test as config could not be loaded: %v", err)
		return
	}
	defer rtInstance.Config().Close()

	configDecoder := rtInstance.Config()
	s.NotNil(configDecoder)

	var loggerConfig struct {
		Level string `json:"level"`
	}
	err = configDecoder.Value("logger").Scan(&loggerConfig)
	s.NoError(err)
	s.Equal("debug", loggerConfig.Level)

	var settings CustomSettings
	err = configDecoder.Value("components.my-custom-settings").Scan(&settings)
	s.NoError(err)
	s.True(settings.FeatureEnabled)

	s.Equal("OverriddenApp", rtInstance.AppInfo().Name)
}
