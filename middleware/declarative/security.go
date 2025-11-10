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
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	middlewarev1 "github.com/origadmin/runtime/api/gen/go/config/middleware/v1"
	"github.com/origadmin/runtime/extension/customize"
	iface "github.com/origadmin/runtime/interfaces/security/declarative"
	internal "github.com/origadmin/runtime/internal/security/declarative"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/middleware"
)

var (
	// globalPolicyManager holds the singleton instance of the PolicyManager.
	// It must be set by the application's main entry point before the runtime starts.
	globalPolicyManager PolicyManager
)

// SetPolicyManager sets the global policy manager for the declarative security middleware.
// This function must be called during the application's initialization phase.
func SetPolicyManager(pm PolicyManager) {
	if globalPolicyManager != nil {
		log.Warn("global policy manager is being overwritten")
	}
	globalPolicyManager = pm
}

func init() {
	// Register the security middleware factory.
	// This allows the runtime to create the security middleware by name ("security").
	middleware.Register(middleware.Security, &securityFactory{})
}

// securityFactory is the factory for creating the declarative security middleware.
type securityFactory struct{}

// NewMiddlewareServer creates a new server-side security middleware.
func (f *securityFactory) NewMiddlewareServer(cfg *middlewarev1.Middleware, opts ...middleware.Option) (kratosmiddleware.Middleware, bool) {
	if globalPolicyManager == nil {
		log.Error("declarative security middleware factory failed: PolicyManager has not been set. Call declarative.SetPolicyManager() during initialization.")
		return nil, false
	}

	// Unmarshal the middleware-specific configuration from the 'customize' field.
	mwCfg := &Config{}
	if cfg.Customize != nil {
		if err := customize.UnmarshalTo(cfg.Customize, mwCfg); err != nil {
			log.Errorf("failed to unmarshal security middleware config: %v", err)
			return nil, false
		}
	}

	// Create the middleware options.
	securityOpts := []Option{
		WithPolicyManager(globalPolicyManager),
	}
	if mwCfg.DefaultPolicy != "" {
		securityOpts = append(securityOpts, WithDefaultPolicy(mwCfg.DefaultPolicy))
	}

	return SecurityMiddleware(securityOpts...), true
}

// NewMiddlewareClient creates a new client-side security middleware.
// This is not applicable for the declarative security middleware.
func (f *securityFactory) NewMiddlewareClient(cfg *middlewarev1.Middleware, opts ...middleware.Option) (kratosmiddleware.Middleware, bool) {
	return nil, false
}

// Config defines the configuration for the declarative security middleware.
type Config struct {
	DefaultPolicy string `json:"defaultPolicy"`
}

// PolicyManager defines the interface for managing and retrieving security policies.
type PolicyManager interface {
	// GetPolicy retrieves a SecurityPolicy by its name.
	GetPolicy(name string) (iface.SecurityPolicy, error)
}

// Option is a function that configures the SecurityMiddleware.
type Option func(*options)

type options struct {
	policyManager PolicyManager
	defaultPolicy string
}

// WithPolicyManager sets the PolicyManager for the middleware.
func WithPolicyManager(pm PolicyManager) Option {
	return func(o *options) {
		o.policyManager = pm
	}
}

// WithDefaultPolicy sets the default policy name to use if no policy is specified in the metadata.
func WithDefaultPolicy(policyName string) Option {
	return func(o *options) {
		o.defaultPolicy = policyName
	}
}

// SecurityMiddleware creates a Kratos middleware for declarative security.
// It uses policy names from route metadata to apply authentication and authorization.
func SecurityMiddleware(opts ...Option) kratosmiddleware.Middleware {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.policyManager == nil {
		log.Error("declarative security middleware created without a PolicyManager")
		// Return a "failing" middleware if not configured correctly.
		return func(handler kratosmiddleware.Handler) kratosmiddleware.Handler {
			return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
				return nil, errors.New(500, "MIDDLEWARE_MISCONFIGURED", "security middleware is not configured with a policy manager")
			}
		}
	}

	return func(handler kratosmiddleware.Handler) kratosmiddleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.New(500, "TRANSPORT_CONTEXT_MISSING", "transport context is missing")
			}

			// 1. Get policy name from metadata
			policyName := getPolicyNameFromMetadata(tr)

			// Handle "public" policy (skip authentication and authorization)
			if policyName == "public" {
				log.Debugf("Policy 'public' detected for method %s, skipping security checks.", tr.Operation())
				return handler(ctx, req)
			}

			// Use default policy if not specified
			if policyName == "" && o.defaultPolicy != "" {
				policyName = o.defaultPolicy
				log.Debugf("No policy specified for method %s, using default policy '%s'.", tr.Operation(), policyName)
			}

			if policyName == "" {
				return nil, errors.New(401, "UNAUTHORIZED", "no security policy specified for this method")
			}

			// 2. Get SecurityPolicy instance from PolicyManager
			policy, err := o.policyManager.GetPolicy(policyName)
			if err != nil {
				log.Errorf("Failed to get security policy '%s': %v", policyName, err)
				return nil, errors.New(500, "SECURITY_POLICY_ERROR", fmt.Sprintf("failed to get security policy '%s'", policyName))
			}

			// 3. Authenticate
			credentialSource := internal.NewKratosTransportHeaderAdapter(tr.RequestHeader())
			principal, authErr := policy.Authenticate(ctx, credentialSource)
			if authErr != nil {
				log.Debugf("Authentication failed for method %s with policy '%s': %v", tr.Operation(), policyName, authErr)
				return nil, errors.Unauthorized("UNAUTHENTICATED", "authentication failed")
			}

			// 4. Authorize
			fullMethodName := tr.Operation() // Kratos operation usually is the full method name
			authorized, authzErr := policy.Authorize(ctx, principal, fullMethodName)
			if authzErr != nil {
				log.Errorf("Authorization check failed for method %s with policy '%s': %v", tr.Operation(), policyName, authzErr)
				return nil, errors.New(500, "AUTHORIZATION_ERROR", fmt.Sprintf("authorization check failed: %v", authzErr))
			}
			if !authorized {
				log.Debugf("Authorization denied for principal %s on method %s with policy '%s'.", principal.GetID(), tr.Operation(), policyName)
				return nil, errors.Forbidden("FORBIDDEN", "authorization denied")
			}

			// 5. Inject Principal into context and proceed
			ctx = iface.NewContextWithPrincipal(ctx, principal)
			return handler(ctx, req)
		}
	}
}

// getPolicyNameFromMetadata extracts the security policy name from Kratos transport metadata.
func getPolicyNameFromMetadata(tr transport.Transporter) string {
	// For HTTP, metadata is usually in the header. For gRPC, it's in grpc.metadata.
	// Kratos's transport.Transporter.Operation() often contains the full method name,
	// and the policy name would be attached to the route's metadata during code generation.
	// This is a placeholder; the actual mechanism to retrieve the policy name from metadata
	// might depend on how 'origen' or 'protoc-gen-go-http' injects it.
	// For now, we assume it's directly accessible or can be derived from the operation.

	// A more robust solution would involve Kratos's router context or a custom metadata key.
	// For example, if the policy name is stored in a context value.
	// For the purpose of this design, we assume it's directly available via some Kratos mechanism
	// or can be retrieved from the operation's options.

	// The design document states: "中间件从 Kratos 的 transport.Context 中获取路由信息，
	// 并找到之前代码生成时注入的元数据，得到策略名"
	// This implies a mechanism to get route-specific metadata.
	// Kratos's `transport.Transporter` has `Operation()`, which is the full method name.
	// We need a way to map `Operation()` to the policy name.
	// This mapping is typically done by the router itself, or by a custom metadata injector.

	// For the sake of this example, let's assume the policy name is stored in a custom HTTP header
	// or gRPC metadata key.
	// In a real scenario, `origen` would generate code that makes this accessible.

	// If using HTTP, check a custom header
	if ht, ok := tr.(*http.Transport); ok {
		if policy := ht.Request().Header.Get("X-Security-Policy"); policy != "" {
			return policy
		}
	}

	// Fallback or default logic if not found in custom header
	// For now, let's return an empty string if not explicitly set.
	return ""
}
