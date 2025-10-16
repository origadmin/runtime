package service

import (
	"github.com/origadmin/runtime/errors"
)

const Module = "service"

// Predefined errors using the common StructuredError
var (
	ErrNilServerConfig     = newServiceError("server configuration is nil")
	ErrMissingServerConfig = newServiceError("missing server configuration")

	ErrNilClientConfig     = newServiceError("client configuration is nil")
	ErrMissingClientConfig = newServiceError("missing client configuration")

	ErrMissingProtocol     = newServiceError("protocol is not specified")
	ErrUnsupportedProtocol = newServiceError("unsupported protocol")
)

// Helper functions to create new errors
func newServiceError(message string) *errors.Structured {
	return errors.NewStructured(Module, message).WithCaller()
}
