/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/metadata"
	middlewareMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
)

func MetadataClient(ms []Middleware, ok bool, cmm *configv1.Middleware_Metadata) []Middleware {
	if !ok {
		return ms
	}

	var options []middlewareMetadata.Option
	if prefix := cmm.GetPrefix(); prefix != "" {
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cmm.GetData(); len(metaSource) > 0 {
		data := make(metadata.Metadata, len(metaSource))
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	return append(ms, middlewareMetadata.Client(options...))
}

func MetadataServer(ms []Middleware, ok bool, cmm *configv1.Middleware_Metadata) []Middleware {
	if !ok {
		return ms
	}

	var options []middlewareMetadata.Option
	if prefix := cmm.GetPrefix(); prefix != "" {
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cmm.GetData(); len(metaSource) > 0 {
		data := metadata.Metadata{}
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	return append(ms, middlewareMetadata.Server(options...))
}
