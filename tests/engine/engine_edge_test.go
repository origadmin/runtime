package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/engine"
)

// TestEngine_CircularDependency verifies detection of infinite instantiation loops
func TestEngine_CircularDependency(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// A -> B -> A
	reg.Register(runtime.CategoryClient, func(ctx context.Context, h engine.Handle) (any, error) {
		_, err := h.Locator().In(runtime.CategoryServer).Get(ctx, "B")
		return nil, err
	}, engine.WithDefaultEntry("A"))

	reg.Register(runtime.CategoryServer, func(ctx context.Context, h engine.Handle) (any, error) {
		_, err := h.Locator().In(runtime.CategoryClient).Get(ctx, "A")
		return nil, err
	}, engine.WithDefaultEntry("B"))

	_ = reg.Load(ctx, "src")

	_, err := reg.In(runtime.CategoryClient).Get(ctx, "A")
	if err == nil {
		t.Errorf("Should have detected circular dependency")
	}
	fmt.Printf("Detected expected circular error: %v\n", err)
}

// TestEngine_PriorityCompetition verifies that newer registrations with SAME priority 
// take precedence over older ones (standard stack behavior)
func TestEngine_PriorityCompetition(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// First registration
	reg.Register(runtime.CategoryLogger, func(ctx context.Context, h engine.Handle) (any, error) {
		return &mockComponent{Name: "old-logger"}, nil
	})

	// Second registration (same priority)
	reg.Register(runtime.CategoryLogger, func(ctx context.Context, h engine.Handle) (any, error) {
		return &mockComponent{Name: "new-logger"}, nil
	})

	_ = reg.Load(ctx, "src")

	inst, _ := reg.In(runtime.CategoryLogger).Get(ctx, "")
	if inst.(*mockComponent).Name != "new-logger" {
		t.Errorf("Expected newer registration to take precedence, got %v", inst.(*mockComponent).Name)
	}
}

// TestEngine_LifecycleProtection verifies that Register cannot be called after Load
func TestEngine_LifecycleProtection(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	_ = reg.Load(ctx, "src")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should have panicked when registering after Load")
		}
	}()

	reg.Register(runtime.CategoryLogger, simpleProvider)
}
