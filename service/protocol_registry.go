package service

import (
	"context"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/factory"
	tkerrors "github.com/origadmin/toolkits/errors"
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

// getProtocolName extracts the protocol name from the transportv1.Transport configuration.
func getProtocolName(cfg *transportv1.Transport) (string, error) {
	if cfg == nil {
		return "", tkerrors.Errorf("transport configuration is nil")
	}
	switch cfg.Protocol.(type) {
	case *transportv1.Transport_Grpc:
		return "grpc", nil
	case *transportv1.Transport_Http:
		return "http", nil
	default:
		return "", tkerrors.Errorf("unknown or missing protocol in transport configuration")
	}
}

// NewServer creates a new server instance based on the provided configuration and options.
// It automatically looks up the appropriate protocol factory from the default registry.
func NewServer(cfg *transportv1.Transport, opts ...Option) (interfaces.Server, error) {
	protocolName, err := getProtocolName(cfg)
	if err != nil {
		return nil, err
	}

	f, ok := GetProtocolFactory(protocolName)
	if !ok {
		return nil, tkerrors.Errorf("unsupported protocol: %s", protocolName)
	}

	server, err := f.NewServer(cfg, opts...)
	if err != nil {
		return nil, tkerrors.Errorf("failed to create server for protocol %s: %w", protocolName, err)
	}

	return server, nil
}

// NewClient creates a new client instance based on the provided configuration and options.
func NewClient(ctx context.Context, cfg *transportv1.Transport, opts ...Option) (interfaces.Client, error) {
	protocolName, err := getProtocolName(cfg)
	if err != nil {
		return nil, err
	}

	f, ok := GetProtocolFactory(protocolName)
	if !ok {
		return nil, tkerrors.Errorf("unsupported protocol: %s", protocolName)
	}

	client, err := f.NewClient(ctx, cfg, opts...)
	if err != nil {
		return nil, tkerrors.Errorf("failed to create client for protocol %s: %w", protocolName, err)
	}

	return client, nil
}
