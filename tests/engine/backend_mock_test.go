package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/container"
)

// --- Mock Domain Interfaces ---

type (
	Authenticator interface{ Authenticate() bool }
	Skipper       interface{ Skip() bool }
)

type mockAuthn struct{}
func (m *mockAuthn) Authenticate() bool { return true }

type mockSkipper struct{}
func (m *mockSkipper) Skip() bool { return true }

type mockMiddleware struct {
	name string
	auth Authenticator
	skip Skipper
}

func TestBackendDeepDependencyInjection(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()
	root := "dummy_config"

	// 1. Register Infrastructure (Authn)
	reg.Register("infrastructure", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		return &mockAuthn{}, nil
	}, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: "jwt", Value: nil}}}, nil
	}))

	// 2. Register Skipper
	reg.Register("skipper", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		return &mockSkipper{}, nil
	}, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: "default", Value: nil}}}, nil
	}))

	// 3. Register Middleware (Complex DI)
	reg.Register("middleware", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		// Deep Dependency Discovery
		auth, err := engine.Get[Authenticator](ctx, h.In("infrastructure"), "jwt")
		if err != nil {
			return nil, err
		}
		skip, err := engine.Get[Skipper](ctx, h.In("skipper"), "")
		if err != nil {
			return nil, err
		}
		return &mockMiddleware{name: "authz-mw", auth: auth, skip: skip}, nil
	}, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: "authz-mw", Value: nil}}}, nil
	}), engine.WithPriority(500), engine.WithScopes("server"))

	// 4. Load Engine
	if err := reg.Load(ctx, root); err != nil {
		t.Fatalf("Engine load failed: %v", err)
	}

	// 5. Verify
	mwH := reg.In("middleware", engine.WithInScope("server"))
	mw, err := engine.Get[*mockMiddleware](ctx, mwH, "authz-mw")
	if err != nil {
		t.Fatalf("Failed to create middleware stack: %v", err)
	}

	if mw.auth == nil || mw.skip == nil {
		t.Errorf("DI failed: auth=%v, skip=%v", mw.auth, mw.skip)
	}

	t.Log("Successfully verified correct Registration -> Load -> Get flow.")
}
