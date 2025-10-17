package interfaces

import (
	"github.com/go-kratos/kratos/v2/middleware"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/runtime/middleware/v1"
	"github.com/origadmin/runtime/interfaces/factory"
	"github.com/origadmin/runtime/interfaces/options"
)

// KMiddleware is an alias for the Kratos middleware type.
type KMiddleware = middleware.Middleware

type (
	// MiddlewareBuilder is an interface that defines a method for registering a buildImpl.
	MiddlewareBuilder interface {
		factory.Registry[MiddlewareFactory]
		BuildClientMiddlewares(*middlewarev1.Middlewares, ...options.Option) []KMiddleware
		BuildServerMiddlewares(*middlewarev1.Middlewares, ...options.Option) []KMiddleware
	}

	// MiddlewareFactory is an interface that defines a method for creating a new buildImpl.
	// It receives the middleware-specific Protobuf configuration and the generic options.Option slice.
	// Each factory is responsible for parsing the options it cares about (e.g., by using log.FromOptions).
	MiddlewareFactory interface {
		// NewMiddlewareClient builds a client-side middleware.
		NewMiddlewareClient(*middlewarev1.Middleware, ...options.Option) (KMiddleware, bool)
		// NewMiddlewareServer builds a server-side middleware.
		NewMiddlewareServer(*middlewarev1.Middleware, ...options.Option) (KMiddleware, bool)
	}
)
