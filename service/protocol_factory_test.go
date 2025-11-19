package service

import (
	"context"
	"errors"
	"testing"

	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	projectContext "github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// MockServer implements interfaces.Server for testing purposes.
type MockServer struct{}

func (m *MockServer) Start(ctx context.Context) error { return nil }
func (m *MockServer) Stop(ctx context.Context) error  { return nil }

// MockClient implements interfaces.Client for testing purposes.
type MockClient struct{}

// MockProtocolFactory implements interfaces.ProtocolFactory for testing purposes.
type MockProtocolFactory struct {
	NewServerError error
	NewClientError error
}

// NewServer method of MockProtocolFactory
func (m *MockProtocolFactory) NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error) {
	if m.NewServerError != nil {
		return nil, m.NewServerError
	}
	return &MockServer{}, nil
}

// NewClient method of MockProtocolFactory
func (m *MockProtocolFactory) NewClient(ctx projectContext.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error) {
	if m.NewClientError != nil {
		return nil, m.NewClientError
	}
	return &MockClient{}, nil
}

// Helper to reset the registry for isolated tests
func resetProtocolRegistry() {
	defaultFactory.Reset()
}

func TestRegisterAndGetProtocol(t *testing.T) {
	resetProtocolRegistry()

	aFactory := &MockProtocolFactory{}
	RegisterProtocol("mock_protocol", aFactory)

	factory, ok := defaultFactory.Get("mock_protocol")
	if !ok {
		t.Errorf("Expected protocol 'mock_protocol' to be registered, but it was not found.")
	}

	if factory != aFactory {
		t.Errorf("Expected retrieved factory to be the same as the registered one, got %v, want %v", factory, aFactory)
	}

	// Test for a non-existent protocol
	_, ok = defaultFactory.Get("non_existent_protocol")
	if ok {
		t.Errorf("Expected protocol 'non_existent_protocol' not to be found, but it was.")
	}
}

func TestNewServer(t *testing.T) {
	resetProtocolRegistry()

	factoryErr := errors.New("internal factory error")

	tests := []struct {
		name                string
		cfg                 *transportv1.Server
		factory             *MockProtocolFactory
		expectedMsgContains string
	}{
		{
			name:                "nil config",
			cfg:                 nil,
			expectedMsgContains: "server configuration is nil", // Updated to match actual error message
		},
		{
			name:                "missing protocol in config",
			cfg:                 &transportv1.Server{Protocol: ""},
			expectedMsgContains: "protocol is not specified in server configuration",
		},
		{
			name:                "unsupported protocol (no factory registered)",
			cfg:                 &transportv1.Server{Protocol: "grpc"},
			expectedMsgContains: "unsupported protocol: grpc",
		},
		{
			name:                "factory returns error",
			cfg:                 &transportv1.Server{Protocol: "grpc"},
			factory:             &MockProtocolFactory{NewServerError: factoryErr},
			expectedMsgContains: "failed to create server for protocol grpc",
		},
		{
			name:                "successful server creation",
			cfg:                 &transportv1.Server{Protocol: "grpc"},
			factory:             &MockProtocolFactory{},
			expectedMsgContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetProtocolRegistry()

			// Register the factory before creating the server
			if tt.factory != nil && tt.cfg != nil && tt.cfg.Protocol != "" {
				RegisterProtocol(tt.cfg.Protocol, tt.factory)
			}

			server, err := NewServer(tt.cfg)

			if tt.expectedMsgContains != "" {
				// For error cases
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.expectedMsgContains)
				}
				if !contains(err.Error(), tt.expectedMsgContains) {
					t.Errorf("Expected error to contain %q, got %v", tt.expectedMsgContains, err)
				}
			} else {
				// For success cases
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if server == nil {
					t.Error("Expected a server instance, got nil")
				}
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	resetProtocolRegistry()

	factoryErr := errors.New("internal factory client error")

	tests := []struct {
		name                string
		cfg                 *transportv1.Client
		factory             *MockProtocolFactory
		expectedMsgContains string
	}{
		{
			name:                "nil config",
			cfg:                 nil,
			expectedMsgContains: "client configuration is nil", // Updated to match actual error message
		},
		{
			name: "missing protocol in config",
			cfg: &transportv1.Client{
				// Protocol field is required
			},
			expectedMsgContains: "protocol is not specified in client configuration",
		},
		{
			name:                "unsupported protocol (no factory registered)",
			cfg:                 &transportv1.Client{Protocol: "test"},
			expectedMsgContains: "unsupported protocol: test",
		},
		{
			name:                "factory returns error",
			cfg:                 &transportv1.Client{Protocol: "test"},
			factory:             &MockProtocolFactory{NewClientError: factoryErr},
			expectedMsgContains: "failed to create client for protocol test",
		},
		{
			name:                "successful client creation",
			cfg:                 &transportv1.Client{Protocol: "test"},
			factory:             &MockProtocolFactory{},
			expectedMsgContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetProtocolRegistry()

			// Register the factory before creating the client
			if tt.factory != nil && tt.cfg != nil && tt.cfg.Protocol != "" {
				RegisterProtocol(tt.cfg.Protocol, tt.factory)
			}

			client, err := NewClient(context.Background(), tt.cfg)

			if tt.expectedMsgContains != "" {
				// For error cases
				if err == nil {
					t.Fatalf("Expected error containing %q, got nil", tt.expectedMsgContains)
				}
				if !contains(err.Error(), tt.expectedMsgContains) {
					t.Errorf("Expected error to contain %q, got %v", tt.expectedMsgContains, err)
				}
			} else {
				// For success cases
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if client == nil {
					t.Error("Expected a client instance, got nil")
				}
			}
		})
	}
}

// contains checks if sub is a substring of s (safe for empty strings)
func contains(s, sub string) bool {
	if sub == "" {
		return true
	}
	return len(s) >= len(sub) && (s == sub || (len(s) > 0 && len(sub) > 0 && indexOf(s, sub) >= 0))
}

// indexOf returns the index of sub in s, or -1 if not found.
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
