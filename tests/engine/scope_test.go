package engine_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/origadmin/runtime/engine/container"
	"github.com/origadmin/runtime/engine/metadata"
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
	// Visible to all because no Scopes are specified.
	reg.Register("middleware",
		func(root any) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "propagation", Value: nil}},
			}, nil
		},
		func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			if h.Scope() == metadata.ServerScope {
				return &ServerMiddleware{Name: "ServerProp"}, nil
			}
			if h.Scope() == metadata.ClientScope {
				return &ClientMiddleware{Name: "ClientProp"}, nil
			}
			return "GlobalProp", nil
		})

	// Initialize the container
	_ = reg.Init(ctx, "dummy_config")

	// 2. Retrieval from Server Scope
	hServer := reg.In("middleware", engine.WithScope(metadata.ServerScope))
	instServer, err := hServer.Get(ctx, "propagation")
	assert.NoError(t, err)
	assert.IsType(t, &ServerMiddleware{}, instServer)
	assert.Equal(t, "ServerProp", instServer.(*ServerMiddleware).Name)

	// 3. Retrieval from Client Scope
	hClient := reg.In("middleware", engine.WithScope(metadata.ClientScope))
	instClient, err := hClient.Get(ctx, "propagation")
	assert.NoError(t, err)
	assert.IsType(t, &ClientMiddleware{}, instClient)
	assert.Equal(t, "ClientProp", instClient.(*ClientMiddleware).Name)

	// 4. Verification of Physical Isolation
	assert.NotEqual(t, instServer, instClient)

	// 5. Explicit Override & Fallback
	reg.Register("middleware",
		func(root any) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "override", Value: nil}},
			}, nil
		},
		func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			return "GlobalOverride", nil
		},
		engine.WithScopes(metadata.GlobalScope))

	hScoped := reg.In("middleware", engine.WithScope(metadata.ServerScope))
	instOverride, err := hScoped.Get(ctx, "override")
	assert.NoError(t, err)
	assert.Equal(t, "GlobalOverride", instOverride)
}
