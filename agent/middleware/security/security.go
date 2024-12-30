/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package security

import (
	"github.com/go-kratos/kratos/v2/metadata"

	"github.com/origadmin/runtime/context"
	securityv1 "github.com/origadmin/runtime/gen/go/security/v1"
	"github.com/origadmin/toolkits/security"
)

const (
	reason          string = "FORBIDDEN"
	StringBoolTrue         = "true"
	StringBoolFalse        = "false"
)

var (
	ErrInvalidToken          = securityv1.ErrorSecurityErrorReasonBearerTokenMissing("bearer token missing")
	ErrInvalidClaims         = securityv1.ErrorSecurityErrorReasonInvalidClaims("invalid bearer token")
	ErrMissingClaims         = securityv1.ErrorSecurityErrorReasonInvalidClaims("missing scheme")
	ErrTokenNotFound         = securityv1.ErrorSecurityErrorReasonTokenNotFound("token not found")
	ErrMissingToken          = securityv1.ErrorSecurityErrorReasonBearerTokenMissing("bearer token missing")
	ErrInvalidAuthentication = securityv1.ErrorSecurityErrorReasonInvalidAuthentication("unauthenticated")
	ErrInvalidAuthorization  = securityv1.ErrorSecurityErrorReasonInvalidAuthorization("unauthorized")
)

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

func aggregateTokenParsers(outer ...func(ctx context.Context) string) func(ctx context.Context) string {
	fns := []func(ctx context.Context) string{
		security.TokenFromContext,
	}
	for i := range outer {
		if outer[i] == nil {
			continue
		}
		fns = append(fns, outer[i])
	}
	return func(ctx context.Context) string {
		return tokenParser(ctx, fns)
	}
}
func FromMetaDataKey(ctx context.Context, key string) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(key)
	}
	return ""
}

func FromMetaData(key string) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		return FromMetaDataKey(ctx, key)
	}
}
