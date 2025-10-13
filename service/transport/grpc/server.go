package grpc

import (
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
)

// NewGRPCServer creates a new concrete gRPC server instance based on the provided configuration.
// It returns *transgrpc.Server, not the generic interfaces.Server.
func NewGRPCServer(grpcConfig *transportv1.GrpcServerConfig, serverOpts *ServerOptions) (*transgrpc.Server, error) {
	// Initialize the Kratos gRPC server options using the adapter function.
	kratosOpts, err := initGrpcServerOptions(grpcConfig, serverOpts)
	if err != nil {
		return nil, err
	}

	// Create the Kratos gRPC server instance.
	srv := transgrpc.NewServer(kratosOpts...)

	return srv, nil
}
