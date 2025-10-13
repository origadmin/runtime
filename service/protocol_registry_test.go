package service

import (
	"context"
	"errors"
	"testing"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
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

	factoryErr := errors.New("internal factory error")

	tests := []struct {
		name                string
		cfg                 *transportv1.Server
		factory             *MockProtocolFactory
		expectedReason      string
		expectedMsgContains string
		checkFactoryErrInMD bool
	}{
		{
			name:                "nil config",
			cfg:                 nil,
			expectedReason:      "ERR_NIL_SERVER_CONFIG",
			expectedMsgContains: "server configuration is nil",
		},
		{
			name:                "missing protocol in config",
			cfg:                 &transportv1.Server{Protocol: ""},
			expectedReason:      "ERR_MISSING_SERVER_CONFIG",
			expectedMsgContains: "protocol is not specified",
		},
		{
			name:                "unsupported protocol (no factory registered)",
			cfg:                 &transportv1.Server{Protocol: "grpc"},
			expectedReason:      "ERR_UNSUPPORTED_PROTOCOL",
			expectedMsgContains: "unsupported protocol: grpc",
		},
		{
			name:                "factory returns error",
			cfg:                 &transportv1.Server{Protocol: "grpc"},
			factory:             &MockProtocolFactory{NewServerError: factoryErr},
			expectedReason:      "ERR_SERVER_CREATION_FAILED",
			expectedMsgContains: "failed to create server for protocol grpc",
			checkFactoryErrInMD: true,
		},
		{
			name:                "successful server creation",
			cfg:                 &transportv1.Server{Protocol: "grpc"},
			factory:             &MockProtocolFactory{},
			expectedReason:      "",
			expectedMsgContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetProtocolRegistry()
			if tt.factory != nil {
				// Extract protocol name from cfg for registration
				protocolName, err := getServerProtocolName(tt.cfg) // Use the helper
				// Ignore error here if we expect an error from getServerProtocolName
				if err == nil && protocolName != "" {
					RegisterProtocol(protocolName, tt.factory)
				}
			}

			server, err := NewServer(tt.cfg)

			if tt.expectedReason != "" {
				if err == nil {
					t.Fatalf("Expected error with reason %q, got nil", tt.expectedReason)
				}
				ke := kerrors.FromError(err)
				if ke == nil {
					t.Fatalf("Expected Kratos error, got: %v", err)
				}
				if ke.Reason != tt.expectedReason {
					t.Errorf("Expected reason %q, got %q", tt.expectedReason, ke.Reason)
				}
				if tt.expectedMsgContains != "" && (ke.Message == "" || !contains(ke.Message, tt.expectedMsgContains)) {
					t.Errorf("Expected message to contain %q, got %q", tt.expectedMsgContains, ke.Message)
				}
				if tt.checkFactoryErrInMD {
					if ke.Metadata == nil || ke.Metadata["error"] != factoryErr.Error() {
						t.Errorf("Expected metadata 'error' to be %q, got %v", factoryErr.Error(), ke.Metadata)
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

	factoryErr := errors.New("internal factory client error")

	tests := []struct {
		name                string
		cfg                 *transportv1.Client
		factory             *MockProtocolFactory
		expectedReason      string
		expectedMsgContains string
		checkFactoryErrInMD bool
	}{
		{
			name:                "nil config",
			cfg:                 nil,
			expectedReason:      "ERR_NIL_CLIENT_CONFIG",
			expectedMsgContains: "client configuration is nil",
		},
		{
			name:                "missing protocol in config",
			cfg:                 &transportv1.Client{}, // Empty protocol field
			expectedReason:      "ERR_MISSING_CLIENT_CONFIG",
			expectedMsgContains: "protocol is not specified",
		},
		{
			name:                "unsupported protocol (no factory registered)",
			cfg:                 &transportv1.Client{Protocol: "grpc"},
			expectedReason:      "ERR_UNSUPPORTED_PROTOCOL",
			expectedMsgContains: "unsupported protocol: grpc",
		},
		{
			name:                "factory returns error",
			cfg:                 &transportv1.Client{Protocol: "grpc"},
			factory:             &MockProtocolFactory{NewClientError: factoryErr},
			expectedReason:      "ERR_CLIENT_CREATION_FAILED",
			expectedMsgContains: "failed to create client for protocol grpc",
			checkFactoryErrInMD: true,
		},
		{
			name:                "successful client creation",
			cfg:                 &transportv1.Client{Protocol: "grpc"},
			factory:             &MockProtocolFactory{},
			expectedReason:      "",
			expectedMsgContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetProtocolRegistry()
			if tt.factory != nil {
				protocolName, err := getClientProtocolName(tt.cfg)
				if err == nil && protocolName != "" {
					RegisterProtocol(protocolName, tt.factory)
				}
			}

			client, err := NewClient(projectContext.Background(), tt.cfg)

			if tt.expectedReason != "" {
				if err == nil {
					t.Fatalf("Expected error with reason %q, got nil", tt.expectedReason)
				}
				ke := kerrors.FromError(err)
				if ke == nil {
					t.Fatalf("Expected Kratos error, got: %v", err)
				}
				if ke.Reason != tt.expectedReason {
					t.Errorf("Expected reason %q, got %q", tt.expectedReason, ke.Reason)
				}
				if tt.expectedMsgContains != "" && (ke.Message == "" || !contains(ke.Message, tt.expectedMsgContains)) {
					t.Errorf("Expected message to contain %q, got %q", tt.expectedMsgContains, ke.Message)
				}
				if tt.checkFactoryErrInMD {
					if ke.Metadata == nil || ke.Metadata["error"] != factoryErr.Error() {
						t.Errorf("Expected metadata 'error' to be %q, got %v", factoryErr.Error(), ke.Metadata)
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
