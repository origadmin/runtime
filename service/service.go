// Package service implements the functions, types, and interfaces for the module.
package service

import (
	"github.com/origadmin/runtime/context"
)

// ServerRegistrar defines the single, universal entry point for service registration.
// It is the responsibility of the transport-specific factories to pass the correct server type
// (e.g., *GRPCServer or *HTTPServer) to the Register method.
// The implementation of this interface, typically within the user's project, is expected
// to perform a type switch to handle the specific server type.
type ServerRegistrar interface {
	Register(ctx context.Context, srv any) error
}

// ServerRegisterFunc is a function that implements ServerRegistrar.
type ServerRegisterFunc func(ctx context.Context, srv any) error

// Register implements the ServerRegistrar interface for ServerRegisterFunc.
func (f ServerRegisterFunc) Register(ctx context.Context, srv any) error {
	return f(ctx, srv)
}
