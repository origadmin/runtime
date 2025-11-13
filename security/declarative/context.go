package declarative

import (
	"context"
	"fmt"
	"net/http"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc/metadata"

	"github.com/origadmin/runtime/interfaces/security/declarative" // This import path remains the same
)

// valueSource defines the common interface for extracting key-value pairs from different sources.
type valueSource interface {
	Get(key string) string
	Values(key string) []string
	GetAll() map[string][]string
}

// HTTPHeaderProvider adapts an http.Header to the valueSource interface.
type HTTPHeaderProvider struct {
	Header http.Header
}

// Values returns the values associated with the given key from the HTTP header.
func (p HTTPHeaderProvider) Values(key string) []string {
	return p.Header[key]
}

// Get returns the first value associated with the given key from the HTTP header.
func (p HTTPHeaderProvider) Get(key string) string {
	return p.Header.Get(key)
}

// GetAll returns all key-value pairs from the HTTP header.
func (p HTTPHeaderProvider) GetAll() map[string][]string {
	return p.Header
}

// GRPCCredentialsProvider adapts a grpc.Metadata to the valueSource interface.
type GRPCCredentialsProvider struct {
	MD metadata.MD
}

// Values returns the values associated with the given key from the gRPC metadata.
func (p GRPCCredentialsProvider) Values(key string) []string {
	return p.MD.Get(key)
}

// Get returns the first value associated with the given key from the gRPC metadata.
func (p GRPCCredentialsProvider) Get(key string) string {
	values := p.MD.Get(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// GetAll returns all key-value pairs from the gRPC metadata.
func (p GRPCCredentialsProvider) GetAll() map[string][]string {
	return p.MD
}

// contextAwareSecurityRequest implements declarative.SecurityRequest by wrapping a valueSource
// and providing context-specific information like Kind, Operation, Method, and RouteTemplate.
type contextAwareSecurityRequest struct {
	kind          string
	operation     string
	method        string
	routeTemplate string
	delegate      valueSource
}

// Kind returns the type of the request as a string (e.g., "grpc", "http").
func (c *contextAwareSecurityRequest) Kind() string {
	return c.kind
}

// GetOperation returns the primary identifier for the logical operation being performed.
func (c *contextAwareSecurityRequest) GetOperation() string {
	return c.operation
}

// GetMethod returns the HTTP verb (e.g., "GET", "POST") if the request is an HTTP call.
func (c *contextAwareSecurityRequest) GetMethod() string {
	return c.method
}

// GetRouteTemplate returns the matched HTTP route template (e.g., "/v1/users/{id}")
// if the request is an HTTP call and a route template was matched.
func (c *contextAwareSecurityRequest) GetRouteTemplate() string {
	return c.routeTemplate
}

// Get returns the first value associated with the given key from the delegate valueSource.
func (c *contextAwareSecurityRequest) Get(key string) string {
	if c.delegate != nil {
		return c.delegate.Get(key)
	}
	return ""
}

// Values returns the values associated with the given key from the delegate valueSource.
func (c *contextAwareSecurityRequest) Values(key string) []string {
	if c.delegate != nil {
		return c.delegate.Values(key)
	}
	return nil
}

// GetAll returns all key-value pairs from the delegate valueSource.
func (c *contextAwareSecurityRequest) GetAll() map[string][]string {
	if c.delegate != nil {
		return c.delegate.GetAll()
	}
	return nil
}

// FromHTTPRequest creates an HTTPHeaderProvider from http.Request.
func FromHTTPRequest(r *http.Request) *HTTPHeaderProvider {
	return &HTTPHeaderProvider{Header: r.Header}
}

// FromGRPCMetadata creates a GRPCCredentialsProvider from metadata.MD.
func FromGRPCMetadata(md metadata.MD) *GRPCCredentialsProvider {
	return &GRPCCredentialsProvider{MD: md}
}

// FromServerContext extracts a SecurityRequest from the server context.
func FromServerContext(ctx context.Context) (declarative.SecurityRequest, error) {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, kratoserrors.New(500, "TRANSPORT_CONTEXT_MISSING", "transport context is missing")
	}

	var (
		kind          = string(tr.Kind())
		operation     = tr.Operation()
		method        string
		routeTemplate string
		delegate      valueSource
		err           error
	)

	switch tr.Kind() {
	case transport.KindHTTP:
		if ht, ok := tr.(kratoshttp.Transporter); ok {
			req := ht.Request()
			delegate = FromHTTPRequest(req)
			method = req.Method
			// Kratos HTTP transport does not directly expose the matched route template
			// via the transport.Transporter interface or kratoshttp.Transporter.
			// If route template is needed, it would typically be stored in the request context
			// by a routing middleware. For now, it remains empty.
		} else {
			err = kratoserrors.New(500, "INVALID_HTTP_TRANSPORT", "invalid HTTP transport type")
		}
	case transport.KindGRPC:
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			delegate = FromGRPCMetadata(md)
			// For gRPC, method and routeTemplate are typically derived from the operation
			// or are not directly applicable in the same way as HTTP.
		} else {
			err = kratoserrors.New(400, "NO_METADATA", "no metadata found in context")
		}
	default:
		err = fmt.Errorf("unsupported transport type: %v", tr.Kind())
	}

	if err != nil {
		return nil, err
	}

	return &contextAwareSecurityRequest{
		kind:          kind,
		operation:     operation,
		method:        method,
		routeTemplate: routeTemplate,
		delegate:      delegate,
	}, nil
}
