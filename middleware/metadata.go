/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/metadata"
	middlewareMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"

	configv1 "github.com/origadmin/runtime/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
)

func MetadataClient(ms []Middleware, ok bool, cmm *configv1.Middleware_Metadata) []Middleware {
	if !ok {
		log.Debug("[MetadataClient] Middleware is not enabled")
		return ms
	}

	log.Debug("[MetadataClient] Middleware is enabled")

	var options []middlewareMetadata.Option
	if prefix := cmm.GetPrefix(); prefix != "" {
		log.Debug("[MetadataClient] Propagated prefix: ", prefix)
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cmm.GetData(); len(metaSource) > 0 {
		log.Debug("[MetadataClient] Metadata source: ", metaSource)
		data := make(metadata.Metadata, len(metaSource))
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	log.Debug("[MetadataClient] Options: ", options)
	log.Debug("[MetadataClient] Metadata client middleware enabled")
	return append(ms, middlewareMetadata.Client(options...))
}

func MetadataServer(ms []Middleware, ok bool, cmm *configv1.Middleware_Metadata) []Middleware {
	if !ok {
		log.Debug("[MetadataServer] Middleware is not enabled")
		return ms
	}

	log.Debug("[MetadataServer] Middleware is enabled")

	var options []middlewareMetadata.Option
	if prefix := cmm.GetPrefix(); prefix != "" {
		log.Debug("[MetadataServer] Propagated prefix: ", prefix)
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cmm.GetData(); len(metaSource) > 0 {
		log.Debug("[MetadataServer] Metadata source: ", metaSource)
		data := metadata.Metadata{}
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	log.Debug("[MetadataServer] Options: ", options)
	log.Debug("[MetadataServer] Metadata server middleware enabled")
	return append(ms, middlewareMetadata.Server(options...))
}
