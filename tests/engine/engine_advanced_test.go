package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
)

// TestEngine_InjectFeature verifies direct instance injection
func TestEngine_InjectFeature(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// 1. Inject a named instance
	myInst := &mockComponent{Name: "manual"}
	reg.Inject(runtime.CategoryLogger, "manual_logger", myInst)

	// 2. Inject a default instance (with higher priority)
	defaultInst := &mockComponent{Name: "default_injected"}
	reg.Inject(runtime.CategoryLogger, "", defaultInst)

	// Verify manual logger
	inst, err := reg.In(runtime.CategoryLogger).Get(ctx, "manual_logger")
	if err != nil {
		t.Fatalf("Failed to get injected logger: %v", err)
	}
	if inst.(*mockComponent).Name != "manual" {
		t.Errorf("Expected 'manual', got %v", inst.(*mockComponent).Name)
	}

	// Verify default injected logger
	inst, err = reg.In(runtime.CategoryLogger).Get(ctx, "")
	if err != nil {
		t.Fatalf("Failed to get default injected logger: %v", err)
	}
	if inst.(*mockComponent).Name != "default_injected" {
		t.Errorf("Expected 'default_injected', got %v", inst.(*mockComponent).Name)
	}
}

// TestEngine_SeedingFeature verifies DefaultEntries can pre-seed entries even without config
func TestEngine_SeedingFeature(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	reg.Register(runtime.CategoryDatabase, simpleProvider, engine.WithDefaultEntries("predefined_db"))

	// Load with EMPTY source
	err := reg.Load(ctx, nil)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should still find predefined_db
	inst, err := reg.In(runtime.CategoryDatabase).Get(ctx, "predefined_db")
	if err != nil {
		t.Fatalf("Seeding failed: %v", err)
	}
	if inst.(*mockComponent).Name != "predefined_db" {
		t.Errorf("Expected 'predefined_db', got %v", inst.(*mockComponent).Name)
	}
}

// TestEngine_Iteration verifies Iter() functionality
func TestEngine_Iteration(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	reg.Register(runtime.CategoryCache, simpleProvider, engine.WithConfigResolverOption(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{
				{Name: "c1", Value: "v1"},
				{Name: "c2", Value: "v2"},
			},
		}, nil
	}))

	_ = reg.Load(ctx, "src")

	count := 0
	names := make(map[string]bool)
	it := reg.In(runtime.CategoryCache).Iter(ctx)
	for it.Next() {
		name, _ := it.Value()
		count++
		names[name] = true
	}
	if err := it.Err(); err != nil {
		t.Fatalf("Iteration failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 items in iteration, got %v", count)
	}
	if !names["c1"] || !names["c2"] {
		t.Errorf("Iteration missing expected names")
	}
}
