package grpc

import (
	"fmt" // 导入 fmt 包，用于 tkerrors.Errorf/Wrapf
	"time"

	"google.golang.org/grpc"
	tkerrors "github.com/origadmin/toolkits/errors"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/service/selector"
	"github.com/origadmin/runtime/service/tls"
)

// adaptClientConfig 将服务配置转换为特定协议的选项
func adaptClientConfig(cfg *configv1.Service) ([]grpc.ClientOption, error) {
	// 1. 验证配置
	if cfg == nil {
		return nil, tkerrors.Errorf("client config is required for creation") // 修正为 Errorf
	}

	grpcCfg := cfg.GetGrpc()
	if grpcCfg == nil {
		return nil, tkerrors.Errorf("grpc client config is required for creation") // 修正为 Errorf
	}

	var opts []grpc.ClientOption

	// 2. 处理端点 (gRPC typically uses WithEndpoint or WithTarget)
	if endpoint := grpcCfg.GetEndpoint(); endpoint != "" {
		opts = append(opts, grpc.WithEndpoint(endpoint))
	}

	// 3. 处理超时 (gRPC client timeout is usually per-call, not a global client option)
	// For now, we'll skip global timeout for gRPC client options as it's not directly analogous
	// to http.WithTimeout as a client option. It's usually handled via context.WithTimeout for calls.

	// 4. 处理TLS
	if tlsCfg := grpcCfg.GetTlsConfig(); tlsCfg != nil {
		tlsConfig, err := tls.NewClientTLSConfig(tlsCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation") // 修正为 Wrapf
		}
		if tlsConfig != nil {
			opts = append(opts, grpc.WithTransportCredentials(tls.NewClientCredentials(tlsConfig)))
		}
	}

	// 5. 处理选择器
	if selectorCfg := cfg.GetSelector(); selectorCfg != nil {
		filter, err := selector.NewFilter(selectorCfg)
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid selector config for client creation") // 修正为 Wrapf
		}
		opts = append(opts, grpc.WithDefaultServiceConfig(filter.String())) // Simplified, usually more complex for gRPC
	}

	// Add insecure by default if TLS is not used and not explicitly set
	if grpcCfg.GetUseTls() == false {
		opts = append(opts, grpc.WithInsecure())
	}

	return opts, nil
}
