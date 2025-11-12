package security_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/security/declarative"
	impl "github.com/origadmin/runtime/security/declarative"
)

const (
	separatedValidToken   = "separated-valid-token"
	separatedInvalidToken = "separated-invalid-token"
	separatedAuthHeader   = "Authorization"
	separatedTestResource = "/test"
	separatedTestAction   = "read"
)

// --- 1. Mocks and Test Components for Separated Middleware ---

// MockPrincipal implements declarative.Principal
type SeparatedMockPrincipal struct {
	ID     string
	Roles  []string
	Claims map[string]*anypb.Any
}

func (mp *SeparatedMockPrincipal) GetID() string {
	return mp.ID
}

func (mp *SeparatedMockPrincipal) GetRoles() []string {
	return mp.Roles
}

func (mp *SeparatedMockPrincipal) GetClaims() map[string]any {
	// For this test, we don't need to return actual claims, just satisfy the interface.
	return nil
}

// mockAuthenticator implements the Authenticator interface for testing.
type SeparatedMockAuthenticator struct {
	validToken string
}

func (m *SeparatedMockAuthenticator) Authenticate(ctx context.Context, cred declarative.Credential) (declarative.Principal, error) {
	var token securityv1.BearerCredential
	if err := cred.ParsedPayload(&token); err != nil {
		return nil, errors.Unauthorized("UNAUTHORIZED", "invalid token format or type mismatch")
	}
	if token.Token != m.validToken {
		return nil, errors.Unauthorized("UNAUTHORIZED", "invalid token")
	}
	return &SeparatedMockPrincipal{ID: "separated-test-user", Roles: []string{"user"}}, nil
}

func (m *SeparatedMockAuthenticator) Supports(cred declarative.Credential) bool {
	return strings.ToLower(cred.Type()) == "jwt"
}

// mockAuthorizer implements the Authorizer interface for testing.
type SeparatedMockAuthorizer struct {
	// For simplicity, this mock always authorizes if a principal is present.
	// In a real scenario, it would check principal roles, resource, action against policies.
}

func (ma *SeparatedMockAuthorizer) Authorize(ctx context.Context, p declarative.Principal, resourceIdentifier string, action string) (bool, error) {
	// For this test, we'll just check if the principal exists.
	// A real authorizer would perform more complex logic.
	if p == nil || p.GetID() == "" {
		return false, errors.Forbidden("FORBIDDEN", "principal not found or invalid")
	}
	// Simulate a successful authorization for any valid principal
	return true, nil
}

// --- 2. Separated Middleware/Interceptor Implementations ---

// kratosAuthnMiddleware is the core logic for our authentication middleware/interceptor for Kratos.
func kratosAuthnMiddleware(authenticator declarative.Authenticator, extractor declarative.CredentialExtractor) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			provider, err := impl.FromServerContext(ctx)
			if err != nil {
				return nil, err
			}

			if provider == nil {
				return nil, errors.InternalServer("UNKNOWN_REQUEST_TYPE", "request type not supported by authentication middleware")
			}

			// 1. Extract Credential
			cred, err := extractor.Extract(ctx, provider) // Pass ctx here
			if err != nil {
				return nil, errors.Unauthorized("CREDENTIAL_MISSING", err.Error())
			}

			// 2. Authenticate
			principal, err := authenticator.Authenticate(ctx, cred)
			if err != nil {
				return nil, err // Return the error from the authenticator
			}

			// 3. Inject Principal into context
			newCtx := declarative.PrincipalWithContext(ctx, principal)

			return handler(newCtx, req)
		}
	}
}

// kratosAuthzMiddleware is the core logic for our authorization middleware/interceptor for Kratos.
func kratosAuthzMiddleware(authorizer declarative.Authorizer, resourceIdentifier string, action string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 1. Get Principal from context (assumes authentication middleware ran before)
			principal, ok := declarative.PrincipalFromContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHENTICATED", "principal not found in context, authentication required")
			}

			// 2. Authorize
			allowed, err := authorizer.Authorize(ctx, principal, resourceIdentifier, action)
			if err != nil {
				return nil, err // Return the error from the authorizer
			}
			if !allowed {
				return nil, errors.Forbidden("PERMISSION_DENIED", "principal not authorized for this resource/action")
			}

			return handler(ctx, req)
		}
	}
}

// standardAuthnMiddleware is the core logic for our authentication middleware/interceptor for standard net/http.
func standardAuthnMiddleware(authenticator declarative.Authenticator, extractor declarative.CredentialExtractor) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var provider declarative.ValueProvider

			if r, ok := req.(*http.Request); ok {
				provider = impl.FromHTTPRequest(r)
			} else {
				return nil, errors.InternalServer("UNEXPECTED_REQUEST_TYPE", "standardAuthnMiddleware expects *http.Request")
			}

			if provider == nil {
				return nil, errors.InternalServer("UNKNOWN_REQUEST_TYPE", "request type not supported by authentication middleware")
			}

			// 1. Extract Credential
			cred, err := extractor.Extract(ctx, provider) // Pass ctx here
			if err != nil {
				return nil, errors.Unauthorized("CREDENTIAL_MISSING", err.Error())
			}

			// 2. Authenticate
			principal, err := authenticator.Authenticate(ctx, cred)
			if err != nil {
				return nil, err
			}

			// 3. Inject Principal into context
			newCtx := declarative.PrincipalWithContext(ctx, principal)

			if r, ok := req.(*http.Request); ok {
				*r = *r.WithContext(newCtx)
				return handler(newCtx, r)
			}
			return nil, errors.InternalServer("UNEXPECTED_STATE", "standardAuthnMiddleware failed to update http.Request")
		}
	}
}

// standardAuthzMiddleware is the core logic for our authorization middleware/interceptor for standard net/http.
func standardAuthzMiddleware(authorizer declarative.Authorizer, resourceIdentifier string, action string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 1. Get Principal from context (assumes authentication middleware ran before)
			principal, ok := declarative.PrincipalFromContext(ctx)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHENTICATED", "principal not found in context, authentication required")
			}

			// 2. Authorize
			allowed, err := authorizer.Authorize(ctx, principal, resourceIdentifier, action)
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, errors.Forbidden("PERMISSION_DENIED", "principal not authorized for this resource/action")
			}

			return handler(ctx, req)
		}
	}
}

// authnFilterFunc is a standard net/http middleware (FilterFunc) that performs authentication
// and injects the Principal into the http.Request's context.
func authnFilterFunc(authenticator declarative.Authenticator, extractor declarative.CredentialExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			provider := impl.FromHTTPRequest(r)

			cred, err := extractor.Extract(r.Context(), provider) // Pass r.Context() here
			if err != nil {
				http.Error(w, errors.Unauthorized("CREDENTIAL_MISSING", err.Error()).Error(), http.StatusUnauthorized)
				return
			}

			principal, err := authenticator.Authenticate(r.Context(), cred)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			newCtx := declarative.PrincipalWithContext(r.Context(), principal)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

// authzFilterFunc is a standard net/http middleware (FilterFunc) that performs authorization.
func authzFilterFunc(authorizer declarative.Authorizer, resourceIdentifier string, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := declarative.PrincipalFromContext(r.Context())
			if !ok {
				http.Error(w, errors.Unauthorized("UNAUTHENTICATED", "principal not found in context, authentication required").Error(), http.StatusUnauthorized)
				return
			}

			allowed, err := authorizer.Authorize(r.Context(), principal, resourceIdentifier, action)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden) // Use Forbidden for authorization errors
				return
			}
			if !allowed {
				http.Error(w, errors.Forbidden("PERMISSION_DENIED", "principal not authorized for this resource/action").Error(), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// TestServiceHTTP defines a simple HTTP service for Kratos.
type SeparatedTestServiceHTTP interface {
	TestCall(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

// testServiceHTTPImpl implements TestServiceHTTP.
type SeparatedTestServiceHTTPImpl struct {
	t *testing.T
}

func (s *SeparatedTestServiceHTTPImpl) TestCall(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	p, ok := declarative.PrincipalFromContext(ctx)
	assert.True(s.t, ok, "Principal should be in Kratos HTTP context (native handler)")
	if ok {
		assert.Equal(s.t, "separated-test-user", p.GetID())
	}
	return &emptypb.Empty{}, nil
}

func RegisterSeparatedTestServiceHTTPServer(s *kratoshttp.Server, srv SeparatedTestServiceHTTP) {
	r := s.Route("/")
	r.GET(separatedTestResource, separatedTestServiceCallHTTPHandler(srv))
}

func separatedTestServiceCallHTTPHandler(srv SeparatedTestServiceHTTP) func(ctx kratoshttp.Context) error {
	return func(ctx kratoshttp.Context) error {
		var in emptypb.Empty
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.TestCall(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

// --- 3. Test Scenarios for Separated Middleware ---

func TestSeparatedMiddlewareIntegration(t *testing.T) {
	// --- Setup common components ---
	auth := &SeparatedMockAuthenticator{validToken: separatedValidToken}
	authorizer := &SeparatedMockAuthorizer{} // Our mock authorizer
	extractor := impl.NewHeaderCredentialExtractor()

	// Use the CompositeAuthenticator
	compositeAuth := impl.NewCompositeAuthenticator(auth)

	// Create independent authentication and authorization middleware chains
	kratosAuthnChain := kratosAuthnMiddleware(compositeAuth, extractor)
	kratosAuthzChain := kratosAuthzMiddleware(authorizer, separatedTestResource, separatedTestAction)

	standardAuthnChain := standardAuthnMiddleware(compositeAuth, extractor)
	standardAuthzChain := standardAuthzMiddleware(authorizer, separatedTestResource, separatedTestAction)

	// --- Scenario 1: Standard net/http Server ---
	t.Run("Standard net/http - Separated Middleware", func(t *testing.T) {
		finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p, ok := declarative.PrincipalFromContext(r.Context())
			assert.True(t, ok, "Principal should be in context")
			if ok {
				assert.Equal(t, "separated-test-user", p.GetID())
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		})

		middlewareAdapter := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				kratosHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
					next.ServeHTTP(w, req.(*http.Request))
					return "ok", nil
				}

				// Chain authn and authz middleware
				chainedHandler := standardAuthzChain(standardAuthnChain(kratosHandler))

				_, err := chainedHandler(r.Context(), r)

				if err != nil {
					if se := errors.FromError(err); se != nil {
						http.Error(w, se.Message, int(se.Code))
					} else {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			})
		}

		server := httptest.NewServer(middlewareAdapter(finalHandler))
		defer server.Close()

		// Test with valid token (Authn & Authz pass)
		req, _ := http.NewRequest("GET", server.URL, nil)
		req.Header.Set(separatedAuthHeader, "Bearer "+separatedValidToken)
		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Test with invalid token (Authn fails)
		req, _ = http.NewRequest("GET", server.URL, nil)
		req.Header.Set(separatedAuthHeader, "Bearer "+separatedInvalidToken)
		resp, err = server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()

		// Test with no token (Authn fails)
		req, _ = http.NewRequest("GET", server.URL, nil)
		resp, err = server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	// --- Scenario 2: Kratos HTTP Server ---
	t.Run("Kratos HTTP - Separated Middleware", func(t *testing.T) {
		srv := kratoshttp.NewServer(
			kratoshttp.Address("127.0.0.1:0"),
			kratoshttp.Middleware(kratosAuthnChain, kratosAuthzChain), // Chain authn and authz middleware
		)

		testService := &SeparatedTestServiceHTTPImpl{t: t}
		RegisterSeparatedTestServiceHTTPServer(srv, testService)

		go func() {
			if err := srv.Start(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}()
		defer srv.Stop(context.Background())
		time.Sleep(100 * time.Millisecond)

		endpoint, err := srv.Endpoint()
		assert.NoError(t, err)
		assert.NotNil(t, endpoint)

		client := &http.Client{}

		// Test with valid token (Authn & Authz pass)
		req, _ := http.NewRequest("GET", endpoint.String()+separatedTestResource, nil)
		req.Header.Set(separatedAuthHeader, "Bearer "+separatedValidToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Test with invalid token (Authn fails)
		req, _ = http.NewRequest("GET", endpoint.String()+separatedTestResource, nil)
		req.Header.Set(separatedAuthHeader, "Bearer "+separatedInvalidToken)
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()

		// Test with no token (Authn fails)
		req, _ = http.NewRequest("GET", endpoint.String()+separatedTestResource, nil)
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	// --- Scenario 3: Kratos HTTP Server with http.HandlerFunc Context Propagation ---
	t.Run("Kratos HTTP - HandlerFunc Context Propagation - Separated Middleware", func(t *testing.T) {
		httpHandler := func(w http.ResponseWriter, r *http.Request) {
			p, ok := declarative.PrincipalFromContext(r.Context())
			assert.True(t, ok, "Principal should be in Kratos HTTP context (HandlerFunc) via FilterFunc")
			if ok {
				assert.Equal(t, "separated-test-user", p.GetID())
			}
			w.WriteHeader(http.StatusOK)
		}

		// Chain authn and authz filters
		chainedFilter := authnFilterFunc(compositeAuth, extractor)(authzFilterFunc(authorizer, separatedTestResource, separatedTestAction)(http.HandlerFunc(httpHandler)))

		srv := kratoshttp.NewServer(
			kratoshttp.Address("127.0.0.1:0"),
		)
		srv.Handle(separatedTestResource, chainedFilter)

		go func() {
			if err := srv.Start(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}()
		defer srv.Stop(context.Background())
		time.Sleep(100 * time.Millisecond)

		endpoint, err := srv.Endpoint()
		assert.NoError(t, err)
		assert.NotNil(t, endpoint)

		client := &http.Client{}

		// Test with valid token (Authn & Authz pass)
		req, _ := http.NewRequest("GET", endpoint.String()+separatedTestResource, nil)
		req.Header.Set(separatedAuthHeader, "Bearer "+separatedValidToken)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()

		// Test with invalid token (Authn fails)
		req, _ = http.NewRequest("GET", endpoint.String()+separatedTestResource, nil)
		req.Header.Set(separatedAuthHeader, "Bearer "+separatedInvalidToken)
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()

		// Test with no token (Authn fails)
		req, _ = http.NewRequest("GET", endpoint.String()+separatedTestResource, nil)
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		resp.Body.Close()
	})

	// --- Scenario 4: Kratos gRPC Server ---
	t.Run("Kratos gRPC - Separated Middleware", func(t *testing.T) {
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			p, ok := declarative.PrincipalFromContext(ctx)
			assert.True(t, ok, "Principal should be in Kratos gRPC context")
			if ok {
				assert.Equal(t, "separated-test-user", p.GetID())
			}
			return &emptypb.Empty{}, nil
		}

		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("failed to listen: %v", err)
		}
		addr := lis.Addr().String()
		lis.Close()

		srv := kratosgrpc.NewServer(
			kratosgrpc.Address(addr),
			kratosgrpc.Middleware(kratosAuthnChain, kratosAuthzChain), // Chain authn and authz middleware
		)

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
			if err := srv.Start(context.Background()); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
				panic(err)
			}
		}()
		defer srv.Stop(context.Background())
		time.Sleep(100 * time.Millisecond)

		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		assert.NoError(t, err)
		defer conn.Close()

		// Test with valid token (Authn & Authz pass)
		ctxValid := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(separatedAuthHeader, "Bearer "+separatedValidToken))
		var reply emptypb.Empty
		err = conn.Invoke(ctxValid, "/test.Service/TestCall", &emptypb.Empty{}, &reply)
		assert.NoError(t, err)

		// Test with invalid token (Authn fails)
		ctxInvalid := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(separatedAuthHeader, "Bearer "+separatedInvalidToken))
		err = conn.Invoke(ctxInvalid, "/test.Service/TestCall", &emptypb.Empty{}, &reply)
		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())

		// Test with no token (Authn fails)
		ctxNoToken := context.Background()
		err = conn.Invoke(ctxNoToken, "/test.Service/TestCall", &emptypb.Empty{}, &reply)
		assert.Error(t, err)
		st, ok = status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}
