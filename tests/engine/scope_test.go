package engine_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/origadmin/runtime/engine/container"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
)

// Mock Middleware types
type ServerMiddleware struct{ Name string }
type ClientMiddleware struct{ Name string }

func TestScopeIsolationAndMatching(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()

	// 1. Unified Registration (Universal)
	reg.Register("middleware",
		func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			if h.Scope() == "server" {
				return &ServerMiddleware{Name: "ServerProp"}, nil
			}
			if h.Scope() == "client" {
				return &ClientMiddleware{Name: "ClientProp"}, nil
			}
			return "GlobalProp", nil
		},
		engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "propagation", Value: nil}},
			}, nil
		}))

	// 2. Load configuration into container
	_ = reg.Load(ctx, "dummy_config")

	// 3. Retrieval from Server Perspective
	hServer := reg.In("middleware", engine.WithInScope("server"))
	instServer, err := hServer.Get(ctx, "propagation")
	assert.NoError(t, err)
	assert.IsType(t, &ServerMiddleware{}, instServer)
	assert.Equal(t, "ServerProp", instServer.(*ServerMiddleware).Name)

	// 4. Retrieval from Client Perspective
	hClient := reg.In("middleware", engine.WithInScope("client"))
	instClient, err := hClient.Get(ctx, "propagation")
	assert.NoError(t, err)
	assert.IsType(t, &ClientMiddleware{}, instClient)
	assert.Equal(t, "ClientProp", instClient.(*ClientMiddleware).Name)

	// 5. Verification of Physical Isolation
	assert.NotEqual(t, instServer, instClient)

	// 6. Explicit Override & Fallback to _global
	reg.Register("middleware",
		func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			return "GlobalOverride", nil
		},
		engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "override", Value: nil}},
			}, nil
		}),
		engine.WithScopes(component.GlobalScope))

	// Re-load to trigger binding
	_ = reg.Load(ctx, "dummy_config")

	hScoped := reg.In("middleware", engine.WithInScope("server"))
	instOverride, err := hScoped.Get(ctx, "override")
	assert.NoError(t, err)
	assert.Equal(t, "GlobalOverride", instOverride)
}
