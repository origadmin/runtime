package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/storage"
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
	// StorageProvider returns the configured storage provider.
	StorageProvider() storage.Provider
	// Cleanup returns the cleanup function.
	Cleanup() func()
}

// resultImpl implements the Result interface.
type resultImpl struct {
	config           interfaces.Config
	structuredConfig interfaces.StructuredConfig
	appInfo          *interfaces.AppInfo
	container        interfaces.Container
	logger           log.Logger
	storageProvider  storage.Provider // Add storage provider field
	cleanup          func()
}

// AppInfo returns the application information.
func (r *resultImpl) AppInfo() *interfaces.AppInfo {
	return r.appInfo
}

// Container returns the initialized component provider.
func (r *resultImpl) Container() interfaces.Container {
	return r.container
}

// Config returns the configuration decoder.
func (r *resultImpl) Config() interfaces.Config {
	return r.config
}

// StructuredConfig returns the structured configuration decoder.
func (r *resultImpl) StructuredConfig() interfaces.StructuredConfig {
	return r.structuredConfig
}

// Logger returns the logger instance.
func (r *resultImpl) Logger() log.Logger {
	return r.logger
}

// StorageProvider returns the configured storage provider.
func (r *resultImpl) StorageProvider() storage.Provider {
	return r.storageProvider
}

// Cleanup returns the cleanup function.
func (r *resultImpl) Cleanup() func() {
	return r.cleanup
}
