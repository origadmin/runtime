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

const (
	// GlobalScope is the default fallback scope for the system.
	GlobalScope = component.GlobalScope
)

// --- Category Conventions ---

const (
	CategoryInfrastructure Category = "infrastructure"
	CategoryLogger         Category = "logger"
	CategoryRegistry       Category = "registry"
	CategoryClient         Category = "client"
	CategoryServer         Category = "server"
	CategoryMiddleware     Category = "middleware"
	CategoryDatabase       Category = "database"
	CategoryCache          Category = "cache"
	CategoryObjectStore    Category = "objectstore"
	CategoryQueue          Category = "queue"
	CategoryTask           Category = "task"
	CategoryMail           Category = "mail"
	CategoryStorage        Category = "storage"
)

// --- Scope Conventions ---

const (
	ServerScope Scope = "server"
	ClientScope Scope = "client"
)

// --- Priority Conventions ---

const (
	PriorityInfrastructure Priority = 100
	PriorityRegistry       Priority = 200
	PriorityStorage        Priority = 300
	PriorityClientStack    Priority = 400
	PriorityServerStack    Priority = 500
)

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

// WithResolver specifies a local config resolver for a component.
func WithResolver(res component.Resolver) RegisterOption {
	return engine.WithResolverOption(res)
}
