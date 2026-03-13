package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/container"
)

// TestPriority_ConfigResolver verifies:
// LoadOption > RegisterOption > GlobalWithCategoryResolvers
func TestPriority_ConfigResolver(t *testing.T) {
	ctx := context.Background()
	cat := component.Category("test-resolver")

	globalResolver := func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: "global", Value: "global"}}}, nil
	}
	registerResolver := func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: "register", Value: "register"}}}, nil
	}
	loadResolver := func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: "load", Value: "load"}}}, nil
	}

	// 1. Test Global
	c1 := container.NewContainer(container.WithCategoryResolvers(map[component.Category]component.ConfigResolver{cat: globalResolver}))
	c1.Register(cat, simpleProvider)
	_ = c1.Load(ctx, "src")
	if _, err := c1.In(cat).Get(ctx, "global"); err != nil {
		t.Errorf("Global resolver should be used when others are missing")
	}

	// 2. Test Register overrides Global
	c2 := container.NewContainer(container.WithCategoryResolvers(map[component.Category]component.ConfigResolver{cat: globalResolver}))
	c2.Register(cat, simpleProvider, engine.WithConfigResolverOption(registerResolver))
	_ = c2.Load(ctx, "src")
	if _, err := c2.In(cat).Get(ctx, "register"); err != nil {
		t.Errorf("Register resolver should override global resolver")
	}

	// 3. Test Load overrides Register
	c3 := container.NewContainer(container.WithCategoryResolvers(map[component.Category]component.ConfigResolver{cat: globalResolver}))
	c3.Register(cat, simpleProvider, engine.WithConfigResolverOption(registerResolver))
	_ = c3.Load(ctx, "src", engine.WithLoadResolver(loadResolver))
	if _, err := c3.In(cat).Get(ctx, "load"); err != nil {
		t.Errorf("Load resolver should override register resolver")
	}
}

// TestPriority_RequirementResolver verifies:
// ConfigEntry > ModuleConfig > RegistrationOptions > Container.Requirement
func TestPriority_RequirementResolver(t *testing.T) {
	ctx := context.Background()
	cat := component.Category("test-req")
	purpose := "dep"

	resContainer := func(ctx context.Context, h component.Handle, p string) (any, error) { return "container", nil }
	resRegister := func(ctx context.Context, h component.Handle, p string) (any, error) { return "register", nil }
	resModule := func(ctx context.Context, h component.Handle, p string) (any, error) { return "module", nil }
	resEntry := func(ctx context.Context, h component.Handle, p string) (any, error) { return "entry", nil }

	setup := func(opts ...component.RegisterOption) component.Container {
		c := engine.NewContainer()
		c.Requirement(cat, purpose, resContainer)
		c.Register(cat, func(ctx context.Context, h component.Handle) (any, error) {
			val, _ := h.Require(purpose)
			return val, nil
		}, opts...)
		return c
	}

	// 1. Container Level
	c1 := setup()
	_ = c1.Load(ctx, "src")
	val, _ := c1.In(cat).Get(ctx, "")
	if val != "container" {
		t.Errorf("Expected container level resolver, got %v", val)
	}

	// 2. Register Level overrides Container
	c2 := setup(engine.WithRequirementResolverOption(resRegister))
	_ = c2.Load(ctx, "src")
	val, _ = c2.In(cat).Get(ctx, "")
	if val != "register" {
		t.Errorf("Expected register level resolver, got %v", val)
	}

	// 3. Module Level overrides Register
	c3 := setup(engine.WithRequirementResolverOption(resRegister))
	_ = c3.Load(ctx, "src", engine.WithLoadResolver(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries:             []component.ConfigEntry{{Name: "default", Value: "v"}},
			RequirementResolver: resModule,
		}, nil
	}))
	val, _ = c3.In(cat).Get(ctx, "")
	if val != "module" {
		t.Errorf("Expected module level resolver, got %v", val)
	}

	// 4. Entry Level overrides Module
	c4 := setup(engine.WithRequirementResolverOption(resRegister))
	_ = c4.Load(ctx, "src", engine.WithLoadResolver(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{
			Entries: []component.ConfigEntry{{
				Name:                "default",
				Value:               "v",
				RequirementResolver: resEntry,
			}},
			RequirementResolver: resModule,
		}, nil
	}))
	val, _ = c4.In(cat).Get(ctx, "")
	if val != "entry" {
		t.Errorf("Expected entry level resolver, got %v", val)
	}
}

// TestPriority_RegistrationHierarchy verifies priority numbers and registration order
func TestPriority_RegistrationHierarchy(t *testing.T) {
	ctx := context.Background()
	cat := component.Category("hierarchy")

	c := engine.NewContainer()

	// 1. Priority 100
	c.Register(cat, func(ctx context.Context, h component.Handle) (any, error) {
		return "priority-100", nil
	}, engine.WithPriority(100))

	// 2. Priority 200
	c.Register(cat, func(ctx context.Context, h component.Handle) (any, error) {
		return "priority-200", nil
	}, engine.WithPriority(200))

	// 3. Another Priority 200 (should override previous 200)
	c.Register(cat, func(ctx context.Context, h component.Handle) (any, error) {
		return "priority-200-new", nil
	}, engine.WithPriority(200))

	_ = c.Load(ctx, "src")

	val, _ := c.In(cat).Get(ctx, "")
	if val != "priority-200-new" {
		t.Errorf("Expected priority-200-new to take precedence, got %v", val)
	}
}

// TestPriority_InOptionVerifies verifies InOption chain effect
func TestPriority_InOption(t *testing.T) {
	ctx := context.Background()
	cat := component.Category("in-option")
	c := engine.NewContainer()
	c.Register(cat, simpleProvider)
	_ = c.Load(ctx, "src", engine.WithLoadResolver(func(ctx context.Context, source any, opts *component.LoadOptions) (*component.ModuleConfig, error) {
		return &component.ModuleConfig{Entries: []component.ConfigEntry{{Name: "A", Value: "v1"}, {Name: "B", Value: "v2"}}}, nil
	}))

	reg := c.In(cat, func(r component.Registry) component.Registry {
		return r.Skip("A").(component.Registry)
	})

	_, errA := reg.Get(ctx, "A")
	if errA == nil {
		t.Error("Expected error when getting skipped component A")
	}

	_, errB := reg.Get(ctx, "B")
	if errB != nil {
		t.Errorf("Should still be able to get component B: %v", errB)
	}
}
