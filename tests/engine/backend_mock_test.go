package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/metadata"
)

// --- Mock Interfaces ---

type Authorizer interface {
	IsAllowed(ctx context.Context, sub, obj, act string) (bool, error)
}

type Skipper interface {
	ShouldSkip(ctx context.Context) bool
}

// --- Implementations ---

type mockAuthorizer struct {
	driver string
}

func (m *mockAuthorizer) IsAllowed(ctx context.Context, sub, obj, act string) (bool, error) {
	return true, nil
}

type mockSkipper struct{}

func (m *mockSkipper) ShouldSkip(ctx context.Context) bool { return false }

type mockMiddleware struct {
	name string
	auth Authorizer
	skip Skipper
}

// --- Configs ---

type BackendRootConfig struct {
	Security    *SecurityConfig
	Skipper     *SkipperConfig
	Middlewares *MiddlewareConfig
}

type SecurityConfig struct {
	Active  string
	Configs []engine.ConfigEntry
}

type SkipperConfig struct {
	Configs []engine.ConfigEntry
}

type MiddlewareConfig struct {
	Configs []engine.ConfigEntry
}

func TestBackendDeepDependencyInjection(t *testing.T) {
	ctx := context.Background()

	// 1. Truth
	root := &BackendRootConfig{
		Security: &SecurityConfig{
			Active: "casbin",
			Configs: []engine.ConfigEntry{
				{Name: "casbin", Value: "file-driver"},
			},
		},
		Skipper: &SkipperConfig{
			Configs: []engine.ConfigEntry{{Name: "default", Value: nil}},
		},
		Middlewares: &MiddlewareConfig{
			Configs: []engine.ConfigEntry{
				{Name: "authz-mw", Value: "authz-mw-cfg"},
			},
		},
	}

	reg := engine.NewContainer()

	// 2. Register

	reg.Register("security", func(r any) (*engine.ModuleConfig, error) {
		raw := r.(*BackendRootConfig).Security
		return &engine.ModuleConfig{Entries: raw.Configs, Active: raw.Active}, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		var driver string
		engine.BindConfig(h, &driver)
		return &mockAuthorizer{driver: driver}, nil
	}, engine.WithPriority(100))

	reg.Register("skipper", func(r any) (*engine.ModuleConfig, error) {
		raw := r.(*BackendRootConfig).Skipper
		return &engine.ModuleConfig{Entries: raw.Configs}, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		return &mockSkipper{}, nil
	}, engine.WithPriority(100))

	reg.Register(metadata.CategoryMiddleware, func(r any) (*engine.ModuleConfig, error) {
		raw := r.(*BackendRootConfig).Middlewares
		return &engine.ModuleConfig{Entries: raw.Configs}, nil
	}, func(ctx context.Context, h engine.Handle, opts ...options.Option) (any, error) {
		auth, err := engine.Cast[Authorizer](ctx, h.In("security"), "")
		if err != nil {
			return nil, err
		}
		skip, err := engine.Cast[Skipper](ctx, h.In("skipper"), "")
		if err != nil {
			return nil, err
		}
		return &mockMiddleware{name: "authz-mw", auth: auth, skip: skip}, nil
	}, engine.WithPriority(500), engine.WithScope(metadata.ServerScope))

	// 3. Activate Engine (SINGLE entrance)
	if err := reg.Init(ctx, root); err != nil {
		t.Fatalf("Engine init failed: %v", err)
	}

	// 4. Verify
	mwH := reg.In(metadata.CategoryMiddleware, engine.WithScope(metadata.ServerScope))
	mw, err := engine.Cast[*mockMiddleware](ctx, mwH, "authz-mw")
	if err != nil {
		t.Fatalf("Failed to create middleware stack: %v", err)
	}

	if mw.auth.(*mockAuthorizer).driver != "file-driver" {
		t.Errorf("Expected driver 'file-driver', got %s", mw.auth.(*mockAuthorizer).driver)
	}

	t.Log("Successfully verified correct Registration -> Init -> Get flow.")
}
