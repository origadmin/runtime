package bootstrap

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/interfaces"
)

// Result defines the interface for the raw data produced by the bootstrap process.
// It provides access to the loaded configuration and application information.
type Result interface {
	// AppInfo returns the application information.
	// Bootstrap is responsible for loading this data, but not for its validation.
	AppInfo() *appv1.App
	// Config returns the raw configuration decoder.
	Config() interfaces.Config
	// StructuredConfig returns the structured configuration decoder,
	// which has merged defaults and applied transformers.
	StructuredConfig() interfaces.StructuredConfig
}
