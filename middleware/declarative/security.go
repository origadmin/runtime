/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package declarative provides a Kratos middleware for declarative security.
package declarative

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	kratosmiddleware "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport" // Keep import for transport.KindHTTP/GRPC if needed elsewhere, but tr.Operation() is removed

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/interfaces/options"
	iface "github.com/origadmin/runtime/interfaces/security/declarative"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware"
	securityImpl "github.com/origadmin/runtime/security/declarative" // Renamed alias
)

// factory is the factory for creating the declarative security middleware.
type factory struct{}

func init() {
	middleware.Register(middleware.DeclarativeSecurity, &factory{})
}

// NewMiddlewareServer creates a new server-side security middleware.
func (f *factory) NewMiddlewareServer(cfg *middlewarev1.Middleware,
	opts ...options.Option) (kratosmiddleware.Middleware, bool) {
	o := FromOptions(opts) // Use FromOptions from this package

	if cfg.GetSecurity() == nil {
		return nil, false
	}

	if o.policyProvider == nil {
		return nil, false
	}
	// Check if CredentialExtractor is provided
	if o.credentialExtractor == nil {
		return nil, false
	}

	// Unmarshal the middleware-specific configuration from the 'customize' field.
	o.defaultPolicy = "public"

	if cfg.GetSecurity().GetDefaultPolicy() != "" {
		o.defaultPolicy = cfg.GetSecurity().GetDefaultPolicy()
	}

	return SecurityMiddleware(o), true
}

// NewMiddlewareClient creates a new client-side security middleware.
// This is not applicable for the declarative security middleware.
func (f *factory) NewMiddlewareClient(_cfg *middlewarev1.Middleware, _opts ...options.Option) (kratosmiddleware.
Middleware, bool) {
	return nil, false
}

// SecurityMiddleware creates a Kratos middleware for declarative security.
// It uses policy names from route metadata to apply authentication and authorization.
func SecurityMiddleware(o *Options) kratosmiddleware.Middleware {
	return func(handler kratosmiddleware.Handler) kratosmiddleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// Get ValueProvider from the context using the utility function
			valueProvider, vpErr := securityImpl.FromServerContext(ctx)
			if vpErr != nil {
				// If we can't get a ValueProvider, it means transport context is missing or invalid.
				// We can't proceed with security checks without it.
				log.Errorf("Failed to get ValueProvider from context: %v", vpErr)
				return nil, errors.New(500, "VALUE_PROVIDER_ERROR", "failed to get value provider from context")
			}

			// Extract fullMethodName from the ValueProvider
			fullMethodName := valueProvider.Get("fullMethodName")
			if fullMethodName == "" {
				log.Errorf("fullMethodName not found in ValueProvider")
				return nil, errors.New(500, "SECURITY_POLICY_ERROR", "fullMethodName not found in request context")
			}

			// 1. Get policy name from metadata
			var policyName string
			// The full gRPC method name is the canonical identifier for the policy.
			// For HTTP requests, the router might store the matched template in the transport.
			// We prioritize the gRPC method name but could fall back to HTTP path if needed.
			policyName, err = o.policyProvider.GetPolicyNameForMethod(ctx, fullMethodName)
			if err != nil {
				log.Errorf("Failed to get policy name for method %s: %v", fullMethodName, err)
				return nil, errors.New(500, "SECURITY_POLICY_ERROR", "failed to determine security policy")
			}

			// Handle "public" policy (skip authentication and authorization)
			if policyName == "public" {
				log.Debugf("Policy 'public' detected for method %s, skipping security checks.", fullMethodName)
				return handler(ctx, req)
			}

			// Use default policy if not specified
			if policyName == "" && o.defaultPolicy != "" {
				policyName = o.defaultPolicy
				log.Debugf("No policy specified for method %s, using default policy '%s'.", fullMethodName, policyName)
			}

			if policyName == "" {
				return nil, errors.New(401, "UNAUTHORIZED", "no security policy specified for this method")
			}

			// 2. Get SecurityPolicy instance from PolicyManager
			policy, err := o.policyProvider.GetPolicy(ctx, policyName)
			if err != nil {
				log.Errorf("Failed to get security policy '%s': %v", policyName, err)
				return nil, errors.New(500, "SECURITY_POLICY_ERROR", fmt.Sprintf("failed to get security policy '%s'", policyName))
			}

			// 3. Authenticate
			// Extract credential using the provided CredentialExtractor
			cred, extractErr := o.credentialExtractor.Extract(ctx, valueProvider)
			if extractErr != nil {
				log.Debugf("Authentication failed for method %s with policy '%s': %v", fullMethodName, policyName, extractErr)
				return nil, errors.Unauthorized("UNAUTHENTICATED", "failed to extract credential")
			}

			principal, authErr := policy.Authenticate(ctx, cred) // Use extracted credential
			if authErr != nil {
				log.Debugf("Authentication failed for method %s with policy '%s': %v", fullMethodName, policyName, authErr)
				return nil, errors.Unauthorized("UNAUTHENTICATED", "authentication failed")
			}

			// 4. Authorize
			// For now, using a placeholder "access" action. This should ideally be derived from the request.
			authorized, authzErr := policy.Authorize(ctx, principal, fullMethodName, "access") // Added action parameter
			if authzErr != nil {
				log.Errorf("Authorization check failed for method %s with policy '%s': %v", fullMethodName, policyName, authzErr)
				return nil, errors.New(500, "AUTHORIZATION_ERROR", fmt.Sprintf("authorization check failed: %v", authzErr))
			}
			if !authorized {
				log.Debugf("Authorization denied for principal %s on method %s with policy '%s'.", principal.GetID(), fullMethodName, policyName)
				return nil, errors.Forbidden("FORBIDDEN", "authorization denied")
			}

			// 5. Inject Principal into context and proceed
			ctx = iface.PrincipalWithContext(ctx, principal) // Corrected function name
			return handler(ctx, req)
		}
	}
}
