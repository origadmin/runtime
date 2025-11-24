package bootstrap

import (
	"github.com/origadmin/runtime/interfaces"
)

// Result defines the interface for the raw data produced by the bootstrap process.
// It provides access to the loaded configuration.
type Result interface {
	// Config returns the raw configuration decoder.
	Config() interfaces.Config
	// StructuredConfig returns the structured configuration decoder,
	// which has merged defaults and applied transformers.
	StructuredConfig() interfaces.StructuredConfig
}
