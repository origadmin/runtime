package grpc

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// adaptClientConfig converts service configuration to Kratos gRPC client options.
func adaptClientConfig(cfg *configv1.Service) ([]transgrpc.ClientOption, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation")
	}

	grpcCfg := cfg.GetGrpc()
	if grpcCfg == nil {
		return nil, tkerrors.Errorf("grpc client config is required for creation")
	}

	var kratosClientOpts []transgrpc.ClientOption

	if tlsCfg := grpcCfg.GetTlsConfig(); tlsCfg != nil {
		tlsConfig, err := tls.NewClientTLSConfig(tlsCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		if tlsConfig != nil {
			kratosClientOpts = append(kratosClientOpts, transgrpc.WithTLSConfig(tlsConfig))
		}
	}

	if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
		nodeFilter, err := selector.NewFilter(selectorCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation")
		}
		kratosClientOpts = append(kratosClientOpts, transgrpc.WithNodeFilter(nodeFilter))
	}

	return kratosClientOpts, nil
}
