package runtime

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/container"
)

// Mock structs mimicking protobuf generated code
type MockConfig struct {
	Name string
}

func (m *MockConfig) GetName() string { return m.Name }

type MockContainer struct {
	Active  string
	Default *MockConfig
	Configs []*MockConfig
}

func (m *MockContainer) GetActive() string         { return m.Active }
func (m *MockContainer) GetDefault() *MockConfig   { return m.Default }
func (m *MockContainer) GetConfigs() []*MockConfig { return m.Configs }

func TestDefaultResolvers_Functionality(t *testing.T) {
	ctx := context.Background()

	// Use individual container for this test to avoid shared state
	customResolver := func(source any, cat component.Category) (*component.ModuleConfig, error) {
		if c, ok := source.(*MockContainer); ok {
			res := &component.ModuleConfig{Active: c.Active}
			for _, cfg := range c.Configs {
				res.Entries = append(res.Entries, component.ConfigEntry{Name: cfg.Name, Value: cfg})
			}
			return res, nil
		}
		return nil, nil
	}

	reg := container.NewContainer(container.WithCategoryResolvers(map[component.Category]component.Resolver{
		"mocks": customResolver,
	}))

	reg.Register("mocks", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		return "mock", nil
	})

	c := &MockContainer{
		Configs: []*MockConfig{{Name: "A"}, {Name: "B"}},
	}

	if err := reg.Load(ctx, c); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	h := reg.In("mocks")
	foundA := false
	for name := range h.Iter(ctx) {
		if name == "A" {
			foundA = true
		}
	}
	if !foundA {
		t.Error("Expected to find instance A")
	}
}

func TestDefaultResolvers_PassThrough(t *testing.T) {
	ctx := context.Background()
	regRaw := container.NewContainer(nil)
	single := &MockConfig{Name: "Solo"}
	var capturedConfig any

	regRaw.Register("raw_config", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		capturedConfig = h.Config()
		return "ok", nil
	})

	if err := regRaw.Load(ctx, single, engine.ForCategory("raw_config")); err != nil {
		t.Fatal(err)
	}

	// Trigger instantiation
	_, err := regRaw.In("raw_config").Get(ctx, component.DefaultName)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if capturedConfig != single {
		t.Errorf("Pass-through mode should return original source, got %T", capturedConfig)
	}
}

func TestContainer_LifecycleLock(t *testing.T) {
	reg := container.NewContainer(nil)
	ctx := context.Background()

	reg.Register("first", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		return "ok", nil
	})

	if err := reg.Load(ctx, "root"); err != nil {
		t.Fatal(err)
	}

	// Subsequent registration should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Register after Load should have panicked")
		}
	}()

	reg.Register("second", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
		return "fail", nil
	})
}
