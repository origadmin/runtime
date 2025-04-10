/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package agent implements the functions, types, and interfaces for the module.
package agent

import (
	"fmt"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

type httpAgent struct {
	prefix  string
	version string
	server  *transhttp.Server
}

func (obj *httpAgent) HTTPServer() *transhttp.Server {
	return obj.server
}

func (obj *httpAgent) SetPrefix(prefix string) {
	obj.prefix = prefix
}

func (obj *httpAgent) SetVersion(version string) {
	obj.version = version
}

func (obj *httpAgent) Route() *transhttp.Router {
	return obj.server.Route(obj.URI())
}

func (obj *httpAgent) URI() string {
	return fmt.Sprintf("%s/%s", obj.prefix, obj.version)
}

func NewHTTP(server *transhttp.Server) HTTPAgent {
	return &httpAgent{
		prefix:  DefaultPrefix,
		version: DefaultVersion,
		server:  server,
	}
}
