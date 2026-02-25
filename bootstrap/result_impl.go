package bootstrap

import (
	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	"github.com/origadmin/runtime/contracts"
)

// resultImpl implements the Result interface.
type resultImpl struct {
	config           contracts.ConfigLoader
	structuredConfig contracts.StructuredConfig
	appConfig        *appv1.App
	rawConfig        any
}

func (b *resultImpl) AppConfig() *appv1.App {
	return b.appConfig
}

// Config returns the raw configuration decoder.
func (b *resultImpl) Config() contracts.ConfigLoader {
	return b.config
}

// StructuredConfig returns the structured configuration decoder.
func (b *resultImpl) StructuredConfig() contracts.StructuredConfig {
	return b.structuredConfig
}

// RawConfig returns the original bootstrap configuration object.
func (b *resultImpl) RawConfig() any {
	return b.rawConfig
}
