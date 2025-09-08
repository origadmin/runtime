package service

import (
	"github.com/origadmin/framework/runtime/api/gen/go/apierrors"
	"github.com/origadmin/framework/runtime/interfaces"
	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/errors"
)

// ServerBuilder defines the generic interface for building a server.
// T is the concrete server type, which must implement interfaces.Server.
type ServerBuilder[T interfaces.Server] interface {
	Build(cfg *configv1.Service, opts ...Option) (T, error)
}

// serverBuilder is the concrete implementation of ServerBuilder.
type serverBuilder[T interfaces.Server] struct{}

// Build creates a new server instance based on the provided configuration and options.
func (sb *serverBuilder[T]) Build(cfg *configv1.Service, opts ...Option) (T, error) {
	if cfg == nil || cfg.Type == "" {
		return *new(T), errors.NewMessage(apierrors.ErrorReason_INVALID_PARAMETER, "service configuration or type is missing") // Changed INVALID_ARGUMENT to INVALID_PARAMETER
	}

	factory, ok := getProtocolFactory(cfg.Type)
	if !ok {
		return *new(T), errors.NewMessageWithMeta(apierrors.ErrorReason_NOT_FOUND, map[string]string{"protocol": cfg.Type}, "no protocol factory found for protocol: %s", cfg.Type)
	}

	server, err := factory.NewServer(cfg, opts...)
	if err != nil {
		return *new(T), errors.WrapAndConvert(err, apierrors.ErrorReason_INTERNAL_SERVER_ERROR, "failed to create server for protocol %s: %v", cfg.Type, err)
	}

	// Type assertion to ensure the created server matches the generic type T
	concreteServer, ok := server.(T)
	if !ok {
		return *new(T), errors.NewMessage(apierrors.ErrorReason_INTERNAL_SERVER_ERROR, "created server does not match expected type")
	}

	return concreteServer, nil
}

// NewBuilder returns a new generic ServerBuilder instance.
func NewBuilder[T interfaces.Server]() ServerBuilder[T] {
	return &serverBuilder[T]{}
}
