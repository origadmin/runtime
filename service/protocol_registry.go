package service

import (
	"context"
	"fmt"
	"net/http"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
)

// Protocol is the name of the gRPC, HTTP, or other protocol.
const (
	ProtocolGRPC = "grpc"
	ProtocolHTTP = "http"
)

// defaultRegistry is the default, package-level instance of the protocol registry.
var defaultRegistry = factory.New[ProtocolFactory]()

// RegisterProtocol registers a new protocol factory with the default registry.
// This function is the public API for registration and is safe for concurrent use.
func RegisterProtocol(name string, f ProtocolFactory) {
	defaultRegistry.Register(name, f)
}

// GetProtocolFactory retrieves a registered protocol factory by name from the default registry.
// This function is the public API for retrieval.
func GetProtocolFactory(name string) (ProtocolFactory, bool) {
	return defaultRegistry.Get(name)
}

// getServerProtocolName extracts the protocol name from the transportv1.Server configuration.
func getServerProtocolName(cfg *transportv1.Server) (string, error) {
	if cfg == nil {
		ke := kerrors.New(http.StatusBadRequest, "ERR_NIL_SERVER_CONFIG", "server configuration is nil")
		return "", ke
	}
	if cfg.Protocol == "" {
		ke := kerrors.New(http.StatusBadRequest, "ERR_MISSING_SERVER_CONFIG", "protocol is not specified in server configuration")
		return "", ke
	}
	return cfg.Protocol, nil
}

// getClientProtocolName extracts the protocol name from the transportv1.Client configuration.
func getClientProtocolName(cfg *transportv1.Client) (string, error) {
	if cfg == nil {
		ke := kerrors.New(http.StatusBadRequest, "ERR_NIL_CLIENT_CONFIG", "client configuration is nil")
		return "", ke
	}
	if cfg.Protocol == "" {
		ke := kerrors.New(http.StatusBadRequest, "ERR_MISSING_CLIENT_CONFIG", "protocol is not specified in client configuration")
		return "", ke
	}
	return cfg.Protocol, nil
}

// NewServer creates a new server instance based on the provided configuration and options.
// It automatically looks up the appropriate protocol factory from the default registry.
func NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error) {
	protocolName, err := getServerProtocolName(cfg)
	if err != nil {
		return nil, err
	}

	f, ok := GetProtocolFactory(protocolName)
	if !ok {
		ke := kerrors.New(http.StatusBadRequest, "ERR_UNSUPPORTED_PROTOCOL", fmt.Sprintf("unsupported protocol: %s", protocolName))
		return nil, ke
	}

	server, err := f.NewServer(cfg, opts...)
	if err != nil {
		ke := kerrors.New(http.StatusInternalServerError, "ERR_SERVER_CREATION_FAILED", fmt.Sprintf("failed to create server for protocol %s", protocolName))
		ke.Metadata = map[string]string{"error": err.Error()}
		return nil, ke
	}

	return server, nil
}

// NewClient creates a new client instance based on the provided configuration and options.
func NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error) {
	protocolName, err := getClientProtocolName(cfg)
	if err != nil {
		return nil, err
	}

	f, ok := GetProtocolFactory(protocolName)
	if !ok {
		ke := kerrors.New(http.StatusBadRequest, "ERR_UNSUPPORTED_PROTOCOL", fmt.Sprintf("unsupported protocol: %s", protocolName))
		return nil, ke
	}

	client, err := f.NewClient(ctx, cfg, opts...)
	if err != nil {
		ke := kerrors.New(http.StatusInternalServerError, "ERR_CLIENT_CREATION_FAILED", fmt.Sprintf("failed to create client for protocol %s", protocolName))
		ke.Metadata = map[string]string{"error": err.Error()}
		return nil, ke
	}

	return client, nil
}
