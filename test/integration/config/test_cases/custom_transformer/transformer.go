package custom_transformer

import (
	"fmt"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
)

// TestTransformer is a custom transformer for testing purposes.
type TestTransformer struct {
	interfaces.StructuredConfig
	Suffix string
	c      interfaces.ConfigLoader
	cfg    *testconfigs.TestConfig
}

func (t *TestTransformer) DecodedConfig() any {
	return t.cfg
}

func (t *TestTransformer) Transform(c interfaces.ConfigLoader, sc interfaces.StructuredConfig) (interfaces.StructuredConfig, error) {
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
