package custom_transformer

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/config"
)

type TestConfig struct {
	App *appv1.App `json:"app"`
}

func (c *TestConfig) GetApp() *appv1.App {
	return c.App
}

type TestTransformer struct{}

func (t *TestTransformer) Transform(c config.KConfig) (any, error) {
	var cfg TestConfig
	if err := c.Scan(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

var _ ConfigTransformer = (*TestTransformer)(nil)

type ConfigTransformer interface {
	Transform(config.KConfig) (any, error)
}
