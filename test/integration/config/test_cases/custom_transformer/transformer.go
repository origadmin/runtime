package custom_transformer

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/source/memory"
	"github.com/origadmin/runtime/bootstrap"
	testconfigs "github.com/origadmin/runtime/test/integration/config/proto"
	"gopkg.in/yaml.v3"
)

// TestTransformer is a custom transformer for testing purposes.
type TestTransformer struct {
	suffix string
}

// NewTestTransformer creates a new TestTransformer.
func NewTestTransformer(suffix string) *TestTransformer {
	return &TestTransformer{suffix: suffix}
}

// Transform implements the bootstrap.Transformer interface.
func (t *TestTransformer) Transform(ctx context.Context, c config.Config) (config.Config, error) {
	// Decode the current configuration into our TestConfig struct.
	var cfg testconfigs.TestConfig
	if err := c.Scan(&cfg); err != nil {
		return nil, fmt.Errorf("failed to scan config in transformer: %w", err)
	}

	// Perform the transformation: append suffix to App.Name.
	if cfg.App != nil {
		cfg.App.Name = cfg.App.Name + t.suffix
	} else {
		// If App is nil, create it and set the name.
		cfg.App = &testconfigs.App{}
		cfg.App.Name = "TransformedApp" + t.suffix
	}

	// Re-encode the modified struct back to a map[string]interface{} for memory source.
	var modifiedMap map[string]interface{}
	modifiedBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified config in transformer: %w", err)
	}
	if err := yaml.Unmarshal(modifiedBytes, &modifiedMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal modified bytes to map: %w", err)
	}

	// Create a new config.Config from the modified map using a memory source.
	newConfig := config.New(config.WithSource(memory.NewSource(modifiedMap)))
	if err := newConfig.Load(); err != nil {
		return nil, fmt.Errorf("failed to load new config from memory source: %w", err)
	}

	return newConfig, nil
}

// String returns the name of the transformer.
func (t *TestTransformer) String() string {
	return "test-transformer"
}

// Register the custom transformer with the bootstrap package.
func init() {
	bootstrap.RegisterTransformer("test-transformer", func(options map[string]interface{}) (bootstrap.Transformer, error) {
		suffix, ok := options["suffix"].(string)
		if !ok {
			return nil, fmt.Errorf("transformer 'test-transformer' requires a 'suffix' option of type string")
		}
		return NewTestTransformer(suffix), nil
	})
}
