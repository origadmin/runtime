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
	"github.com/origadmin/runtime/middleware/security/internal/helper"
)

const (
	reason          string = "FORBIDDEN"
	StringBoolTrue         = "true"
	StringBoolFalse        = "false"
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
	if md, ok := metadata.FromServerContext(ctx); ok && md.Get(key) == StringBoolTrue {
		return true
	}
	return false
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
