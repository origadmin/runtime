package grpc

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	kgprc "github.com/go-kratos/kratos/v2/transport/grpc"
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/optionutil"
	"github.com/origadmin/runtime/service"
	"google.golang.org/grpc"
)

// grpcProtocolFactory implements the service.ProtocolFactory for gRPC.
type grpcProtocolFactory struct{}

// init registers this factory with the framework's protocol registry.
func init() {
	service.RegisterProtocol("grpc", &grpcProtocolFactory{})
}

// NewServer creates a new gRPC server instance.
func (f *grpcProtocolFactory) NewServer(cfg *transportv1.Server, opts ...options.Option) (interfaces.Server, error) {
	// 1. Extract the specific gRPC server config from the container.
	grpcConfig := cfg.GetGrpc()
	if grpcConfig == nil {
		return nil, fmt.Errorf("gRPC server config is missing in transport container")
	}

	// 2. Apply options to a Options struct to get a configured context.
	initialServiceCfg := &service.Options{}
	configuredContext := optionutil.Apply(initialServiceCfg, opts...)

	// 3. Retrieve the fully configured Options from the context.
	configuredServiceCfg, ok := optionutil.ConfigFromContext[*service.Options](configuredContext)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve configured service options from context")
	}

	// 4. Get the registrar from the configured service options.
	grpcRegistrar, ok := configuredServiceCfg.Registrar.(service.GRPCRegistrar)
	if !ok && configuredServiceCfg.Registrar != nil {
		return nil, fmt.Errorf("invalid registrar: expected service.GRPCRegistrar, got %T", configuredServiceCfg.Registrar)
	}

	// 5. Get the middleware provider.
	if configuredServiceCfg.MiddlewareProvider == nil {
		return nil, fmt.Errorf("middleware provider not found in options")
	}

	var kOpts []kgprc.ServerOption
	var mws []middleware.Middleware

	// Build middleware chain using the provider
	for _, name := range grpcConfig.Middlewares {
		m, ok := configuredServiceCfg.MiddlewareProvider.GetMiddleware(name)
		if !ok {
			return nil, fmt.Errorf("middleware '%s' not found via provider", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		kOpts = append(kOpts, kgprc.Middleware(mws...))
	}

	// Apply other server options from protobuf config
	if grpcConfig.Network != "" {
		kOpts = append(kOpts, kgprc.Network(grpcConfig.Network))
	}
	if grpcConfig.Addr != "" {
		kOpts = append(kOpts, kgprc.Address(grpcConfig.Addr))
	}
	if grpcConfig.Timeout != nil {
		kOpts = append(kOpts, kgprc.Timeout(grpcConfig.Timeout.AsDuration()))
	}
	if grpcConfig.ShutdownTimeout != nil {
		kOpts = append(kOpts, kgprc.ShutdownTimeout(grpcConfig.ShutdownTimeout.AsDuration()))
	}
	// TODO: Add TLS configuration

	// Create the gRPC server instance
	srv := kgprc.NewServer(kOpts...)

	// Register business logic
	if grpcRegistrar != nil {
		grpcRegistrar.RegisterGRPC(srv)
	}

	return srv, nil
}

// NewClient creates a new gRPC client instance.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Client, opts ...options.Option) (interfaces.Client, error) {
	// 1. Extract the specific gRPC client config from the container.
	grpcConfig := cfg.GetGrpc()
	if grpcConfig == nil {
		return nil, fmt.Errorf("gRPC client config is missing in transport container")
	}

	// 2. Apply options to get the configured context and service options.
	initialServiceCfg := &service.Options{}
	configuredContext := optionutil.Apply(initialServiceCfg, opts...)
	configuredServiceCfg, ok := optionutil.ConfigFromContext[*service.Options](configuredContext)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve configured service options from context")
	}

	// 3. Get the middleware provider.
	if configuredServiceCfg.MiddlewareProvider == nil {
		return nil, fmt.Errorf("middleware provider not found in options")
	}

	var dialOpts []grpc.DialOption
	var mws []middleware.Middleware

	// Build client interceptors (middlewares) using the provider
	for _, name := range grpcConfig.Middlewares {
		m, ok := configuredServiceCfg.MiddlewareProvider.GetMiddleware(name)
		if !ok {
			return nil, fmt.Errorf("client middleware '%s' not found via provider", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(kgprc.ClientInterceptor(mws...)))
		// TODO: Add stream interceptors if needed
	}

	// Apply other client options from protobuf config
	if grpcConfig.MaxRecvMsgSize > 0 {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxRecvMsgSize(int(grpcConfig.MaxRecvMsgSize))))
	}
	if grpcConfig.MaxSendMsgSize > 0 {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxSendMsgSize(int(grpcConfig.MaxSendMsgSize))))
	}

	// Determine target endpoint: prioritize endpoint from options over direct target from config
	target := grpcConfig.Target
	if configuredServiceCfg.ClientEndpoint != "" {
		target = configuredServiceCfg.ClientEndpoint
	}

	// Apply selector filter if provided via options
	if configuredServiceCfg.ClientSelectorFilter != nil {
		// TODO: Kratos gRPC client needs a selector builder to use a NodeFilter.
		// This part requires a more complex integration with a discovery/selector builder.
	}

	// Set dial timeout
	dialCtx := ctx
	if grpcConfig.DialTimeout != nil {
		var cancel context.CancelFunc
		dialCtx, cancel = context.WithTimeout(ctx, grpcConfig.DialTimeout.AsDuration())
		defer cancel()
	}

	// Create the gRPC client connection
	conn, err := grpc.DialContext(dialCtx, target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC client to %s: %w", target, err)
	}

	return conn, nil
}
