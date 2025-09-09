package grpc

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc" // Import Kratos gRPC transport

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/service/selector" // Import custom selector
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// adaptClientConfig converts service configuration to Kratos gRPC client options.
func adaptClientConfig(cfg *configv1.Service) ([]transgrpc.ClientOption, error) { // Change return type
	// 1. 验证配置
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation") // 使用内部错误
	}

	grpcCfg := cfg.GetGrpc()
	if grpcCfg == nil {
		return nil, tkerrors.Errorf("grpc client config is required for creation") // 使用内部错误
	}

	var kratosClientOpts []transgrpc.ClientOption

	// 4. 处理TLS
	if tlsCfg := grpcCfg.GetTlsConfig(); tlsCfg != nil {
		tlsConfig, err := tls.NewClientTLSConfig(tlsCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation") // 使用内部错误包装
		}
		if tlsConfig != nil {
			kratosClientOpts = append(kratosClientOpts, transgrpc.WithTLSConfig(tlsConfig)) // 使用直接的 Kratos TLS 选项
		}
	}

	// Removed: Add insecure by default if TLS is not used and not explicitly set
	// This logic will be moved to NewClient to choose between Dial and DialInsecure

	// 5. 处理选择器 (使用 transgrpc.WithSelector)
	if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
		nodeFilter, err := selector.NewFilter(selectorCfg) // Returns selector.NodeFilter (which is filter.Filter)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation") // 使用内部错误包装
		}
		// Create a Kratos selector from the custom filter
		// kratosSelector.NewSelector expects ...filter.Filter
		// Since nodeFilter is already a filter.Filter, we can pass it directly
		kratosClientOpts = append(kratosClientOpts, transgrpc.WithNodeFilter(nodeFilter))
	}

	return kratosClientOpts, nil
}
