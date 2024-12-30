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

func injectServerTransportContext(ctx context.Context, scheme string, token string) context.Context {
	if header, ok := transport.FromServerContext(ctx); ok {
		header.RequestHeader().Set(security.HeaderAuthorize, formatToken(scheme, token))
	}
	return ctx
}

func injectClientTransportContext(ctx context.Context, scheme string, token string) context.Context {
	if header, ok := transport.FromClientContext(ctx); ok {
		header.RequestHeader().Set(security.HeaderAuthorize, formatToken(scheme, token))
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
	case security.ContextTypeServerHeader:
		return injectServerTransportContext(ctx, scheme, token)
	case security.ContextTypeClientHeader:
		return injectClientTransportContext(ctx, scheme, token)
	case security.ContextTypeContext:
		return security.NewTokenContext(ctx, formatToken(scheme, token))
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

func extractServerTransportContext(ctx context.Context) string {
	if header, ok := transport.FromServerContext(ctx); ok {
		return header.RequestHeader().Get(security.HeaderAuthorize)
	}
	return ""
}

func extractClientTransportContext(ctx context.Context) string {
	if header, ok := transport.FromClientContext(ctx); ok {
		return header.RequestHeader().Get(security.HeaderAuthorize)
	}
	return ""
}

func extractTokenFromContext(ctx context.Context, tokenType security.TokenType) string {
	switch tokenType {
	case security.ContextTypeMetadata:
		return extractTokenMetadataContext(ctx)
	case security.ContextTypeServerHeader:
		return extractServerTransportContext(ctx)
	case security.ContextTypeClientHeader:
		return extractClientTransportContext(ctx)
	case security.ContextTypeContext:
		return security.TokenFromContext(ctx)
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
		return security.ClaimsFromContext(ctx), nil
	}
	return nil, status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+tokenType.String())
}

func formatToken(scheme string, tokenStr string) string {
	return fmt.Sprintf("%s %s", scheme, tokenStr)
}

func FromTransportClient(authorize string, scheme string) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		if tr, ok := transport.FromClientContext(ctx); ok {
			token := tr.RequestHeader().Get(authorize)
			splits := strings.SplitN(token, " ", 2)
			if len(splits) > 1 && strings.EqualFold(splits[0], scheme) {
				return splits[1]
			}
		}
		return ""
	}
}

func FromTransportServer(authorize string, scheme string) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		if tr, ok := transport.FromServerContext(ctx); ok {
			token := tr.RequestHeader().Get(authorize)
			splits := strings.SplitN(token, " ", 2)
			if len(splits) > 1 && strings.EqualFold(splits[0], scheme) {
				return splits[1]
			}
		}
		return ""
	}
}
