package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
)

// TestEngine_BasicBinding verifies the fundamental Register -> Load -> Get flow
func TestEngine_BasicBinding(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// Register a simple logger
	reg.Register(component.CategoryLogger, simpleProvider)

	// Load with a direct source string
	err := reg.Load(ctx, "my-logger-config")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Get default instance
	inst, err := reg.In(component.CategoryLogger).Get(ctx, "")
	if err != nil {
		t.Fatalf("Failed to get logger: %v", err)
	}

	if inst.(*mockComponent).Name != "logger" {
		t.Errorf("Expected dynamic name 'logger' (category fallback), got %v", inst.(*mockComponent).Name)
	}
}

// TestEngine_NamingPriority verifies the naming logic: Named "default" > Active > Single Entry
func TestEngine_NamingPriority(t *testing.T) {
	ctx := context.Background()

	t.Run("Explicit_Named_Default", func(t *testing.T) {
		reg := engine.NewContainer()
		reg.Register(component.CategoryCache, simpleProvider, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{
					{Name: "redis", Value: "cfg1"},
					{Name: "default", Value: "cfg2"},
				},
			}, nil
		}))
		_ = reg.Load(ctx, "src")
		val, _ := reg.In(component.CategoryCache).Get(ctx, "")
		if val.(*mockComponent).Name != "default" {
			t.Errorf("Should pick explicitly named 'default', got %v", val.(*mockComponent).Name)
		}
	})

	t.Run("Active_Override", func(t *testing.T) {
		reg := engine.NewContainer()
		reg.Register(component.CategoryCache, simpleProvider, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{
					{Name: "redis", Value: "cfg1"},
					{Name: "memcached", Value: "cfg2"},
				},
				Active: "memcached",
			}, nil
		}))
		_ = reg.Load(ctx, "src")
		val, _ := reg.In(component.CategoryCache).Get(ctx, "")
		if val.(*mockComponent).Name != "memcached" {
			t.Errorf("Should pick explicit Active entry, got %v", val.(*mockComponent).Name)
		}
	})

	t.Run("Single_Entry_Fallback", func(t *testing.T) {
		reg := engine.NewContainer()
		reg.Register(component.CategoryCache, simpleProvider, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{{Name: "only_one", Value: "cfg1"}},
			}, nil
		}))
		_ = reg.Load(ctx, "src")
		val, _ := reg.In(component.CategoryCache).Get(ctx, "")
		if val.(*mockComponent).Name != "only_one" {
			t.Errorf("Should pick the only available entry, got %v", val.(*mockComponent).Name)
		}
	})

	t.Run("Multiple_Should_Fail_Default", func(t *testing.T) {
		reg := engine.NewContainer()
		reg.Register(component.CategoryCache, simpleProvider, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
			return &component.ModuleConfig{
				Entries: []component.ConfigEntry{
					{Name: "redis", Value: "cfg1"},
					{Name: "memcached", Value: "cfg2"},
				},
			}, nil
		}))
		_ = reg.Load(ctx, "src")
		_, err := reg.In(component.CategoryCache).Get(ctx, "")
		if err == nil {
			t.Errorf("Should return error for default request when multiple entries exist without explicit default")
		}
	})
}
