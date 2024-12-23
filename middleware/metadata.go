/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package middleware implements the functions, types, and interfaces for the module.
package middleware

import (
	"github.com/go-kratos/kratos/v2/metadata"
	middlewareMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"

	middlewarev1 "github.com/origadmin/runtime/gen/go/middleware/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware/selector"
)

func MetadataClient(selector selector.Selector, cfg *middlewarev1.Middleware_Metadata) selector.Selector {
	log.Debug("[MetadataClient] KMiddleware is enabled")
	var options []middlewareMetadata.Option
	if prefix := cfg.GetPrefix(); prefix != "" {
		log.Debug("[MetadataClient] Propagated prefix: ", prefix)
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cfg.GetData(); len(metaSource) > 0 {
		log.Debug("[MetadataClient] Metadata source: ", metaSource)
		data := make(metadata.Metadata, len(metaSource))
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	log.Debug("[MetadataClient] Metadata client middleware enabled")
	return selector.Append("Metadata", middlewareMetadata.Client(options...))
}

func MetadataServer(selector selector.Selector, cfg *middlewarev1.Middleware_Metadata) selector.Selector {
	log.Debug("[MetadataServer] KMiddleware is enabled")

	var options []middlewareMetadata.Option
	if prefix := cfg.GetPrefix(); prefix != "" {
		log.Debug("[MetadataServer] Propagated prefix: ", prefix)
		options = append(options, middlewareMetadata.WithPropagatedPrefix(prefix))
	}
	if metaSource := cfg.GetData(); len(metaSource) > 0 {
		log.Debug("[MetadataServer] Metadata source: ", metaSource)
		data := metadata.Metadata{}
		for k, v := range metaSource {
			data[k] = []string{v}
		}
		options = append(options, middlewareMetadata.WithConstants(data))
	}
	log.Debug("[MetadataServer] Metadata server middleware enabled")
	return selector.Append("Metadata", middlewareMetadata.Server(options...))
}
