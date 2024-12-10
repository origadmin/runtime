/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package security

import (
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/origadmin/toolkits/security"

	"github.com/origadmin/runtime/context"
	securityv1 "github.com/origadmin/runtime/gen/go/security/v1"
	"github.com/origadmin/runtime/middleware/security/helper"
)

const (
	reason string = "FORBIDDEN"
)

var (
	ErrInvalidToken  = securityv1.ErrorAuthErrorReasonBearerTokenMissing("bearer token missing")
	ErrInvalidClaims = securityv1.ErrorAuthErrorReasonInvalidClaims("invalid bearer token")
	ErrMissingToken  = securityv1.ErrorAuthErrorReasonBearerTokenMissing("bearer token missing")
	ErrInvalidAuth   = securityv1.ErrorAuthErrorReasonUnauthenticated("unauthenticated")
)

type securityCtx struct{}

func PolicyFromContext(ctx context.Context) security.Policy {
	if policy, ok := ctx.Value(securityCtx{}).(security.Policy); ok {
		return policy
	}
	return nil
}

func NewPolicyContext(ctx context.Context, claims security.Policy) context.Context {
	return context.WithValue(ctx, securityCtx{}, claims)
}

type claimCtx struct{}

func ClaimsFromContext(ctx context.Context) security.Claims {
	if claims, ok := ctx.Value(claimCtx{}).(security.Claims); ok {
		return claims
	}
	return nil
}

func NewClaimsContext(ctx context.Context, claims security.Claims) context.Context {
	return context.WithValue(ctx, claimCtx{}, claims)
}

type skipCtx struct{}

func NewSkipContext(ctx context.Context) context.Context {
	if SkipFromContext(ctx) {
		return ctx
	}
	return context.WithValue(ctx, skipCtx{}, true)
}

func SkipFromContext(ctx context.Context) bool {
	if _, ok := ctx.Value(skipCtx{}).(bool); ok {
		return true
	}
	return false
}

func mergePublic(public []string, paths ...string) []string {
	// Create a map to track unique paths
	pathMap := make(map[string]bool)

	// Add all existing paths to the map
	for _, p := range public {
		pathMap[p] = true
	}

	// Add all new paths to the map
	for _, path := range paths {
		pathMap[path] = true
	}

	// Convert the map back to a slice
	var uniquePaths []string
	for path := range pathMap {
		uniquePaths = append(uniquePaths, path)
	}

	return uniquePaths
}

func defaultSkipper(paths ...string) func(path string) bool {
	paths = mergePublic(paths)
	return func(path string) bool {
		pathLen := len(path)
		for _, p := range paths {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func IsSkipped(ctx context.Context, key string) bool {
	if SkipFromContext(ctx) {
		return true
	}
	if md, ok := metadata.FromServerContext(ctx); ok {
		if md.Get(key) == "true" || md.Get(key) == "1" {
			return true
		}
	}
	return false
}

func WithSkipContextClient(ctx context.Context, key string) context.Context {
	return metadata.AppendToClientContext(ctx, key, "true")
}

func WithSkipContextServer(ctx context.Context, key string) context.Context {
	md, ok := metadata.FromServerContext(ctx)
	if ok {
		md.Set(key, "true")
		return metadata.NewServerContext(NewSkipContext(ctx), md)
	}
	return ctx
}

func tokenParser(ctx context.Context, fns []func(ctx context.Context) string) string {
	for i := range fns {
		if s := fns[i](ctx); s != "" {
			return s
		}
	}
	return ""
}

func FromMetaData(key string) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return helper.FromMD(ctx, key)
	}
}

func defaultTokenParser(outer ...func(ctx context.Context) string) func(ctx context.Context) string {
	fns := []func(ctx context.Context) string{
		security.FromToken,
	}
	fns = append(fns, outer...)
	return func(ctx context.Context) string {
		return tokenParser(ctx, fns)
	}
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
