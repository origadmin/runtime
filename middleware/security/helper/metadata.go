/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package helper implements the functions, types, and interfaces for the module.
package helper

import (
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/origadmin/toolkits/security"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/origadmin/runtime/context"
)

// ContextWithToken .
func ContextWithToken(ctx context.Context, expectedScheme string, tokenStr string, ctxType security.ContextType) context.Context {
	switch ctxType {
	case "context":
		return injectTokenToGrpcContext(ctx, expectedScheme, tokenStr)
	case "metadata":
		return injectTokenToKratosContext(ctx, expectedScheme, tokenStr)
	default:
		return injectTokenToGrpcContext(ctx, expectedScheme, tokenStr)
	}
}

// TokenFromContext .
func TokenFromContext(ctx context.Context, expectedScheme string, ctxType security.ContextType) (string, error) {
	val := extractTokenFromContext(ctx, ctxType)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}

	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization string")
	}

	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated with "+expectedScheme)
	}

	return splits[1], nil
}

func extractTokenFromGrpcContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	return md.Get("Authorization")[0]
}

func extractTokenFromKratosContext(ctx context.Context) string {
	if header, ok := transport.FromServerContext(ctx); ok {
		return header.RequestHeader().Get("Authorization")
	}
	return ""
}

func extractTokenFromContext(ctx context.Context, ctxType security.ContextType) string {
	switch ctxType {
	case "context":
		return extractTokenFromGrpcContext(ctx)
	case "metadata":
		return extractTokenFromKratosContext(ctx)
	default:
		return extractTokenFromGrpcContext(ctx)
	}
}

func formatToken(expectedScheme string, tokenStr string) string {
	return fmt.Sprintf("%s %s", expectedScheme, tokenStr)
}

func injectTokenToKratosContext(ctx context.Context, expectedScheme string, tokenStr string) context.Context {
	if header, ok := transport.FromClientContext(ctx); ok {
		header.RequestHeader().Set("Authorization", formatToken(expectedScheme, tokenStr))
	} else {
		//log.Error("authn token injection failure in kratos context")
	}
	return ctx
}

func injectTokenToGrpcContext(ctx context.Context, expectedScheme string, tokenStr string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	md.Set("Authorization", formatToken(expectedScheme, tokenStr))
	return ctx
}
