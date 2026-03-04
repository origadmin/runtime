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
	GlobalScope       = component.GlobalScope
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

// WithExtractor specifies a local config extractor.
func WithExtractor(e component.Extractor) RegisterOption {
	return engine.WithExtractor(e)
}
