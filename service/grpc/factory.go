package grpc

import (
	"time"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context" // Use project's context

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/goexts/generic/configure"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/tls"
	tkerrors "github.com/origadmin/toolkits/errors"
)

const (
	defaultTimeout = 5 * time.Second
)

// grpcProtocolFactory implements service.ProtocolFactory for gRPC.
type grpcProtocolFactory struct{}

// NewClient creates a new gRPC client instance by delegating to the direct implementation.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	return NewClient(ctx, cfg, opts...)
}

// NewServer creates a new gRPC server instance using Kratos transport.
func (f *grpcProtocolFactory) NewServer(cfg *configv1.Service, opts ...service.Option) (interfaces.Server, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/grpc"))
	ll.Debugf("Creating new gRPC server instance with config: %+v", cfg)

	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: context.Background()}}
	configure.Apply(svcOpts, opts)

	var serverOpts []transgrpc.ServerOption
	serverOpts = append(serverOpts, transgrpc.Middleware(recovery.Recovery()))

	if cfg.GetGrpc() != nil {
		grpcCfg := cfg.GetGrpc()

		if grpcCfg.GetUseTls() {
			tlsConfig, err := tls.NewServerTLSConfig(grpcCfg.GetTlsConfig())
			if err != nil {
				return nil, tkerrors.Wrapf(err, "invalid TLS config for server creation")
			}
			if tlsConfig != nil {
				serverOpts = append(serverOpts, transgrpc.TLSConfig(tlsConfig))
			}
		}

		if grpcCfg.GetNetwork() != "" {
			serverOpts = append(serverOpts, transgrpc.Network(grpcCfg.GetNetwork()))
		}

		if grpcCfg.GetAddr() != "" {
			serverOpts = append(serverOpts, transgrpc.Address(grpcCfg.GetAddr()))
		}

		timeout := defaultTimeout
		if grpcCfg.GetTimeout() != 0 {
			timeout = time.Duration(grpcCfg.GetTimeout() * 1e6)
		}
		serverOpts = append(serverOpts, transgrpc.Timeout(timeout))

		ll.Debugw("gRPC server configured", "endpoint", grpcCfg.GetEndpoint())
	}

	serverOptsFromContext := FromServerOptions(svcOpts)

	serverOpts = append(serverOpts, serverOptsFromContext...)

	return transgrpc.NewServer(serverOpts...), nil
}

func init() {
	service.RegisterProtocol("grpc", &grpcProtocolFactory{})
}
