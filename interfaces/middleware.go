package interfaces

import "github.com/go-kratos/kratos/v2/middleware"

// MiddlewareProvider defines the contract for retrieving a middleware instance by name.
// This allows decoupling the transport factories from a global middleware registry.
type MiddlewareProvider interface {
	GetMiddleware(name string) (middleware.Middleware, bool)
}
