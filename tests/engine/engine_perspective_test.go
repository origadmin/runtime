package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
)

// TestEngine_PerspectiveDuality simulates the Gateway/Service directional requirement
func TestEngine_PerspectiveDuality(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	// Provider that behaves differently based on Work Tag (h.Tag())
	reg.Register(component.CategoryMiddleware, func(ctx context.Context, h engine.Handle) (any, error) {
		tag := h.Tag()
		scope := h.Scope()

		if scope == component.ServerScope {
			if tag == "gateway" {
				return &mockComponent{Type: "root-none"}, nil
			}
			if tag == "service" {
				return &mockComponent{Type: "node-extracter"}, nil
			}
		}
		if scope == component.ClientScope {
			return &mockComponent{Type: "any-injecter"}, nil
		}
		return nil, nil
	}, engine.WithDefaultEntry("propagation"), engine.WithScopes(component.ServerScope, component.ClientScope))

	_ = reg.Load(ctx, "source")

	// CASE A: Gateway Identity
	gP := reg.In(component.CategoryMiddleware, engine.WithInTags("gateway"))
	gS, _ := gP.In(component.CategoryMiddleware, engine.WithInScope(component.ServerScope)).Get(ctx, "propagation")
	if gS.(*mockComponent).Type != "root-none" {
		t.Errorf("Gateway-Server identity failure: expected root-none, got %v", gS.(*mockComponent).Type)
	}

	// CASE B: Service Identity
	sP := reg.In(component.CategoryMiddleware, engine.WithInTags("service"))
	sS, _ := sP.In(component.CategoryMiddleware, engine.WithInScope(component.ServerScope)).Get(ctx, "propagation")
	if sS.(*mockComponent).Type != "node-extracter" {
		t.Errorf("Service-Server identity failure: expected node-extracter, got %v", sS.(*mockComponent).Type)
	}

	// Verify Isolation: In the same container, they are different instances
	if gS == sS {
		t.Errorf("Identity Isolation failure: Gateway and Service instances should be different")
	}
}

// TestEngine_ScopeIsolation verifies that components in Client/Server scopes are strictly separate
func TestEngine_ScopeIsolation(t *testing.T) {
	ctx := context.Background()
	reg := engine.NewContainer()

	reg.Register(component.CategoryClient, func(ctx context.Context, h engine.Handle) (any, error) {
		return &mockComponent{Name: "client-inst"}, nil
	}, engine.WithScopes(component.ClientScope))

	reg.Register(component.CategoryServer, func(ctx context.Context, h engine.Handle) (any, error) {
		return &mockComponent{Name: "server-inst"}, nil
	}, engine.WithScopes(component.ServerScope))

	_ = reg.Load(ctx, "source")

	// Check cross-scope access (should fail)
	_, err := reg.In(component.CategoryClient, engine.WithInScope(component.ServerScope)).Get(ctx, "")
	if err == nil {
		t.Errorf("Should not be able to get Client component from Server scope")
	}

	// Correct access
	inst, _ := reg.In(component.CategoryClient, engine.WithInScope(component.ClientScope)).Get(ctx, "")
	if inst.(*mockComponent).Name != "client-inst" {
		t.Errorf("Failed to get component from correct scope")
	}
}
