/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package runtime

import (
	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/runtime/engine"
)

// --- Engine Metadata Types Aliases ---

type Category = component.Category
type Scope = component.Scope
type Priority = component.Priority
type Provider = component.Provider
type Handle = component.Handle

const (
	// GlobalScope is the default fallback scope for the system.
	GlobalScope = component.GlobalScope
)

// --- Category Conventions ---

const (
	CategoryInfrastructure Category = "infrastructure"
	CategoryLogger                  = component.CategoryLogger
	CategoryRegistrar               = component.CategoryRegistrar
	CategoryDiscovery               = component.CategoryDiscovery
	CategoryClient                  = component.CategoryClient
	CategoryServer                  = component.CategoryServer
	CategoryMiddleware              = component.CategoryMiddleware
	CategoryDatabase                = component.CategoryDatabase
	CategoryCache                   = component.CategoryCache
	CategoryObjectStore             = component.CategoryObjectStore
	CategoryQueue                   = component.CategoryQueue
	CategoryTask                    = component.CategoryTask
	CategoryMail                    = component.CategoryMail
	CategoryStorage                 = component.CategoryStorage
	CategorySecurity                = component.CategorySecurity
	CategorySkipper                 = component.CategorySkipper
)

// --- Scope Conventions ---

const (
	ServerScope = component.ServerScope
	ClientScope = component.ClientScope
)

// --- Priority Conventions ---

const (
	PriorityInfrastructure Priority = 100
	PriorityRegistry       Priority = 200
	PriorityStorage        Priority = 300
	PriorityClientStack    Priority = 400
	PriorityServerStack    Priority = 500
)

// --- Global Registration (init phase) ---

// Register registers a component capability to the global pool.
// Typically used in init() functions of component packages.
func Register(cat Category, p Provider, opts ...RegisterOption) {
	engine.Register(cat, p, opts...)
}

// --- Functional Option Type Aliases ---

type RegisterOption = component.RegisterOption
type InOption = component.InOption
type LoadOption = component.LoadOption

// --- Functional Option Helpers ---

// WithScope specifies the perspective during handle creation (In).
func WithScope(s Scope) InOption {
	return engine.WithInScope(s)
}

// WithScopes specifies the visibilities during registration (Register).
func WithScopes(ss ...Scope) RegisterOption {
	return engine.WithScopes(ss...)
}

// WithPriority specifies the initialization priority.
func WithPriority(p Priority) RegisterOption {
	return engine.WithPriority(p)
}

// WithTags specifies the tags for a component.
func WithTags(tags ...string) RegisterOption {
	return engine.WithTags(tags...)
}

// WithInTags specifies the tags for perspective switching.
func WithInTags(tags ...string) InOption {
	return engine.WithInTags(tags...)
}

// WithResolver specifies a local config resolver for a component.
func WithResolver(res component.Resolver) RegisterOption {
	return engine.WithResolverOption(res)
}
