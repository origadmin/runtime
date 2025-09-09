package service

import (
	"fmt"
	tkerrors "github.com/origadmin/toolkits/errors"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/factory"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	context "github.com/origadmin/runtime/context"
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

// NewServer creates a new server instance based on the provided configuration and options.
// It automatically looks up the appropriate protocol factory from the default registry
// based on the `cfg.Type` field which represents the protocol name.
func NewServer(cfg *configv1.Service, opts ...Option) (interfaces.Server, error) {
	if cfg == nil || cfg.GetProtocol() == "" {
		return nil, tkerrors.Errorf("service configuration or protocol is missing")
	}

	f, ok := GetProtocolFactory(cfg.Protocol)
	if !ok {
		return nil, tkerrors.Errorf("unsupported protocol: %s", cfg.Protocol)
	}

	server, err := f.NewServer(cfg, opts...)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to create server for protocol %s", cfg.Protocol)
	}

	return server, nil
}

// NewClient creates a new client instance based on the provided configuration and options.
func NewClient(ctx context.Context, cfg *configv1.Service, opts ...Option) (interfaces.Client, error) {
	if cfg == nil || cfg.GetProtocol() == "" {
		return nil, tkerrors.Errorf("service configuration or protocol is missing")
	}
	f, ok := GetProtocolFactory(cfg.Protocol)
	if !ok {
		return nil, tkerrors.Errorf("unsupported protocol: %s", cfg.Protocol)
	}

	client, err := f.NewClient(ctx, cfg, opts...)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to create client for protocol %s", cfg.Protocol)
	}

	return client, nil
}
