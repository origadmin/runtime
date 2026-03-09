/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package component

import (
	"context"
	"iter"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
)

type (
	Category string
	Scope    string
	Priority int
)

const (
	CategoryLogger      Category = "logger"
	CategoryRegistrar   Category = "registrar"
	CategoryDiscovery   Category = "discovery"
	CategoryClient      Category = "client"
	CategoryServer      Category = "server"
	CategoryMiddleware  Category = "middleware"
	CategoryDatabase    Category = "database"
	CategoryCache       Category = "cache"
	CategoryObjectStore Category = "objectstore"
	CategoryQueue       Category = "queue"
	CategoryTask        Category = "task"
	CategoryMail        Category = "mail"
	CategoryStorage     Category = "storage"
	CategorySecurity    Category = "security"
	CategorySkipper     Category = "skipper"
)

const (
	// GlobalScope is the default system fallback scope.
	GlobalScope Scope = "_global"
	// ServerScope is the standard scope for server-side components.
	ServerScope Scope = "server"
	// ClientScope is the standard scope for client-side components.
	ClientScope Scope = "client"

	// DefaultName is the system key for the active/default instance.
	DefaultName = "_default"
	// ReservedPrefix defines identifiers owned by the system.
	ReservedPrefix = "_"
)

// IsReserved checks if the metadata string is system-reserved.
func IsReserved(s string) bool {
	return len(s) > 0 && s[0] == '_'
}

type Locator interface {
	Get(ctx context.Context, name string) (any, error)
	Iter(ctx context.Context) iter.Seq2[string, any]
	In(cat Category, opts ...InOption) Locator
	Category() Category
	Scope() Scope
}

type Handle interface {
	Config() any
	Name() string
	Locator() Locator
}

type Provider func(ctx context.Context, h Handle) (any, error)

type Registry interface {
	Register(cat Category, p Provider, opts ...RegisterOption)
	Has(cat Category, opts ...RegisterOption) bool
	Load(ctx context.Context, source any, opts ...LoadOption) error
	In(cat Category, opts ...InOption) Locator
}

type Resolver func(source any, cat Category) (*ModuleConfig, error)

type Registration struct {
	Category Category
	Provider Provider
	Options  []RegisterOption
}

type ModuleConfig struct {
	Entries []ConfigEntry
	Active  string
}

type ConfigEntry struct {
	Name  string
	Value any
}

type RegistrationOptions struct {
	Resolver Resolver
	Scopes   []Scope
	Priority Priority
	Tag      string
}

type RegisterOption func(*RegistrationOptions)

type InOptions struct {
	Scope Scope
	Tags  []string
}

type InOption func(*InOptions)

type LoadOptions struct {
	Category Category
	Scope    Scope
	Name     string
	Resolver Resolver
	Tags     []string
}

type LoadOption func(*LoadOptions)

// Standard configuration interfaces
type AppConfig interface {
	GetApp() *appv1.App
}

type LoggerConfig interface {
	GetLogger() *loggerv1.Logger
}

type MiddlewareConfig interface {
	GetMiddlewares() *middlewarev1.Middlewares
}

type DataConfig interface {
	GetData() *datav1.Data
}
