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

type Registry = component.Registry
type Container = component.Container

// --- Category Conventions ---

const (
	CategoryInfrastructure Category = "infrastructure"
	CategoryLogger         Category = "logger"
	CategoryRegistrar      Category = "registrar"
	CategoryDiscovery      Category = "discovery"
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
	CategorySecurity       Category = "security"
	CategorySkipper        Category = "skipper"
)

// --- Scope Conventions ---

const (
	// GlobalScope is the default fallback scope for the system.
	GlobalScope Scope = ""
	// ServerScope is the standard scope for server-side components.
	ServerScope Scope = "server"
	// ClientScope is the standard scope for client-side components.
	ClientScope Scope = "client"
)

// --- Engine Component Aliases ---

type (
	ConfigResolver      = component.ConfigResolver
	RequirementResolver = component.RequirementResolver
	Registration        = component.Registration
	ModuleConfig        = component.ModuleConfig
	ConfigEntry         = component.ConfigEntry
	RegistrationOptions = component.RegistrationOptions

	InOption   = component.InOption
	Locator    = component.Locator
	LoadOption = component.LoadOption
)

type (
	AppOption      = Option
	RegisterOption = component.RegisterOption
)

// --- Engine Options (Perspective & Load) ---

// WithInScope specifies the perspective scope.
func WithInScope(s Scope) InOption {
	return engine.WithInScope(s)
}

// WithInTags specifies the perspective tags.
func WithInTags(tags ...string) InOption {
	return engine.WithInTags(tags...)
}

// WithResolver specifies a local config resolver for a component.
func WithResolver(res component.ConfigResolver) RegisterOption {
	return engine.WithConfigResolverOption(res)
}

// WithScopes specifies the scopes for a component.
func WithScopes(ss ...Scope) RegisterOption {
	return engine.WithScopes(ss...)
}
