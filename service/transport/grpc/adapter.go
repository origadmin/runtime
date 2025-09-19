package grpc

import (
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
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
func adaptServerConfig(cfg *transportv1.GRPC) ([]transgrpc.ServerOption, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("server config is required for creation")
	}

	// Start with default middlewares
	opts := DefaultServerMiddlewares()

	// Add TLS configuration if needed
	if cfg.GetTls() != nil && cfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewServerTLSConfig(cfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation")
		}
		opts = append(opts, transgrpc.TLSConfig(tlsConfig))
	}

	// Add network and address configurations
	if cfg.GetNetwork() != "" {
		opts = append(opts, transgrpc.Network(cfg.GetNetwork()))
	}

	if cfg.GetAddr() != "" {
		opts = append(opts, transgrpc.Address(cfg.GetAddr()))
	}

	// Configure timeout
	timeout := 5 * time.Second
	if cfg.GetTimeout() != 0 {
		timeout = time.Duration(cfg.GetTimeout() * 1e6)
	}
	opts = append(opts, transgrpc.Timeout(timeout))

	return opts, nil
}

// adaptClientConfig converts service configuration to Kratos gRPC client options.
func adaptClientConfig(cfg *transportv1.GRPC) ([]transgrpc.ClientOption, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation")
	}

	cfg := cfg.GetGrpc()
	if cfg == nil {
		return nil, tkerrors.Errorf("grpc client config is required for creation")
	}

	var opts []transgrpc.ClientOption

	if cfg.GetTls() != nil && cfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(cfg.GetTls())
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
