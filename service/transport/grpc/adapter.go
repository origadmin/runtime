package grpc

import (
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// DefaultServerMiddlewares returns the default server middlewares
func DefaultServerMiddlewares() []transgrpc.ServerOption {
	return []transgrpc.ServerOption{
		transgrpc.Middleware(recovery.Recovery()),
	}
}

// adaptServerConfig converts service configuration to Kratos gRPC server options.
func adaptServerConfig(cfg *configv1.Service) ([]transgrpc.ServerOption, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("server config is required for creation")
	}

	grpcCfg := cfg.GetGrpc()
	if grpcCfg == nil {
		return nil, tkerrors.Errorf("grpc server config is required for creation")
	}

	// Start with default middlewares
	opts := DefaultServerMiddlewares()

	// Add TLS configuration if needed
	if grpcCfg.GetTls() != nil && grpcCfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewServerTLSConfig(grpcCfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation")
		}
		opts = append(opts, transgrpc.TLSConfig(tlsConfig))
	}

	// Add network and address configurations
	if grpcCfg.GetNetwork() != "" {
		opts = append(opts, transgrpc.Network(grpcCfg.GetNetwork()))
	}

	if grpcCfg.GetAddr() != "" {
		opts = append(opts, transgrpc.Address(grpcCfg.GetAddr()))
	}

	// Configure timeout
	timeout := 5 * time.Second
	if grpcCfg.GetTimeout() != 0 {
		timeout = time.Duration(grpcCfg.GetTimeout() * 1e6)
	}
	opts = append(opts, transgrpc.Timeout(timeout))

	return opts, nil
}

// adaptClientConfig converts service configuration to Kratos gRPC client options.
func adaptClientConfig(cfg *configv1.Service) ([]transgrpc.ClientOption, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation")
	}

	grpcCfg := cfg.GetGrpc()
	if grpcCfg == nil {
		return nil, tkerrors.Errorf("grpc client config is required for creation")
	}

	var opts []transgrpc.ClientOption

	if grpcCfg.GetTls() != nil && grpcCfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(grpcCfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		opts = append(opts, transgrpc.WithTLSConfig(tlsConfig))
	}

	if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
		nodeFilter, err := selector.NewFilter(selectorCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation")
		}
		opts = append(opts, transgrpc.WithNodeFilter(nodeFilter))
	}

	return opts, nil
}
