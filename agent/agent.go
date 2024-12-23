/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/transport/http"
)

const (
	ApiVersionV1   = "/api/v1"
	DefaultPrefix  = "/api"
	DefaultVersion = "v1"
)

type Agent interface {
	URI() string
	Server() *http.Server
	Route() *http.Router
}

type agent struct {
	prefix  string
	version string
	server  *http.Server
}

func (obj *agent) SetPrefix(prefix string) {
	obj.prefix = prefix
}

func (obj *agent) SetVersion(version string) {
	obj.version = version
}

func (obj *agent) Server() *http.Server {
	return obj.server
}

func (obj *agent) Route() *http.Router {
	return obj.server.Route(obj.URI())
}

func (obj *agent) URI() string {
	return fmt.Sprintf("%s/%s", obj.prefix, obj.version)
}

func New(server *http.Server) Agent {
	return &agent{
		prefix:  DefaultPrefix,
		version: DefaultVersion,
		server:  server,
	}
}
