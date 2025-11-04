package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/interfaces"
)

// Result defines the interface for the result of the bootstrap process.
// It provides access to the initialized Container, the Config decoder,
// and a cleanup function to release resources.
type Result interface {
	// AppInfo returns the application information.
	AppInfo() *interfaces.AppInfo
	// Container returns the initialized component provider.
	Container() interfaces.Container
	// Config returns the configuration decoder.
	Config() interfaces.Config
	// StructuredConfig returns the structured configuration decoder.
	StructuredConfig() interfaces.StructuredConfig
	// Logger returns the logger instance.
	Logger() log.Logger
	// Cleanup returns the cleanup function.
	Cleanup() func()
}
