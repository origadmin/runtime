package security_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/origadmin/runtime/interfaces/security/declarative"
	impl "github.com/origadmin/runtime/security/declarative"
)

const (
	validToken   = "valid-token"
	invalidToken = "invalid-token"
	authHeader   = "Authorization"
)

// --- 1. Mocks and Test Components ---

// MockPrincipal implements declarative.Principal
type MockPrincipal struct {
	ID     string
	Roles  []string
	Claims map[string]*anypb.Any
}

func (mp *MockPrincipal) GetID() string {
	return mp.ID
}

func (mp *MockPrincipal) GetRoles() []string {
	return mp.Roles
}

func (mp *MockPrincipal) GetClaims() map[string]*anypb.Any {
	return mp.Claims
}

// mockAuthenticator implements the Authenticator interface for testing.
type mockAuthenticator struct {
	validToken string
}

func (m *mockAuthenticator) Authenticate(ctx context.Context, cred declarative.Credential) (declarative.Principal, error) {
	if cred.Raw() == m.validToken {
		return &MockPrincipal{ID: "test-user", Roles: []string{"user"}}, nil
	}
	return nil, errors.Unauthorized("UNAUTHORIZED", "invalid token")
}

func (m *mockAuthenticator) Supports(cred declarative.Credential) bool {
	return strings.ToLower(cred.Type()) == "bearer"
}

// --- 2. Middleware/Interceptor Implementation ---

// securityMiddleware is the core logic for our security middleware/interceptor.
func securityMiddleware(authenticators []declarative.Authenticator, extractor declarative.CredentialExtractor) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var provider declarative.ValueProvider

			// Adapt request to ValueProvider
			if tr, ok := transport.FromServerContext(ctx); ok {
				switch tr.Kind() {
				case transport.KindHTTP:
					if ht, ok := tr.(kratoshttp.Transporter); ok {
						provider = &impl.HTTPHeaderProvider{Header: ht.Request().Header}
					}
				case transport.KindGRPC:
					// No need to cast to gt, as we extract metadata from ctx directly
					md, _ := metadata.FromIncomingContext(ctx) // Extract metadata from the context directly
					provider = &impl.GRPCCredentialsProvider{MD: md}
				}
			} else if r, ok := req.(*http.Request); ok {
				// For standard net/http
				provider = &impl.HTTPHeaderProvider{Header: r.Header}
				ctx = r.Context()
			}

			if provider == nil {
				return nil, errors.InternalServer("UNKNOWN_REQUEST_TYPE", "request type not supported by security middleware")
			}

			// 1. Extract Credential
			cred, err := extractor.Extract(provider)
			if err != nil {
				return nil, errors.Unauthorized("CREDENTIAL_MISSING", err.Error())
			}

			// 2. Find a suitable Authenticator
			var authenticator declarative.Authenticator
			for _, auth := range authenticators {
				if auth.Supports(cred) {
					authenticator = auth
					break
				}
			}
			if authenticator == nil {
				return nil, errors.Unauthorized("AUTHENTICATOR_NOT_FOUND", "no authenticator found for credential type: "+cred.Type())
			}

			// 3. Authenticate
			principal, err := authenticator.Authenticate(ctx, cred)
			if err != nil {
				return nil, err // Return the error from the authenticator
			}

			// 4. Inject Principal into context
			newCtx := declarative.PrincipalWithContext(ctx, principal)

			// For standard net/http, update the request's context
			if r, ok := req.(*http.Request); ok {
				*r = *r.WithContext(newCtx)
				return handler(newCtx, r)
			}

			return handler(newCtx, req)
		}
	}
}

// --- 3. Test Scenarios ---

func TestMiddlewareIntegration(t *testing.T) {
	// --- Setup common components ---
	auth := &mockAuthenticator{validToken: validToken}
	extractor := impl.NewHeaderCredentialExtractor()
	securityChain := securityMiddleware([]declarative.Authenticator{auth}, extractor)

	// --- Scenario 1: Standard net/http Server ---
	t.Run("Standard net/http", func(t *testing.T) {
		// Handler that checks for principal
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := declarative.PrincipalFromContext(r.Context())
			assert.True(t, ok)
			assert.Equal(t, "test-user", p.GetID())
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		})

		// Wrap handler with middleware
		wrappedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The middleware expects to return the handler's result, which we ignore here.
			// We pass the request as the second argument as per our middleware's adaptation.
			_, err := securityChain(func(ctx context.Context, req interface{}) (interface{}, error) {
				handler.ServeHTTP(w, req.(*http.Request))
				return "ok", nil
			})(r.Context(), r)

			if err != nil {
				if se := errors.FromError(err); se != nil {
					http.Error(w, se.Message, int(se.Code))
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})

		server := httptest.NewServer(wrappedHandler)
		defer server.Close()

		// Test with valid token
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set(authHeader, "Bearer "+validToken)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Test with invalid token
		req, _ = http.NewRequest("GET", server.URL, nil)
		req.Header.Set(authHeader, "Bearer "+invalidToken)
		resp, err = server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	// --- Scenario 2: Kratos HTTP Server ---
	t.Run("Kratos HTTP", func(t *testing.T) {
		// Handler that checks for principal
		httpHandler := func(w http.ResponseWriter, r *http.Request) {
			p, ok := declarative.PrincipalFromContext(r.Context())
			assert.True(t, ok)
			assert.Equal(t, "test-user", p.GetID())
			w.WriteHeader(http.StatusOK)
		}

		srv := kratoshttp.NewServer(
			kratoshttp.Address("127.0.0.1:0"), // Use dynamic port
			kratoshttp.Middleware(securityChain),
		)
		srv.Handle("/test", http.HandlerFunc(httpHandler))
		// srv.HandlePrefix("/", router) // No longer needed with direct Handle

		// Start server in background
		go func() {
			if err := srv.Start(context.Background()); err != nil {
				panic(err)
			}
		}()
		defer srv.Stop(context.Background())

		// Wait for server to be ready by getting the endpoint
		var endpoint *url.URL
		var err error
		for i := 0; i < 10; i++ {
			endpoint, err = srv.Endpoint()
			if err == nil && endpoint != nil {
				break
			}
		}
		assert.NoError(t, err)
		assert.NotNil(t, endpoint)

		client := &http.Client{}

		// Test with valid token
		req, _ := http.NewRequest("GET", endpoint.String()+"/test", nil)
		req.Header.Set(authHeader, "Bearer "+validToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Test with invalid token
		req, _ = http.NewRequest("GET", endpoint.String()+"/test", nil)
		req.Header.Set(authHeader, "Bearer "+invalidToken)
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	// --- Scenario 3: Kratos gRPC Server ---
	t.Run("Kratos gRPC", func(t *testing.T) {
		// Simple gRPC handler
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			p, ok := declarative.PrincipalFromContext(ctx)
			assert.True(t, ok)
			assert.Equal(t, "test-user", p.GetID())
			return "OK", nil
		}

		// Create a listener on a dynamic port
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("failed to listen: %v", err)
		}
		addr := lis.Addr().String()
		lis.Close() // Close immediately, Kratos will re-open it

		srv := kratosgrpc.NewServer(
			kratosgrpc.Address(addr),
			kratosgrpc.Middleware(securityChain),
		)

		// Create a minimal service descriptor to register a method
		gsd := &grpc.ServiceDesc{
			ServiceName: "test.Service",
			HandlerType: (*interface{})(nil),
			Methods: []grpc.MethodDesc{
				{
					MethodName: "TestCall",
					Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
						if interceptor == nil {
							return handler(ctx, nil)
						}
						info := &grpc.UnaryServerInfo{
							Server:     srv,
							FullMethod: "/test.Service/TestCall",
						}
						return interceptor(ctx, nil, info, handler)
					},
				},
			},
		}
		srv.RegisterService(gsd, struct{}{})

		go func() {
			if err := srv.Start(context.Background()); err != nil {
				panic(err)
			}
		}()
		defer srv.Stop(context.Background())

		// Wait for server to be ready
		var conn *grpc.ClientConn
		for i := 0; i < 10; i++ {
			conn, err = grpc.Dial(addr, grpc.WithInsecure())
			if err == nil {
				break
			}
		}
		assert.NoError(t, err)
		defer conn.Close()

		// Test with valid token
		ctxValid := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(authHeader, "Bearer "+validToken))
		var reply string
		err = conn.Invoke(ctxValid, "/test.Service/TestCall", struct{}{}, &reply)
		assert.NoError(t, err)
		assert.Equal(t, "OK", reply)

		// Test with invalid token
		ctxInvalid := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(authHeader, "Bearer "+invalidToken))
		err = conn.Invoke(ctxInvalid, "/test.Service/TestCall", struct{}{}, &reply)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}
