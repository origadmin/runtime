package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/container"
	"github.com/origadmin/runtime/helpers/comp"
)

// Mock Backend Types
type BackendMiddleware struct {
	Name     string
	Registry any
}

type backendConfig struct {
	Name string
}

func TestBackendMigrationSimulation(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()

	// 1. Registry (Infrastructure)
	reg.Register("infrastructure", func(ctx context.Context, h component.Handle) (any, error) {
		return "MockRegistry", nil
	}, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "default", Value: nil}},
		}, nil
	}), engine.WithPriority(100))

	// 2. Middleware
	reg.Register("middleware", func(ctx context.Context, h component.Handle) (any, error) {
		cfg := h.Config().(*backendConfig)
		regInst, err := h.Locator().In("infrastructure").Get(ctx, "")
		if err != nil {
			return nil, err
		}

		return &BackendMiddleware{Name: cfg.Name, Registry: regInst}, nil
	}, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "authn", Value: &backendConfig{Name: "authn-mw"}}},
		}, nil
	}), engine.WithPriority(500), engine.WithScopes("server"))

	// 3. Load configuration
	if err := reg.Load(ctx, "dummy_root"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// 4. Verify
	mwH := reg.In("middleware", engine.WithInScope("server"))
	m1, err := comp.Get[*BackendMiddleware](ctx, mwH, "authn")
	if err != nil {
		t.Fatalf("Failed to create authn: %v", err)
	}

	if m1.Name != "authn-mw" || m1.Registry != "MockRegistry" {
		t.Errorf("Unexpected results: %+v", m1)
	}

	t.Log("Backend migration simulation successful.")
}
