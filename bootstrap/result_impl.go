package bootstrap

import (
	bootstrapv1 "github.com/origadmin/runtime/api/gen/go/config/bootstrap/v1"
	"github.com/origadmin/runtime/contracts"
)

// resultImpl implements the Result interface for the bootstrap engine.
type resultImpl struct {
	config           contracts.ConfigLoader
	structuredConfig contracts.StructuredConfig
	bootstrap        *bootstrapv1.Bootstrap
	businessConfig   any
	configPath       string
}

// Bootstrap returns the strong-typed bootstrap metadata.
func (b *resultImpl) Bootstrap() *bootstrapv1.Bootstrap {
	return b.bootstrap
}

// Config returns the decoded business configuration object (any).
func (b *resultImpl) Config() any {
	return b.businessConfig
}

// Loader returns the raw configuration hub.
func (b *resultImpl) Loader() contracts.ConfigLoader {
	return b.config
}

// ConfigPath returns the physical configuration path.
func (b *resultImpl) ConfigPath() string {
	return b.configPath
}

// StructuredConfig returns the legacy structured configuration decoder.
func (b *resultImpl) StructuredConfig() contracts.StructuredConfig {
	return b.structuredConfig
}
