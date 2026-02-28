package bootstrap

import (
	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/config/bootstrap/v1"
	"github.com/origadmin/runtime/config"
	"github.com/origadmin/runtime/contracts"
)

// Result defines the unified contract for the bootstrap engine output.
type Result interface {
	// Bootstrap [Source Phase] Returns the strong-typed bootstrap metadata (sources, service info, etc.)
	Bootstrap() *bootstrapv1.Bootstrap

	// Config [Binding Phase] Returns the final decoded business configuration (any type).
	Config() any

	// Loader [Operations] Returns the enhanced configuration loader (follows Kratos design).
	// It provides rich operations like Value(), Watch(), and Scan().
	Loader() config.KConfig

	// ConfigPath returns the physical path of the loaded configuration file.
	ConfigPath() string

	// StructuredConfig returns the legacy structured configuration decoder.
	// This is kept for backward compatibility with Container and App initialization.
	StructuredConfig() contracts.StructuredConfig
}
