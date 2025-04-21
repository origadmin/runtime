/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package gateway implements the functions, types, and interfaces for the module.
package gateway

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kHttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/runtime/log"
)

var (
	_ transport.Server     = (*Server)(nil)
	_ transport.Endpointer = (*Server)(nil)
)

type Framework interface {
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

type Server struct {
	Engine Framework
	serv   *http.Server

	tlsConf  *tls.Config
	endpoint *url.URL
	timeout  time.Duration
	addr     string

	err error

	filters []kHttp.FilterFunc
	ms      []middleware.Middleware
	dec     kHttp.DecodeRequestFunc
	enc     kHttp.EncodeResponseFunc
	ene     kHttp.EncodeErrorFunc
}

func NewServer(opts ...Option) *Server {
	srv := &Server{
		timeout: 1 * time.Second,
		dec:     kHttp.DefaultRequestDecoder,
		enc:     kHttp.DefaultResponseEncoder,
		ene:     kHttp.DefaultErrorEncoder,
	}

	srv.init(opts...)

	return srv
}

func (s *Server) init(opts ...Option) {
	s.Engine = gin.Default()

	for _, o := range opts {
		o(s)
	}

	s.serv = &http.Server{
		Addr:      s.addr,
		Handler:   s.Engine,
		TLSConfig: s.tlsConf,
	}

	s.endpoint, _ = url.Parse(s.addr)
}

func (s *Server) Endpoint() (*url.URL, error) {
	return s.endpoint, nil
}

func (s *Server) Start(ctx context.Context) error {
	log.Infof("[Gateway] serv listening on: %s", s.addr)

	var err error
	if s.tlsConf != nil {
		err = s.serv.ListenAndServeTLS("", "")
	} else {
		err = s.serv.ListenAndServe()
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[Gateway] serv stopping")
	return s.serv.Shutdown(ctx)
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.Engine.ServeHTTP(res, req)
}
