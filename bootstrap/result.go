package bootstrap

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/interfaces"
)

// Result defines the interface for the raw data produced by the bootstrap process.
// It provides access to the loaded configuration and the raw application protobuf message.
type Result interface {
	// Config returns the raw configuration decoder.
	Config() interfaces.Config
	// StructuredConfig returns the structured configuration decoder,
	// which has merged defaults and applied transformers.
	StructuredConfig() interfaces.StructuredConfig
	// App returns the raw protobuf App message decoded from the bootstrap configuration.
	// This message contains application-specific information as defined in the configuration file.
	App() *appv1.App
}
