package service

import (
	"context"

	commonv1 "github.com/origadmin/runtime/api/gen/go/runtime/common/v1"
	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	"github.com/origadmin/runtime/errors"
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
		return "", errors.WithReason(errors.NewStructured(Module, "server configuration is nil").WithCaller(), commonv1.ErrorReason_VALIDATION_ERROR)
	}
	if cfg.Protocol == "" {
		return "", errors.WithReason(errors.NewStructured(Module, "protocol is not specified in server configuration").WithCaller(), commonv1.ErrorReason_VALIDATION_ERROR)
	}
	return cfg.Protocol, nil
}

// getClientProtocolName extracts the protocol name from the transportv1.Client configuration.
func getClientProtocolName(cfg *transportv1.Client) (string, error) {
	if cfg == nil {
		return "", errors.WithReason(errors.NewStructured(Module, "client configuration is nil").WithCaller(), commonv1.ErrorReason_VALIDATION_ERROR)
	}
	if cfg.Protocol == "" {
		return "", errors.WithReason(errors.NewStructured(Module, "protocol is not specified in client configuration").WithCaller(), commonv1.ErrorReason_VALIDATION_ERROR)
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
		return nil, errors.WithReason(errors.NewStructured(Module, "unsupported protocol: %s", protocolName).WithCaller(), commonv1.ErrorReason_VALIDATION_ERROR)
	}

	server, err := f.NewServer(cfg, opts...)
	if err != nil {
		return nil, errors.WithReason(errors.NewStructured(Module, "failed to create server for protocol %s", protocolName).WithField("error", err.Error()).WithCaller(), commonv1.ErrorReason_INTERNAL_SERVER_ERROR)
	}

	return server, nil
}

// NewClient creates a new client instance based on the provided configuration and options.
// It automatically looks up the appropriate protocol factory from the default registry.
func NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error) {
	protocolName, err := getClientProtocolName(cfg)
	if err != nil {
		return nil, err
	}

	f, ok := GetProtocolFactory(protocolName)
	if !ok {
		return nil, errors.WithReason(errors.NewStructured(Module, "unsupported protocol: %s", protocolName).WithCaller(), commonv1.ErrorReason_VALIDATION_ERROR)
	}

	client, err := f.NewClient(ctx, cfg, opts...)
	if err != nil {
		return nil, errors.WithReason(errors.NewStructured(Module, "failed to create client for protocol %s", protocolName).WithField("error", err.Error()).WithCaller(), commonv1.ErrorReason_INTERNAL_SERVER_ERROR)
	}

	return client, nil
}
