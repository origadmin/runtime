package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/origadmin/runtime"
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/middleware"
)

// TestEngine_Require verifies that category-based resolution can be easily 
// implemented via a RequirementResolver provided by the config.
func TestEngine_Require(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// 1. Register a dependency (e.g., Logger)
	reg.Register(runtime.CategoryLogger, func(ctx context.Context, h component.Handle) (any, error) {
		return &mockComponent{Name: "logger-impl"}, nil
	}, engine.WithConfigResolverOption(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{Name: "default", Value: "log-cfg"}},
		}, nil
	}))

	// 2. Register a dependent component (e.g., Server) that requires Logger
	reg.Register(runtime.CategoryServer, func(ctx context.Context, h component.Handle) (any, error) {
		// Use Require to get the logger
		logger, err := h.Require("log")
		if err != nil {
			return nil, err
		}
		return &mockComponent{Name: "server-impl", Dep: logger}, nil
	}, engine.WithConfigResolverOption(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{
				Name:  "default",
				Value: "srv-cfg",
			}},
			// The resolver provides the mapping logic
			RequirementResolver: func(ctx context.Context, h component.Handle, purpose string) (any, error) {
				if purpose == "log" {
					return h.Locator().In(runtime.CategoryLogger).Get(ctx, "")
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
	srv, err := reg.In(runtime.CategoryServer).Get(ctx, "default")
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

	// 1. Agnostic Provider: It just tries to Require(RequirementCarrier).
	middlewareProvider := func(ctx context.Context, h component.Handle) (any, error) {
		req, _ := h.Require(middleware.RequirementCarrier)
		return &mockComponent{Name: h.Name(), Dep: req}, nil
	}
	reg.Register(runtime.CategoryMiddleware, middlewareProvider)

	// 2. Resolver: Defines which entries get the carrier resolution logic.
	carrierResolver := func(ctx context.Context, h component.Handle, purpose string) (any, error) {
		if purpose == middleware.RequirementCarrier {
			carrier := make(map[string]any)
			it := h.Locator().Iter(ctx)
			for it.Next() {
				name, inst := it.Value()
				carrier[name] = inst
			}
			return carrier, it.Err()
		}
		return nil, fmt.Errorf("unknown requirement: %s", purpose)
	}

	reg.Register(runtime.CategoryMiddleware, nil, engine.WithConfigResolverOption(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
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
	inst, err := reg.In(runtime.CategoryMiddleware).Get(ctx, "selector")
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

// TestEngine_MiddlewareSelectorSequence specifically tests the sequence issue where the Selector 
// is the first entry and its Requirement triggers the instantiation of all subsequent components.
func TestEngine_MiddlewareSelectorSequence(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// Track instantiation order
	var instantiated []string

	middlewareProvider := func(ctx context.Context, h component.Handle) (any, error) {
		instantiated = append(instantiated, h.Name())
		req, _ := h.Require(middleware.RequirementCarrier)
		return &mockComponent{Name: h.Name(), Dep: req}, nil
	}
	reg.Register(runtime.CategoryMiddleware, middlewareProvider)

	carrierResolver := func(ctx context.Context, h component.Handle, purpose string) (any, error) {
		if purpose == middleware.RequirementCarrier {
			carrier := make(map[string]any)
			it := h.Locator().Iter(ctx)
			for it.Next() {
				name, inst := it.Value()
				carrier[name] = inst
			}
			return carrier, it.Err()
		}
		return nil, fmt.Errorf("unknown requirement: %s", purpose)
	}

	reg.Register(runtime.CategoryMiddleware, nil, engine.WithConfigResolverOption(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{
				// Selector is FIRST
				{Name: "selector", Value: "selector-cfg", RequirementResolver: carrierResolver},
				{Name: "auth", Value: "auth-cfg"},
				{Name: "log", Value: "log-cfg"},
			},
		}, nil
	}))

	// Load everything
	if err := reg.Load(ctx, "source-data"); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// 3. Directly get the selector FIRST.
	// This should trigger: selector provider -> Require("carrier") -> Iter -> instantiate(auth), instantiate(log)
	inst, err := reg.In(runtime.CategoryMiddleware).Get(ctx, "selector")
	if err != nil {
		t.Fatalf("Failed to get selector: %v", err)
	}

	mockComp := inst.(*mockComponent)
	carrier := mockComp.Dep.(map[string]any)

	// Verify instantiation occurred for others
	if _, ok := carrier["auth"]; !ok {
		t.Error("Carrier missing 'auth', instantiation might not have been triggered")
	}
	if _, ok := carrier["log"]; !ok {
		t.Error("Carrier missing 'log', instantiation might not have been triggered")
	}

	// Verify sequence: selector must be the first one to start instantiating
	if len(instantiated) < 1 || instantiated[0] != "selector" {
		t.Errorf("Expected 'selector' to be the first in instantiation sequence, got %v", instantiated)
	}

	// The full sequence should be [selector, auth, log] or [selector, log, auth]
	if len(instantiated) != 3 {
		t.Errorf("Expected 3 components to be instantiated, got %d: %v", len(instantiated), instantiated)
	}
}

// TestEngine_Skip verifies the Skip functionality on Locator.
func TestEngine_Skip(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// Register multiple loggers
	reg.Register(runtime.CategoryLogger, simpleProvider, engine.WithConfigResolverOption(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
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
	it := reg.In(runtime.CategoryLogger).Iter(ctx)
	for it.Next() {
		count++
	}
	if count != 3 {
		t.Errorf("Expected 3 loggers, got %d", count)
	}

	// Verify iteration WITH Skip
	count = 0
	// Skip "logger2"
	itSkip := reg.In(runtime.CategoryLogger).Skip("logger2").Iter(ctx)
	for itSkip.Next() {
		name, _ := itSkip.Value()
		if name == "logger2" {
			t.Error("Skip failed: logger2 was returned in iteration")
		}
		count++
	}
	if count != 2 {
		t.Errorf("Expected 2 loggers after skipping one, got %d", count)
	}
}
