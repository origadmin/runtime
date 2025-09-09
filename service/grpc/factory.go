package grpc

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/goexts/generic/configure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/tls"
)

const defaultTimeout = 5 * time.Second

// grpcProtocolFactory implements service.ProtocolFactory for gRPC.
type grpcProtocolFactory struct{}

// NewClient creates a new gRPC client instance using the current recommended grpc.NewClient function.
func (f *grpcProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/grpc"))
	ll.Debugf("Creating new gRPC client with config: %+v", cfg)

	svcOpts := &service.Options{ContextOptions: interfaces.ContextOptions{Context: ctx}}
	configure.Apply(svcOpts, opts)

	var dialOpts []grpc.DialOption
	var endpoint string

	if cfg.GetGrpc() != nil {
		grpcCfg := cfg.GetGrpc()
		endpoint = grpcCfg.GetEndpoint()

		var creds credentials.TransportCredentials
		if grpcCfg.GetUseTls() {
			tlsConfig, err := tls.NewClientTLSConfig(grpcCfg.GetTlsConfig())
			if err != nil {
				return nil, err
			}
			creds = credentials.NewTLS(tlsConfig)
		} else {
			creds = insecure.NewCredentials()
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}

	dialOptsFromContext := FromClientOptions(svcOpts)
	dialOpts = append(dialOpts, dialOptsFromContext...)

	conn, err := grpc.NewClient(endpoint, dialOpts...)
	if err != nil {
		return nil, errors.Newf(500, "INTERNAL_SERVER_ERROR", "create grpc client failed: %v", err)
	}

	return conn, nil
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
				return nil, err
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
