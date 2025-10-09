package grpc

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
	mw "github.com/origadmin/runtime/middleware"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewClient creates a new gRPC client.
// It is the recommended way to create a client when the protocol is known in advance.
func NewClient(ctx context.Context, cfg *transportv1.GrpcClientConfig, opts ...options.Option) (interfaces.Client, error) {
	if cfg == nil {
		return nil, tkerrors.Errorf("gRPC client config is required for creation")
	}

	// 1. Process options to extract client-specific settings (endpoint, selector filter).
	var sOpts options.Context

	// --- Client creation logic below uses the extracted, concrete 'cfg' and 'sOpts' ---

	var dialOpts []grpc.DialOption
	var mws []middleware.Middleware

	// Build client interceptors (middlewares)
	for _, name := range cfg.Middlewares {
		m, ok := mw.Get(name)
		if !ok {
			return nil, fmt.Errorf("client middleware '%s' not found in registry", name)
		}
		mws = append(mws, m)
	}
	if len(mws) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(transgrpc.ClientInterceptor(mws...)))
		// TODO: Add stream interceptors if needed
	}

	// Apply other client options from config
	if cfg.MaxRecvMsgSize > 0 {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxRecvMsgSize(int(cfg.MaxRecvMsgSize))))
	}
	if cfg.MaxSendMsgSize > 0 {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.MaxSendMsgSize(int(cfg.MaxSendMsgSize))))
	}

	// Apply TLS configuration
	if cfg.GetTls() != nil && cfg.GetTls().GetEnabled() {
		tlsConfig, err := tls.NewClientTLSConfig(cfg.GetTls())
		if err != nil {
			return nil, tkerrors.Wrapf(err, "invalid TLS config for client creation")
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(tlsConfig))
	}

	// Determine target endpoint: prioritize endpoint from options (discovery) over direct target
	target := cfg.Target
	if sOpts.Value().clientEndpoint != "" {
		target = sOpts.Value().clientEndpoint
	}

	if target == "" {
		return nil, tkerrors.Errorf("client target endpoint is required for creation")
	}

	// Apply selector filter if provided via options
	if sOpts.Value().clientSelectorFilter != nil {
		// Kratos gRPC client needs a selector builder to use a NodeFilter.
		// This typically involves creating a custom selector.Builder or adapting existing ones.
		// If the target is a discovery endpoint (e.g., "discovery:///service-name"), Kratos's default
		// discovery mechanism might integrate with a NodeFilter. Otherwise, a custom Kratos client
		// option or resolver builder would be required here.
		// For example, if using a custom Kratos selector builder:
		// selectorBuilder := selector.NewBuilderWithFilter(sOpts.Value().clientSelectorFilter)
		// dialOpts = append(dialOpts, transgrpc.WithDiscovery(discovery.NewDiscovery(target)), transgrpc.WithSelector(selectorBuilder))
		// For now, we'll assume the target is a discovery endpoint and the filter will be applied by the Kratos selector.
	}

	// Set dial timeout
	dialCtx, cancel := context.WithTimeout(ctx, service.DefaultTimeout)
	if cfg.DialTimeout != nil {
		dialCtx, cancel = context.WithTimeout(ctx, cfg.DialTimeout.AsDuration())
	}
	defer cancel()

	// Create the gRPC client connection
	conn, err := grpc.DialContext(dialCtx, target, dialOpts...)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to dial gRPC client to %s", target)
	}
	transgrpc.Dial(ctx)
	// Return the client connection (which implements interfaces.Client if type aliased correctly)
	return conn, nil
}
