package service

import (
	"context"
	"strings"
	"testing"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	projectContext "github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/interfaces"
	tkerrors "github.com/origadmin/toolkits/errors"
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

func (m *MockProtocolFactory) NewServer(cfg *configv1.Service, opts ...Option) (interfaces.Server, error) {
	if m.NewServerError != nil {
		return nil, m.NewServerError
	}
	return &MockServer{}, nil
}

func (m *MockProtocolFactory) NewClient(ctx projectContext.Context, cfg *configv1.Service, opts ...Option) (interfaces.Client, error) {
	if m.NewClientError != nil {
		return nil, m.NewClientError
	}
	return &MockClient{}, nil
}

// Helper to reset the registry for isolated tests
func resetProtocolRegistry() {
	defaultRegistry.Reset()
}

func TestRegisterAndGetProtocol(t *testing.T) {
	resetProtocolRegistry()

	aFactory := &MockProtocolFactory{}
	RegisterProtocol("mock_protocol", aFactory)

	factory, ok := GetProtocolFactory("mock_protocol")
	if !ok {
		t.Errorf("Expected protocol 'mock_protocol' to be registered, but it was not found.")
	}

	if factory != aFactory {
		t.Errorf("Expected retrieved factory to be the same as the registered one, got %v, want %v", factory, aFactory)
	}

	// Test for a non-existent protocol
	_, ok = GetProtocolFactory("non_existent_protocol")
	if ok {
		t.Errorf("Expected protocol 'non_existent_protocol' not to be found, but it was.")
	}
}

func TestNewServer(t *testing.T) {
	resetProtocolRegistry()

	tests := []struct {
		name            string
		cfg             *configv1.Service
		factory         *MockProtocolFactory
		expectedErr     string
		checkWrappedErr error
	}{
		{
			name:        "nil config",
			cfg:         nil,
			expectedErr: "service configuration or protocol is missing",
		},
		{
			name:        "missing protocol in config",
			cfg:         &configv1.Service{},
			expectedErr: "service configuration or protocol is missing",
		},
		{
			name:        "unsupported protocol",
			cfg:         &configv1.Service{Protocol: "unknown_protocol"},
			expectedErr: "unsupported protocol: unknown_protocol",
		},
		{
			name:            "factory returns error",
			cfg:             &configv1.Service{Protocol: "mock_error_server"},
			factory:         &MockProtocolFactory{NewServerError: tkerrors.Errorf("internal factory error")},
			expectedErr:     "failed to create server for protocol mock_error_server",
			checkWrappedErr: tkerrors.Errorf("internal factory error"),
		},
		{
			name:        "successful server creation",
			cfg:         &configv1.Service{Protocol: "mock_success_server"},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetProtocolRegistry()
			if tt.factory != nil {
				RegisterProtocol(tt.cfg.Protocol, tt.factory)
			}

			server, err := NewServer(tt.cfg)

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.expectedErr)
				} else if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.expectedErr, err)
				}
				if tt.checkWrappedErr != nil {
					if !tkerrors.Is(err, tt.checkWrappedErr) {
						t.Errorf("Expected wrapped error to be '%v', got '%v'", tt.checkWrappedErr, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if server == nil {
					t.Errorf("Expected a server instance, got nil")
				}
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	resetProtocolRegistry()

	tests := []struct {
		name            string
		cfg             *configv1.Service
		factory         *MockProtocolFactory
		expectedErr     string
		checkWrappedErr error
	}{
		{
			name:        "nil config",
			cfg:         nil,
			expectedErr: "service configuration or protocol is missing",
		},
		{
			name:        "missing protocol in config",
			cfg:         &configv1.Service{},
			expectedErr: "service configuration or protocol is missing",
		},
		{
			name:        "unsupported protocol",
			cfg:         &configv1.Service{Protocol: "unknown_protocol"},
			expectedErr: "unsupported protocol: unknown_protocol",
		},
		{
			name:            "factory returns error",
			cfg:             &configv1.Service{Protocol: "mock_error_client"},
			factory:         &MockProtocolFactory{NewClientError: tkerrors.Errorf("internal factory client error")},
			expectedErr:     "failed to create client for protocol mock_error_client",
			checkWrappedErr: tkerrors.Errorf("internal factory client error"),
		},
		{
			name:        "successful client creation",
			cfg:         &configv1.Service{Protocol: "mock_success_client"},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetProtocolRegistry()
			if tt.factory != nil {
				RegisterProtocol(tt.cfg.Protocol, tt.factory)
			}

			client, err := NewClient(projectContext.Background(), tt.cfg)

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.expectedErr)
				} else if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.expectedErr, err)
				}
				if tt.checkWrappedErr != nil {
					if !tkerrors.Is(err, tt.checkWrappedErr) {
						t.Errorf("Expected wrapped error to be '%v', got '%v'", tt.checkWrappedErr, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if client == nil {
					t.Errorf("Expected a client instance, got nil")
				}
			}
		})
	}
}
