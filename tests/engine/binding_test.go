package engine_test

import (
	"context"
	"testing"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine"
	"github.com/origadmin/runtime/engine/container"
)

// --- Mock Configuration Structs ---

type MockDBConfig struct {
	Name    string
	Dialect string
	DSN     string
}

func (m *MockDBConfig) GetName() string    { return m.Name }
func (m *MockDBConfig) GetDialect() string { return m.Dialect }

type MockData struct {
	Active    string
	Default   *MockDBConfig
	Databases []*MockDBConfig
}

func (m *MockData) GetActive() string           { return m.Active }
func (m *MockData) GetDefault() *MockDBConfig   { return m.Default }
func (m *MockData) GetConfigs() []*MockDBConfig { return m.Databases }

// --- Test ---

func TestEngine_ConfigurationBindingProtocol(t *testing.T) {
	ctx := context.Background()

	// Helper to create a pre-registered registry
	newReg := func() component.Registry {
		reg := container.NewContainer()
		reg.Register("database", func(ctx context.Context, h component.Handle, opts ...options.Option) (any, error) {
			cfg := &MockDBConfig{}
			if err := h.BindConfig(cfg); err != nil {
				return nil, err
			}
			return cfg, nil
		})
		return reg
	}

	t.Run("DimensionReduction_SingleItem", func(t *testing.T) {
		reg := newReg()
		// Single item injection - should be promoted to default
		single := &MockDBConfig{Dialect: "mysql", DSN: "mysql://solo"}

		if err := reg.Load(ctx, single); err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		// Verify by dialect name
		db1, err := engine.Cast[*MockDBConfig](ctx, reg.In("database"), "mysql")
		if err != nil {
			t.Fatal(err)
		}
		if db1.DSN != "mysql://solo" {
			t.Errorf("Wrong DSN: %s", db1.DSN)
		}

		// Verify by default alias
		db2, err := engine.GetDefault[*MockDBConfig](ctx, reg.In("database"))
		if err != nil {
			t.Fatal(err)
		}
		if db2.DSN != "mysql://solo" {
			t.Errorf("Wrong DSN: %s", db2.DSN)
		}
	})

	t.Run("Protocol_DefaultAndNamed", func(t *testing.T) {
		reg := newReg()
		// Container with explicit Default
		data := &MockData{
			Default: &MockDBConfig{Name: "primary", DSN: "mysql://primary"},
			Databases: []*MockDBConfig{
				{Name: "secondary", DSN: "mysql://secondary"},
			},
		}

		if err := reg.Load(ctx, data); err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		// Verify "primary"
		dbP, err := engine.Cast[*MockDBConfig](ctx, reg.In("database"), "primary")
		if err != nil {
			t.Fatal(err)
		}

		// Verify "default" maps to "primary"
		dbD, err := engine.GetDefault[*MockDBConfig](ctx, reg.In("database"))
		if err != nil {
			t.Fatal(err)
		}
		if dbD != dbP {
			t.Errorf("default (%p) should point to same instance as primary (%p)", dbD, dbP)
		}

		// Verify "secondary"
		dbS, err := engine.Cast[*MockDBConfig](ctx, reg.In("database"), "secondary")
		if err != nil {
			t.Fatal(err)
		}
		if dbS.DSN != "mysql://secondary" {
			t.Error("secondary DSN mismatch")
		}
	})
}
