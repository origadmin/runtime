package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
)

// TestEngine_Require verifies that category-based resolution can be easily 
// implemented via a RequirementResolver provided by the config.
func TestEngine_Require(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// 1. Register a dependency (e.g., Logger)
	reg.Register(component.CategoryLogger, func(ctx context.Context, h component.Handle) (any, error) {
		return &mockComponent{Name: "logger-impl"}, nil
	}, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "default", Value: "log-cfg"}},
		}, nil
	}))

	// 2. Register a dependent component (e.g., Server) that requires Logger
	reg.Register(component.CategoryServer, func(ctx context.Context, h component.Handle) (any, error) {
		// Use Require to get the logger
		logger, err := h.Require("log")
		if err != nil {
			return nil, err
		}
		return &mockComponent{Name: "server-impl", Dep: logger}, nil
	}, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{
				Name:  "default",
				Value: "srv-cfg",
			}},
			// The resolver provides the mapping logic
			RequirementResolver: func(ctx context.Context, h component.Handle, purpose string) (any, error) {
				if purpose == "log" {
					return h.Locator().In(component.CategoryLogger).Get(ctx, component.DefaultName)
				}
				return nil, fmt.Errorf("unknown requirement: %s", purpose)
			},
		}, nil
	}))

	// Load
	if err := reg.Load(ctx, "dummy-source"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// 4. Get the server and verify dependency
	srv, err := reg.In(component.CategoryServer).Get(ctx, "default")
	if err != nil {
		t.Fatalf("Failed to get server: %v", err)
	}

	mockSrv := srv.(*mockComponent)
	if mockSrv.Dep == nil {
		t.Fatal("Dependency 'log' was not injected via Require")
	}
	if mockLogger, ok := mockSrv.Dep.(*mockComponent); !ok || mockLogger.Name != "logger-impl" {
		t.Errorf("Injected dependency incorrect. Got %v, want logger-impl", mockSrv.Dep)
	}
}

// TestEngine_MiddlewareSelectorScenario demonstrates the Selector pattern where the resolution logic
// is silently provided by the Resolver via ConfigEntry.
// It verifies that the Carrier correctly contains other components but NOT the requester itself.
func TestEngine_MiddlewareSelectorScenario(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// 1. Agnostic Provider: It just tries to Require("carrier").
	middlewareProvider := func(ctx context.Context, h component.Handle) (any, error) {
		req, _ := h.Require("carrier")
		return &mockComponent{Name: h.Name(), Dep: req}, nil
	}
	reg.Register(component.CategoryMiddleware, middlewareProvider)

	// 2. Resolver: Defines which entries get the carrier resolution logic.
	carrierResolver := func(ctx context.Context, h component.Handle, purpose string) (any, error) {
		if purpose == "carrier" {
			carrier := make(map[string]any)
			// h.Locator() is already scoped and automatically skips the requester instance.
			for name, inst := range h.Locator().Iter(ctx) {
				carrier[name] = inst
			}
			return carrier, nil
		}
		return nil, fmt.Errorf("unknown requirement: %s", purpose)
	}

	reg.Register(component.CategoryMiddleware, nil, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{
				{Name: "auth", Value: "auth-cfg"},
				{Name: "log", Value: "log-cfg"},
				{Name: "selector", Value: "selector-cfg", RequirementResolver: carrierResolver},
			},
		}, nil
	}))

	// Load everything
	if err := reg.Load(ctx, "source-data"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// 3. Trigger instantiation. Getting selector will trigger its RequirementResolver and Iter.
	inst, err := reg.In(component.CategoryMiddleware).Get(ctx, "selector")
	if err != nil {
		t.Fatalf("Failed to get selector: %v", err)
	}

	// 4. Verification
	mockComp := inst.(*mockComponent)
	if mockComp.Dep == nil {
		t.Fatal("Selector should have a carrier, but got nil")
	}
	carrier := mockComp.Dep.(map[string]any)

	// MUST contain other components
	if _, ok := carrier["auth"]; !ok {
		t.Error("Carrier missing expected component: auth")
	}
	if _, ok := carrier["log"]; !ok {
		t.Error("Carrier missing expected component: log")
	}

	// MUST NOT contain itself
	if _, ok := carrier["selector"]; ok {
		t.Error("Carrier should NOT contain the requester (selector) itself")
	}

	// Verification of count
	if len(carrier) != 2 {
		t.Errorf("Expected exactly 2 components in carrier, got %d: %v", len(carrier), carrier)
	}
}

// TestEngine_Skip verifies the Skip functionality on Locator.
func TestEngine_Skip(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// Register multiple loggers
	reg.Register(component.CategoryLogger, simpleProvider, engine.WithResolverOption(func(source any, cat component.Category) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{
				{Name: "logger1", Value: "cfg1"},
				{Name: "logger2", Value: "cfg2"},
				{Name: "logger3", Value: "cfg3"},
			},
		}, nil
	}))

	if err := reg.Load(ctx, "dummy"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify iteration without skip
	count := 0
	for range reg.In(component.CategoryLogger).Iter(ctx) {
		count++
	}
	if count != 3 {
		t.Errorf("Expected 3 loggers, got %d", count)
	}

	// Verify iteration WITH Skip
	count = 0
	// Skip "logger2"
	iter := reg.In(component.CategoryLogger).Skip("logger2").Iter(ctx)
	for name := range iter {
		if name == "logger2" {
			t.Error("Skip failed: logger2 was returned in iteration")
		}
		count++
	}
	if count != 2 {
		t.Errorf("Expected 2 loggers after skipping one, got %d", count)
	}
}
