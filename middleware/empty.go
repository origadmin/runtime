// Package middleware implements the functions, types, and interfaces for the module.
package middleware

// Empty returns an empty middleware.
func Empty() KMiddleware {
	return func(handler KHandler) KHandler {
		return handler
	}
}
