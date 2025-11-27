package service

import (
	"context"

	transportv1 "github.com/origadmin/runtime/api/gen/go/config/transport/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	internalfactory "github.com/origadmin/runtime/internal/factory"
)

// Protocol is the name of the gRPC, HTTP, or other protocol.
const (
	ProtocolGRPC = "grpc"
	ProtocolHTTP = "http"
)

// ProtocolFactory defines the factory standard for creating a specific protocol service instance。
type ProtocolFactory interface {
	NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error)
	NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error)
}

// defaultFactory is the default, package-level instance of the protocol registry.
var defaultFactory = internalfactory.New[ProtocolFactory]()

// RegisterProtocol registers a new protocol factory with the default registry.
// This function is the public API for registration and is safe for concurrent use.
func RegisterProtocol(name string, f ProtocolFactory) {
	defaultFactory.Register(name, f)
}

// NewServer creates a new server instance based on the provided configuration and options.
// It automatically looks up the appropriate protocol factory from the default registry.
func NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error) {
	if cfg == nil {
		return nil, runtimeerrors.NewStructured(Module, "server configuration is nil").WithCaller()
	}
	if cfg.Protocol == "" {
		return nil, runtimeerrors.NewStructured(Module, "protocol is not specified in server configuration")
	}
	protocolName := cfg.Protocol

	f, ok := defaultFactory.Get(protocolName)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "unsupported protocol: %s", protocolName).WithCaller()
	}

	server, err := f.NewServer(cfg, opts...)
	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to create server for protocol %s", protocolName).WithCaller()
	}

	o := fromOptions(opts)
	// Register the user's business logic services if a registrar is provided.
	for _, registrar := range o.registrar {
		if err := registrar.Register(o.ctx, server); err != nil {
			return nil, err
		}
	}

	return server, nil
}

// NewClient creates a new client instance based on the provided configuration and options.
// It automatically looks up the appropriate protocol factory from the default registry.
func NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error) {
	if cfg == nil {
		return nil, runtimeerrors.NewStructured(Module, "client configuration is nil").WithCaller()
	}
	if cfg.Protocol == "" {
		return nil, runtimeerrors.NewStructured(Module, "protocol is not specified in client configuration").WithCaller()
	}
	protocolName := cfg.Protocol

	f, ok := defaultFactory.Get(protocolName)
	if !ok {
		return nil, runtimeerrors.NewStructured(Module, "unsupported protocol: %s", protocolName).WithCaller()
	}

	client, err := f.NewClient(ctx, cfg, opts...)
	if err != nil {
		return nil, runtimeerrors.WrapStructured(err, Module, "failed to create client for protocol %s", protocolName).WithCaller()
	}

	return client, nil
}
