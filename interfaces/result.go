package interfaces

import (
	"github.com/go-kratos/kratos/v2/log"
)

// Result defines the interface for the result of the bootstrap process.
// It provides access to the initialized Container, the Config decoder,
// and a cleanup function to release resources.
type Result interface {
	// AppInfo returns the application information.
	AppInfo() *AppInfo
	// Container returns the initialized component provider.
	Container() Container
	// Config returns the configuration decoder.
	Config() Config
	// StructuredConfig returns the structured configuration decoder.
	StructuredConfig() StructuredConfig
	// Logger returns the logger instance.
	Logger() log.Logger
	// Cleanup returns the cleanup function.
	Cleanup() func()
}
