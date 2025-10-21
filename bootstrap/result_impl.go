// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/interfaces"
)

// resultImpl implements the interfaces.Result interface.
type resultImpl struct {
	appInfo          *interfaces.AppInfo
	container        interfaces.Container
	config           interfaces.Config
	structuredConfig interfaces.StructuredConfig
	logger           log.Logger // Added logger field for configurable logging
	cleanup          func()
}

// StructuredConfig returns the structured configuration decoder.
func (b *resultImpl) StructuredConfig() interfaces.StructuredConfig {
	return b.structuredConfig
}

// Logger returns the logger instance.
func (b *resultImpl) Logger() log.Logger {
	return b.logger
}

// AppInfo returns the application information.
func (b *resultImpl) AppInfo() *interfaces.AppInfo {
	return b.appInfo
}

// Container returns the initialized component provider.
func (b *resultImpl) Container() interfaces.Container {
	return b.container
}

// Config implements interfaces.Result.
func (b *resultImpl) Config() interfaces.Config {
	return b.config
}

// Cleanup implements interfaces.Result.
func (b *resultImpl) Cleanup() func() {
	return b.cleanup
}
