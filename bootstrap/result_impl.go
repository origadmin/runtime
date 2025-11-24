package bootstrap

import (
	"github.com/origadmin/runtime/interfaces"
)

// resultImpl implements the Result interface.
type resultImpl struct {
	config           interfaces.Config
	structuredConfig interfaces.StructuredConfig
}

// Config returns the raw configuration decoder.
func (b *resultImpl) Config() interfaces.Config {
	return b.config
}

// StructuredConfig returns the structured configuration decoder.
func (b *resultImpl) StructuredConfig() interfaces.StructuredConfig {
	return b.structuredConfig
}
