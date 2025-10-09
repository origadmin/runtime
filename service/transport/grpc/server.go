package grpc

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces" // 修正导入路径
	servicetls "github.com/origadmin/runtime/service/tls"
)

// NewGRPCServer creates a new concrete gRPC server instance based on the provided configuration.
// It returns *transgrpc.Server, not the generic interfaces.Server.
func NewGRPCServer(grpcConfig *transportv1.GrpcServerConfig, serverOpts *ServerOptions) (*transgrpc.Server, error) {
	// 3. Prepare the Kratos gRPC server options.
	var kratosOpts []transgrpc.ServerOption

	// 提前检查并获取 Container 实例
	var c interfaces.Container
	if serverOpts.ServiceOptions != nil {
		c = serverOpts.ServiceOptions.Container
	}

	// 如果配置了中间件但 Container 为 nil，则提前返回错误
	if len(grpcConfig.GetMiddlewares()) > 0 && c == nil {
		return nil, fmt.Errorf("application container is required for middleware but not found in options")
	}

	// 9. Register the user's business logic services if a registrar is provided.
	if serverOpts.ServiceOptions != nil && serverOpts.ServiceOptions.Registrar != nil {
		if c == nil {
			return nil, fmt.Errorf("application container is required for registrar but not found in options")
		}
		if grpcRegistrar, ok := serverOpts.ServiceOptions.Registrar.(service.GRPCRegistrar); ok {
			grpcRegistrar.RegisterGRPC(srv)
		} else {
			return nil, fmt.Errorf("invalid registrar: expected service.GRPCRegistrar, got %T", serverOpts.ServiceOptions.Registrar)
		}
	}

	// 4. Apply options from the protobuf configuration.
	if grpcConfig.GetAddr() != "" {
		kratosOpts = append(kratosOpts, transgrpc.Address(grpcConfig.GetAddr()))
	}
	// 5. Configure TLS for server
	if tlsConfig := grpcConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewServerTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create server TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transgrpc.TLSConfig(tlsCfg))
	}

	if grpcConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transgrpc.Timeout(grpcConfig.GetTimeout().AsDuration()))
	}

	// 6. Configure middlewares.
	var mws []middleware.Middleware
	if len(grpcConfig.GetMiddlewares()) > 0 {
		// 此时 c 已经确保不为 nil (因为上面的提前检查)
		for _, name := range grpcConfig.GetMiddlewares() {
			m, ok := c.ServerMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		// If no specific middlewares are configured, use default ones from adapter.go.
		mws = DefaultServerMiddlewares()
	}

	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transgrpc.Middleware(mws...))
	}

	// 7. Apply any external Kratos gRPC server options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(serverOpts.GrpcServerOptions) > 0 {
		kratosOpts = append(kratosOpts, serverOpts.GrpcServerOptions...)
	}

	// 8. Create the Kratos gRPC server instance.
	srv := transgrpc.NewServer(kratosOpts...)

	return srv, nil
}
