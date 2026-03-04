package runtime

import (
	"testing"
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

type MockRoot struct {
	Mocks *MockContainer
}

func (m *MockRoot) GetMocks() *MockContainer { return m.Mocks }

func TestDefaultGlobalResolver_DimensionReduction(t *testing.T) {
	// We use "mocks" category which falls through to resolveGeneric -> resolveFromSource

	// 1. Test Parent Wrapper (MockRoot -> MockContainer)
	root := &MockRoot{
		Mocks: &MockContainer{
			Configs: []*MockConfig{{Name: "A"}, {Name: "B"}},
		},
	}

	mc, err := DefaultGlobalResolver(root, "mocks")
	if err != nil {
		t.Fatalf("Resolver failed: %v", err)
	}
	if len(mc.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(mc.Entries))
	}

	// 2. Test Direct Container (MockContainer) with Default
	container := &MockContainer{
		Default: &MockConfig{Name: "Def"},
	}
	mc, err = DefaultGlobalResolver(container, "mocks")
	if err != nil {
		t.Fatal(err)
	}
	// Expect "default" (from Default field) AND "Def" (from Default's Name)
	if len(mc.Entries) < 2 {
		t.Errorf("Expected at least 2 entries (default + named), got %d", len(mc.Entries))
	}
	foundDefName := false
	foundDefaultKey := false
	for _, e := range mc.Entries {
		if e.Name == "Def" {
			foundDefName = true
		}
		if e.Name == "default" {
			foundDefaultKey = true
		}
	}
	if !foundDefName || !foundDefaultKey {
		t.Error("Missing expected entries for Default config")
	}

	if mc.Active != "default" {
		t.Errorf("Active should be default, got %s", mc.Active)
	}

	// 3. Test Single Item Promotion (Direct Config Item)
	single := &MockConfig{Name: "Solo"}
	mc, err = DefaultGlobalResolver(single, "mocks")
	if err != nil {
		t.Fatal(err)
	}

	// Expect "Solo" (from GetName) AND "default" (promoted)
	if len(mc.Entries) != 2 {
		t.Errorf("Expected 2 entries (self+default), got %d: %+v", len(mc.Entries), mc.Entries)
	}
	if mc.Active != "default" {
		t.Error("Single item should be promoted to active default")
	}
}
