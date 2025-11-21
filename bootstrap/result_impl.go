package bootstrap

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1" // Added for appv1.App
	"github.com/origadmin/runtime/interfaces"
)

// resultImpl implements the Result interface.
type resultImpl struct {
	appInfo          *appv1.App // Updated type
	config           interfaces.Config
	structuredConfig interfaces.StructuredConfig
}

// AppInfo returns the application information.
func (b *resultImpl) AppInfo() *appv1.App {
	return b.appInfo
}

// Config returns the raw configuration decoder.
func (b *resultImpl) Config() interfaces.Config {
	return b.config
}

// StructuredConfig returns the structured configuration decoder.
func (b *resultImpl) StructuredConfig() interfaces.StructuredConfig {
	return b.structuredConfig
}
