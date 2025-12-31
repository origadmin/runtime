package bootstrap

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/interfaces"
)

// resultImpl implements the Result interface.
type resultImpl struct {
	config           interfaces.ConfigLoader
	structuredConfig interfaces.StructuredConfig
	appConfig        *appv1.App
}

func (b *resultImpl) AppConfig() *appv1.App {
	return b.appConfig
}

// Config returns the raw configuration decoder.
func (b *resultImpl) Config() interfaces.ConfigLoader {
	return b.config
}

// StructuredConfig returns the structured configuration decoder.
func (b *resultImpl) StructuredConfig() interfaces.StructuredConfig {
	return b.structuredConfig
}
