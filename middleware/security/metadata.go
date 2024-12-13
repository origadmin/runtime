/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/origadmin/runtime/context"
	"github.com/origadmin/toolkits/security"
)

func injectTokenTransportContext(ctx context.Context, scheme string, token string) context.Context {
	if header, ok := transport.FromClientContext(ctx); ok {
		header.RequestHeader().Set(security.HeaderAuthorize, formatToken(scheme, token))
	} else {
		//log.Error("authn token injection failure in kratos context")
	}
	return ctx
}

func injectTokenMetadataContext(ctx context.Context, scheme string, token string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		// Use pairs to create a new one.
		md = metadata.Pairs()
	}
	md.Set(security.HeaderAuthorize, formatToken(scheme, token))
	return ctx
}

// TokenToTypeContext .
func TokenToTypeContext(ctx context.Context, tokenType security.TokenType, scheme string, token string) context.Context {
	switch tokenType {
	case security.ContextTypeMetadata:
		return injectTokenMetadataContext(ctx, scheme, token)
	case security.ContextTypeHeader:
		return injectTokenTransportContext(ctx, scheme, token)
	case security.ContextTypeContext:
		return NewTokenContext(ctx, formatToken(scheme, token))
	default:
		return injectTokenMetadataContext(ctx, scheme, token)
	}
}

func extractTokenMetadataContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md.Get(security.HeaderAuthorize)[0]
	}
	return ""
}

func extractTokenTransportContext(ctx context.Context) string {
	if header, ok := transport.FromServerContext(ctx); ok {
		return header.RequestHeader().Get(security.HeaderAuthorize)
	}
	return ""
}

func extractTokenFromContext(ctx context.Context, tokenType security.TokenType) string {
	switch tokenType {
	case security.ContextTypeMetadata:
		return extractTokenMetadataContext(ctx)
	case security.ContextTypeHeader:
		return extractTokenTransportContext(ctx)
	case security.ContextTypeContext:
		return TokenFromContext(ctx)
	default:
		return extractTokenMetadataContext(ctx)
	}
}

// TokenFromTypeContext .
func TokenFromTypeContext(ctx context.Context, tokenType security.TokenType, scheme string) (string, error) {
	val := extractTokenFromContext(ctx, tokenType)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+scheme)
	}

	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}

	if !strings.EqualFold(splits[0], scheme) {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+scheme)
	}

	return splits[1], nil
}

func ClaimFromTokenTypeContext(ctx context.Context, tokenType security.TokenType) (security.Claims, error) {
	switch tokenType {
	case security.ContextTypeContext:
		return ClaimsFromContext(ctx), nil
	}
	return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+string(tokenType))
}

func formatToken(scheme string, tokenStr string) string {
	return fmt.Sprintf("%s %s", scheme, tokenStr)
}
