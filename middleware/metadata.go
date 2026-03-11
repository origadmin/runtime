/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and contracts for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/metadata"
	middlewareMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
)

type metadataFactory struct {
}

func (m metadataFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	metadataConfig := cfg.GetMetadata()
	if metadataConfig == nil {
		return nil, false
	}

	var metadataOpts []middlewareMetadata.Option
	if prefixes := metadataConfig.GetPrefixes(); len(prefixes) > 0 {
		metadataOpts = append(metadataOpts, middlewareMetadata.WithPropagatedPrefix(prefixes...))
	}
	if metaSource := metadataConfig.GetData(); len(metaSource) > 0 {
		data := make(metadata.Metadata, len(metaSource))
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		metadataOpts = append(metadataOpts, middlewareMetadata.WithConstants(data))
	}
	return middlewareMetadata.Client(metadataOpts...), true
}

func (m metadataFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...Option) (KMiddleware, bool) {
	metadataConfig := cfg.GetMetadata()
	if metadataConfig == nil {
		return nil, false
	}

	var metadataOpts []middlewareMetadata.Option
	if prefixes := metadataConfig.GetPrefixes(); len(prefixes) > 0 {
		metadataOpts = append(metadataOpts, middlewareMetadata.WithPropagatedPrefix(prefixes...))
	}
	if metaSource := metadataConfig.GetData(); len(metaSource) > 0 {
		data := metadata.Metadata{}
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		metadataOpts = append(metadataOpts, middlewareMetadata.WithConstants(data))
	}
	return middlewareMetadata.Server(metadataOpts...), true
}
