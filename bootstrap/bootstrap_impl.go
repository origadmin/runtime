// Package bootstrap implements the functions, types, and interfaces for the module.
package bootstrap

import (
	"github.com/origadmin/runtime/interfaces"
)

// bootstrapperImpl implements the interfaces.Bootstrapper interface.
type bootstrapperImpl struct {
	appInfo   *interfaces.AppInfo
	container interfaces.Container
	config    interfaces.StructuredConfig
	cleanup   func()
}

func (b *bootstrapperImpl) AppInfo() *interfaces.AppInfo {
	return b.appInfo
}

func (b *bootstrapperImpl) Container() interfaces.Container {
	return b.container
}

// Config implements interfaces.Bootstrapper.
func (b *bootstrapperImpl) Config() interfaces.StructuredConfig {
	return b.config
}

// Cleanup implements interfaces.Bootstrapper.
func (b *bootstrapperImpl) Cleanup() func() {
	return b.cleanup
}
