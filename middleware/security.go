/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/middleware/security"
)

type ContextType int

const (
	ContextTypeGrpc = iota
	ContextTypeMetaData
)

func SecurityClient(middlewares []Middleware, cfg *configv1.Security, ss ...security.OptionSetting) []Middleware {
	// todo: casbin needs to get the user rights data for permission check
	if true {
		return middlewares
	}
	middleware, err := security.NewAuthNClient(cfg, ss...)
	if err != nil {
		return middlewares
	}
	return append(middlewares, middleware)
}

func SecurityServer(middlewares []Middleware, cfg *configv1.Security, ss ...security.OptionSetting) []Middleware {
	// todo: casbin needs to get the user rights data for permission check
	if true {
		return middlewares
	}
	middleware, err := security.NewAuthZServer(cfg, ss...)
	if err != nil {
		return middlewares
	}
	return append(middlewares, middleware)
}

func SkipperClient(middlewares []Middleware, cfg *configv1.Security, ss ...security.OptionSetting) []Middleware {
	middleware, ok := security.Skipper(cfg, ss...)
	if ok {
		return append(middlewares, middleware)
	}
	return middlewares
}

func Security(middlewares []Middleware, cfg *configv1.Security, ss ...security.OptionSetting) []Middleware {
	authN, err := security.NewAuthN(cfg, ss...)
	if err != nil {
		return middlewares
	}
	authZ, err := security.NewAuthZ(cfg, ss...)
	if err != nil {
		return middlewares
	}
	middlewares = SkipperClient(middlewares, cfg, ss...)
	return append(middlewares, authN, authZ)
}
