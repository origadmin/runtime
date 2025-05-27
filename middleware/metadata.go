/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/metadata"
	middlewareMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/middleware/v1"
)

func MetadataClient(ms []KMiddleware, cfg *middlewarev1.Middleware_Metadata) []KMiddleware {
	var options []middlewareMetadata.Option
	if prefix := cfg.GetPrefix(); prefix != "" {
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cfg.GetData(); len(metaSource) > 0 {
		data := make(metadata.Metadata, len(metaSource))
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	return append(ms, middlewareMetadata.Client(options...))
}

func MetadataServer(ms []KMiddleware, cfg *middlewarev1.Middleware_Metadata) []KMiddleware {
	var options []middlewareMetadata.Option
	if prefix := cfg.GetPrefix(); prefix != "" {
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cfg.GetData(); len(metaSource) > 0 {
		data := metadata.Metadata{}
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	return append(ms, middlewareMetadata.Server(options...))
}

type metadataFactory struct {
}

func (m metadataFactory) NewMiddlewareClient(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	cfg := middleware.GetMetadata()
	if cfg.GetEnabled() {
		options := make([]middlewareMetadata.Option, 0)
		if prefix := cfg.GetPrefix(); prefix != "" {
			options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
		}
		if metaSource := cfg.GetData(); len(metaSource) > 0 {
			data := make(metadata.Metadata, len(metaSource))
			for k, v := range metaSource {
				data[k] = []string{v}
			}
			options = append(options, middlewareMetadata.WithConstants(data))
		}
		return middlewareMetadata.Client(options...), true
	}
	return nil, false
}

func (m metadataFactory) NewMiddlewareServer(middleware *middlewarev1.Middleware, options *Options) (KMiddleware, bool) {
	cfg := middleware.GetMetadata()
	if cfg.GetEnabled() {
		options := make([]middlewareMetadata.Option, 0)
		if prefix := cfg.GetPrefix(); prefix != "" {
			options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
		}
		if metaSource := cfg.GetData(); len(metaSource) > 0 {
			data := metadata.Metadata{}
			for k, v := range metaSource {
				data[k] = []string{v}
			}
			options = append(options, middlewareMetadata.WithConstants(data))
		}
		return middlewareMetadata.Server(options...), true
	}
	return nil, false
}
