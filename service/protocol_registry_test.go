package service

import (
	"testing"

	"github.com/origadmin/framework/runtime/interfaces"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context"
)

// MockServer implements interfaces.Server for testing purposes.
type MockServer struct{}

func (m *MockServer) Start(ctx context.Context) error { return nil }
func (m *MockServer) Stop(ctx context.Context) error  { return nil }

// MockProtocolFactory implements interfaces.ProtocolFactory for testing purposes.
type MockProtocolFactory struct{}

func (m *MockProtocolFactory) NewServer(cfg *configv1.Service, opts ...Option) (interfaces.Server, error) {
	return &MockServer{}, nil
}

func TestRegisterAndGetProtocol(t *testing.T) {
	// Clear the registry before each test to ensure isolation
	protocolRegistryLock.Lock()
	protocolRegistry = make(map[string]ProtocolFactory)
	protocolRegistryLock.Unlock()

	aFactory := &MockProtocolFactory{}
	RegisterProtocol("mock_protocol", aFactory)

	factory, ok := getProtocolFactory("mock_protocol")
	if !ok {
		t.Errorf("Expected protocol 'mock_protocol' to be registered, but it was not found.")
	}

	if factory != aFactory {
		t.Errorf("Expected retrieved factory to be the same as the registered one, got %v, want %v", factory, aFactory)
	}

	// Test for a non-existent protocol
	_, ok = getProtocolFactory("non_existent_protocol")
	if ok {
		t.Errorf("Expected protocol 'non_existent_protocol' not to be found, but it was.")
	}
}
