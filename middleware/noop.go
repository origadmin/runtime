// Package middleware implements the functions, types, and interfaces for the module.
package middleware

// Noop returns an empty middleware.
func Noop() KMiddleware {
	return func(handler KHandler) KHandler {
		return handler
	}
}
