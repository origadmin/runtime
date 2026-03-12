package bootstrap

import (
	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/config/bootstrap/v1"
	"github.com/origadmin/runtime/config"
)

// Result defines the unified contract for the bootstrap engine output.
type Result interface {
	// Bootstrap [Source Phase] Returns the strong-typed bootstrap metadata (sources, service info, etc.)
	Bootstrap() *bootstrapv1.Bootstrap

	// Config [Binding Phase] Returns the final decoded business configuration (any type).
	Config() any

	// Decoder [Operations] Returns the enhanced configuration loader (follows Kratos design).
	// It provides rich operations like Value(), Watch(), and Scan().
	Decoder() config.KConfig

	// ConfigPath returns the physical path of the loaded configuration file.
	ConfigPath() string
}
