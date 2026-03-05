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

func TestScopeIsolationAndPerspectiveSwitch(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()

	// 1. Register a component in BOTH server and global scopes
	// They use the SAME provider but will result in DIFFERENT instances
	reg.Register("middleware",
		func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			if h.Scope() == "server" {
				return &ServerMiddleware{Name: "ServerInst"}, nil
			}
			return "GlobalInst", nil
		},
		engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "item", Value: nil}},
			}, nil
		}),
		engine.WithScopes("server", component.GlobalScope))

	// 2. Load configuration
	_ = reg.Load(ctx, "root")

	// 3. Current Perspective: Server
	hServer := reg.In("middleware", engine.WithInScope("server"))

	// A. Should get server instance
	inst1, err := hServer.Get(ctx, "item")
	assert.NoError(t, err)
	assert.Equal(t, "ServerInst", inst1.(*ServerMiddleware).Name)

	// B. Strict Isolation: Should NOT find global things if they aren't in this scope
	// (In this test they are both, but they are separate instances)

	// 4. Perspective Switch: Explicitly move to Global
	// Use hServer.In to switch category/scope perspective
	hGlobal := hServer.In("middleware") // Default is GlobalScope
	assert.Equal(t, component.GlobalScope, hGlobal.Scope())

	inst2, err := hGlobal.Get(ctx, "item")
	assert.NoError(t, err)
	assert.Equal(t, "GlobalInst", inst2)
	
	// 5. Verify they are distinct
	assert.NotEqual(t, inst1, inst2)
}

func TestContainer_LifecycleProtection(t *testing.T) {
	reg := container.NewContainer()
	ctx := context.Background()

	reg.Register("test", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		return "ok", nil
	})

	_ = reg.Load(ctx, "root")

	// Subsequent registration must panic
	assert.Panics(t, func() {
		reg.Register("late", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			return "bad", nil
		})
	})
}
