package engine_test

import (
	"context"
	"testing"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/metadata"
)

// --- 1. Simulate Backend's Strong-Typed Config ---

type BackendConfig struct {
	Middlewares *middlewarev1.Middlewares
	Discovery   *discoveryv1.Discovery
}

// --- 2. Simulating Middleware Implementation ---

type BackendMiddleware struct {
	Name     string
	Type     string
	Registry string
}

func TestBackendMigrationSimulation(t *testing.T) {
	ctx := context.Background()

	// Step 1: Initialize Root Config
	root := &BackendConfig{
		Middlewares: &middlewarev1.Middlewares{
			Configs: []*middlewarev1.Middleware{
				{Name: "authn", Type: "jwt", Enabled: true},
				{Name: "authz", Type: "casbin", Enabled: true},
			},
		},
		Discovery: &discoveryv1.Discovery{
			Name: "consul",
			Type: "consul",
		},
	}

	// Step 2: Initialize Container
	reg := engine.NewContainer()

	// Step 3: Register Registry Factory
	reg.Register(metadata.CategoryRegistry, func(r any) (*engine.ModuleConfig, error) {
		raw := r.(*BackendConfig).Discovery
		return &engine.ModuleConfig{
			Entries: []engine.ConfigEntry{{Name: raw.Name, Value: raw}},
			Active:  raw.Name,
		}, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		var cfg discoveryv1.Discovery
		engine.BindConfig(h, &cfg)
		return cfg.Name, nil
	}, engine.WithPriority(metadata.PriorityRegistry))

	// Step 4: Register Middleware Factory
	reg.Register(metadata.CategoryMiddleware, func(r any) (*engine.ModuleConfig, error) {
		raw := r.(*BackendConfig).Middlewares
		var entries []engine.ConfigEntry
		for _, c := range raw.Configs {
			entries = append(entries, engine.ConfigEntry{Name: c.Name, Value: c})
		}
		return &engine.ModuleConfig{Entries: entries}, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		var cfg middlewarev1.Middleware
		engine.BindConfig(h, &cfg)

		regInst, err := engine.GetDefault[string](ctx, h.In(metadata.CategoryRegistry))
		if err != nil {
			return nil, err
		}

		return &BackendMiddleware{Name: cfg.Name, Registry: regInst}, nil
	}, engine.WithPriority(metadata.PriorityServerStack), engine.WithScopes(metadata.ServerScope))

	// Step 5: Activate and Verify
	if err := reg.Init(ctx, root); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mwH := reg.In(metadata.CategoryMiddleware, engine.WithScope(metadata.ServerScope))

	m1, err := engine.Cast[*BackendMiddleware](ctx, mwH, "authn")
	if err != nil {
		t.Fatalf("Failed to create authn: %v", err)
	}
	if m1.Name != "authn" {
		t.Errorf("Expected authn, got %s", m1.Name)
	}

	t.Log("Backend migration simulation successful.")
}
