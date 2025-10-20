package custom_transformer

import (
	"fmt"

	appv1 "github.com/origadmin/runtime/api/gen/go/runtime/app/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/runtime/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// TestTransformer is a custom transformer for testing purposes.
type TestTransformer struct {
	interfaces.StructuredConfig
	Suffix string
	c      interfaces.Config
	cfg    *testconfigs.TestConfig
}

func (t *TestTransformer) Load() error {
	return t.c.Load()
}

func (t *TestTransformer) Decode(key string, value any) error {
	if err := t.c.Decode(key, value); err != nil {
		return fmt.Errorf("failed to decode config in transformer: %w", err)
	}
	return nil
}

func (t *TestTransformer) Raw() any {
	return t.c.Raw()
}

func (t *TestTransformer) Close() error {
	return t.c.Close()
}

func (t *TestTransformer) DecodeApp() (*appv1.App, error) {
	return t.cfg.App, nil
}

func (t *TestTransformer) DecodeLogger() (*loggerv1.Logger, error) {
	return t.cfg.Logger, nil
}

func (t *TestTransformer) DecodeMiddlewares() (*middlewarev1.Middlewares, error) {
	return t.cfg.Middlewares, nil
}

func (t *TestTransformer) Transform(c interfaces.Config, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
	t.c = c
	// Decode the current configuration into our TestConfig struct.
	var cfg testconfigs.TestConfig
	if err := c.Decode("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to scan config in transformer: %w", err)
	}

	// Perform the transformation: append suffix to App.Name.
	if cfg.App != nil {
		cfg.App.Name = cfg.App.Name + t.Suffix
	} else {
		// If App is nil, create it and set the name.
		cfg.App = &appv1.App{}
		cfg.App.Name = "TransformedApp" + t.Suffix
	}
	t.StructuredConfig = sc
	t.cfg = &cfg
	return t, nil
}

// String returns the name of the transformer.
func (t *TestTransformer) String() string {
	return "test-transformer"
}

// NewTestTransformer creates a new TestTransformer.
func NewTestTransformer(suffix string) *TestTransformer {
	return &TestTransformer{Suffix: suffix}
}

var _ bootstrap.ConfigTransformer = (*TestTransformer)(nil)
