package interfaces

import (
	kratosmiddleware "github.com/go-kratos/kratos/v2/middleware" // Alias for kratos middleware
	"github.com/go-kratos/kratos/v2/transport"
)

type MiddlewareProvider interface {
	BuildClient(cfg MiddlewareConfig, opts ...interface{}) []kratosmiddleware.Middleware
	BuildServer(cfg MiddlewareConfig, opts ...interface{}) []kratosmiddleware.Middleware
}

// ServerBuilder is an interface for building transport servers (HTTP, gRPC).
type ServerBuilder interface {
	DefaultBuild(ServiceConfig, ...interface{}) (transport.Server, error)
	Build(string, ServiceConfig, ...interface{}) (transport.Server, error)
}
