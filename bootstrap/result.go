package bootstrap

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/contracts"
)

// Result defines the interface for the raw data produced by the bootstrap process.
// It provides access to the loaded configuration and the raw application protobuf message.
type Result interface {
	// Config returns the raw configuration decoder.
	Config() contracts.ConfigLoader
	// StructuredConfig returns the structured configuration decoder,
	// which has merged defaults and applied transformers.
	StructuredConfig() contracts.StructuredConfig
	// AppConfig returns the raw protobuf App message decoded from the bootstrap configuration.
	// This message contains application-specific information as defined in the configuration file.
	AppConfig() *appv1.App
	// Bootstrap returns the original bootstrap configuration object as an any type.
	// This allows the runtime to perform interface sniffing on user-defined configurations.
	Bootstrap() any
}
