package component

import (
	"context"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/contracts/options"
	"github.com/origadmin/runtime/engine/metadata"
)

// --- Configuration Sniffing Contracts ---

type (
	AppConfig        interface{ GetApp() *appv1.App }
	LoggerConfig     interface{ GetLogger() *loggerv1.Logger }
	MiddlewareConfig interface {
		GetMiddlewares() *middlewarev1.Middlewares
	}
	DataConfig     interface{ GetData() *datav1.Data }
	RegistryConfig interface {
		GetDiscoveries() *discoveryv1.Discoveries
	}
)

// --- Engine Core Contracts ---

type Handle interface {
	Get(ctx context.Context, name string) (any, error)
	In(category metadata.Category, opts ...RegisterOption) Handle
	BindConfig(target any) error
	Config() any
	Scope() metadata.Scope
	Category() metadata.Category
}

type Provider func(ctx context.Context, h Handle, opts ...options.Option) (any, error)

type ConfigEntry struct {
	Name  string
	Value any
}

type ModuleConfig struct {
	Entries []ConfigEntry
	Active  string
}

type Extractor func(root any) (*ModuleConfig, error)

type Registry interface {
	Handle
	Register(c metadata.Category, e Extractor, p Provider, opts ...RegisterOption)
	Has(c metadata.Category, opts ...RegisterOption) bool
	Init(ctx context.Context, root any) error
}

type RegisterOption func(any)
