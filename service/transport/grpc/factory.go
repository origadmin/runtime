package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kgprc "github.com/go-kratos/kratos/v2/transport/grpc"
	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
	mw "github.com/origadmin/runtime/middleware"
	"google.golang.org/grpc"
)

// grpcProtocolFactory implements the service.ProtocolFactory for gRPC.
type grpcProtocolFactory struct{}

// init registers this factory with the framework's protocol registry.
func init() {
	service.RegisterProtocol("grpc", &grpcProtocolFactory{})
}

// NewServer creates a new gRPC server instance.
// It conforms to the updated ProtocolFactory interface.
func (f *grpcProtocolFactory) NewServer(cfg *transportv1.Server, opts ...service.Option) (interfaces.Server, error) {
	// 1. Extract the specific gRPC server config from the container.
	grpcConfig := cfg.GetGrpc()
	if grpcConfig == nil {
		return nil, fmt.Errorf("gRPC server config is missing in transport container")
	}

	// 2. Process options to extract registrar.
	var sOpts service.Options
	sOpts.Apply(opts...)

	grpcRegistrar, ok := sOpts.Value().registrar.(service.GRPCRegistrar)
	if !ok && sOpts.Value().registrar != nil {
		return nil, fmt.Errorf("invalid registrar: expected service.GRPCRegistrar, got %T", sOpts.Value().registrar)
	}

	// --- All creation logic below uses the extracted, concrete 'grpcConfig' ---

	var kOpts []kgprc.ServerOption
	var mws []middleware.Middleware

	// Build middleware chain
	for _, name := range grpcConfig.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		kOpts = append(kOpts, kgprc.Middleware(mws...))
	}

	// Apply other server options
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
		grpcRegistrar.RegisterGRPC(context.Background(), srv)
	}

	return srv, nil
}

// NewClient creates a new gRPC client instance.
// It conforms to the updated ProtocolFactory interface.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *transportv1.Client, opts ...service.Option) (interfaces.Client, error) {
	// 1. Extract the specific gRPC client config from the container.
	grpcConfig := cfg.GetGrpc()
	if grpcConfig == nil {
		return nil, fmt.Errorf("gRPC client config is missing in transport container")
	}

	// 2. Process options to extract client-specific settings (endpoint, selector filter).
	var sOpts service.Options
	sOpts.Apply(opts...)

	// --- Client creation logic below uses the extracted, concrete 'grpcConfig' and 'sOpts' ---

	var dialOpts []grpc.DialOption
	var mws []middleware.Middleware

	// Build client interceptors (middlewares)
	for _, name := range grpcConfig.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("client middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(kgprc.ClientInterceptor(mws...)))
		// TODO: Add stream interceptors if needed
	}

	// Apply other client options
	if grpcConfig.MaxRecvMsgSize > 0 {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxRecvMsgSize(int(grpcConfig.MaxRecvMsgSize))))
	}
	if grpcConfig.MaxSendMsgSize > 0 {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxSendMsgSize(int(grpcConfig.MaxSendMsgSize))))
	}

	// Determine target endpoint: prioritize endpoint from options (discovery) over direct target
	target := grpcConfig.Target
	if sOpts.Value().clientEndpoint != "" {
		target = sOpts.Value().clientEndpoint
	}

	// Apply selector filter if provided via options
	if sOpts.Value().clientSelectorFilter != nil {
		// Kratos gRPC client needs a selector builder to use a NodeFilter.
		// This typically involves creating a custom selector.Builder or adapting existing ones.
		// For now, we'll assume Kratos's default discovery mechanism can integrate with a NodeFilter
		// if the target is a discovery endpoint (e.g., "discovery:///service-name").
		// If explicit NodeFilter application is needed, a custom Kratos client option or resolver builder
		// would be required here. For example:
		// dialOpts = append(dialOpts, kgprc.WithNodeFilter(sOpts.Value().clientSelectorFilter)) // Hypothetical Kratos option
		// Or, if using a custom Kratos selector builder:
		// selectorBuilder := selector.NewBuilderWithFilter(sOpts.Value().clientSelectorFilter)
		// dialOpts = append(dialOpts, kgprc.WithDiscovery(discovery.NewDiscovery(target)), kgprc.WithSelector(selectorBuilder))
	}


	// Set dial timeout
	dialCtx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	if grpcConfig.DialTimeout != nil {
		dialCtx, cancel = context.WithTimeout(ctx, grpcConfig.DialTimeout.AsDuration())
	}
	defer cancel()

	// Create the gRPC client connection
	conn, err := grpc.DialContext(dialCtx, target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC client to %s: %w", target, err)
	}

	// Return the client connection (which implements interfaces.Client if type aliased correctly)
	return conn, nil
}
