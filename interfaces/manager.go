package interfaces

import (
	"github.com/go-kratos/kratos/v2"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
	"github.com/origadmin/runtime/middleware"
)

type MiddlewareProvider interface {
	BuildClient(cfg *middlewarev1.Middleware, opts ...middleware.Option) []middleware.KMiddleware
	BuildServer(cfg *middlewarev1.Middleware, opts ...middleware.KMiddleware) []middleware.KMiddleware
}

type ServiceProvider interface {
	CreateService(cfg *configv1.Service) (*kratos.App, error)
}
