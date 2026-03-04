/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package component

import (
	"context"
	"iter"

	appv1 "github.com/origadmin/runtime/api/gen/go/config/app/v1"
	datav1 "github.com/origadmin/runtime/api/gen/go/config/data/v1"
	discoveryv1 "github.com/origadmin/runtime/api/gen/go/config/discovery/v1"
	loggerv1 "github.com/origadmin/runtime/api/gen/go/config/logger/v1"
	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/contracts/options"
)

// --- Engine Metadata Types ---

type Scope string
type Category string
type Priority int

const (
	// GlobalScope is the default fallback scope for the system.
	GlobalScope Scope = "_global"
	// ReservedPrefix is the prefix used for system-reserved identifiers.
	ReservedPrefix = "_"
)

// IsReserved checks if the given identifier belongs to the system-reserved namespace.
func IsReserved(s string) bool {
	return len(s) > 0 && s[0] == '_'
}

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
	Iter(ctx context.Context) iter.Seq2[string, any]
	In(category Category, opts ...InOption) Handle
	BindConfig(target any) error
	Config() any
	Scope() Scope
	Category() Category
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
	// Register declares a component capability.
	Register(c Category, p Provider, opts ...RegisterOption)
	// Has checks if a component category is registered.
	Has(c Category, opts ...RegisterOption) bool
	// Load injects a configuration source and triggers binding.
	Load(ctx context.Context, source any, opts ...LoadOption) error
}

// --- Option Definitions ---

type RegistrationOptions struct {
	Extractor Extractor
	Scopes    []Scope
	Priority  Priority
}

type RegisterOption func(*RegistrationOptions)

type InOptions struct {
	Scope Scope
}

type InOption func(*InOptions)

type LoadOptions struct {
	Category Category
	Name     string
}

type LoadOption func(*LoadOptions)
