package engine_test

import (
	"context"
	"github.com/origadmin/runtime/engine"
)

// mockComponent is a generic mock for testing
type mockComponent struct {
	Name string
	Tag  string
	Type string
}

func (m *mockComponent) String() string {
	return m.Name
}

// contains checks if a string slice contains a target
func contains(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}

// simpleProvider creates a mockComponent based on the handle info
func simpleProvider(ctx context.Context, h engine.Handle) (any, error) {
	return &mockComponent{
		Name: h.Name(),
		Tag:  h.Tag(),
	}, nil
}
