package middleware_test

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"strings"
	"testing"
	"time"

	authjwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	jwtv1 "github.com/origadmin/runtime/api/gen/go/config/middleware/jwt/v1"
	"github.com/origadmin/runtime/middleware"
	"github.com/origadmin/toolkits/errors"
)

// TestJWTMiddlewareCreation tests basic middleware creation
func TestJWTMiddlewareCreation(t *testing.T) {
	// Create JWT configuration
	jwtConfig := &jwtv1.JWT{
		Config: &jwtv1.AuthConfig{
			SigningMethod:       "HS256",
			SigningKey:          "test-secret-key",
			Issuer:              "test-issuer",
			Audience:            []string{"test-audience"},
			AccessTokenLifetime: 3600,
		},
		ClaimType: "registered",
	}

	t.Run("ServerMiddleware", func(t *testing.T) {
		serverMW, created := middleware.JwtServer(jwtConfig, nil)
		require.True(t, created, "Server middleware should be created")
		require.NotNil(t, serverMW, "Server middleware should not be nil")
	})

	t.Run("ClientMiddleware", func(t *testing.T) {
		clientMW, created := middleware.JwtClient(jwtConfig, nil)
		require.True(t, created, "Client middleware should be created")
		require.NotNil(t, clientMW, "Client middleware should not be nil")
	})
}

// 在文件顶部添加headerCarrier结构体定义
// 使用stdhttp.Header而不是http.Header
type headerCarrier stdhttp.Header

// 实现transport.Header接口的所有方法
func (h headerCarrier) Get(key string) string {
	return stdhttp.Header(h).Get(key)
}

func (h headerCarrier) Set(key, value string) {
	stdhttp.Header(h).Set(key, value)
}

func (h headerCarrier) Add(key, value string) {
	stdhttp.Header(h).Add(key, value)
}

func (h headerCarrier) Values(key string) []string {
	return stdhttp.Header(h).Values(key)
}

func (h headerCarrier) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

// 定义mockTransport结构体实现transport.Transporter接口
type mockTransport struct {
	kind        string
	endpoint    string
	operation   string
	reqHeader   headerCarrier
	replyHeader headerCarrier
}

// 实现transport.Transporter接口的所有方法
func (tr *mockTransport) Kind() transport.Kind {
	return transport.Kind(tr.kind)
}

func (tr *mockTransport) Endpoint() string {
	return tr.endpoint
}

func (tr *mockTransport) Operation() string {
	return tr.operation
}

func (tr *mockTransport) RequestHeader() transport.Header {
	return tr.reqHeader
}

func (tr *mockTransport) ReplyHeader() transport.Header {
	return tr.replyHeader
}

// TestJWTTokenGeneration tests client-side token generation
func TestJWTTokenGeneration(t *testing.T) {
	jwtConfig := &jwtv1.JWT{
		Config: &jwtv1.AuthConfig{
			SigningMethod:       "HS256",
			SigningKey:          "test-secret-key",
			Issuer:              "test-issuer",
			Audience:            []string{"test-audience"},
			AccessTokenLifetime: 3600,
		},
		ClaimType: "registered",
	}

	// Create client middleware with a custom subject factory
	opts := &middleware.Options{
		SubjectFactory: func() string { return "test-user" },
	}

	t.Run("GenerateToken", func(t *testing.T) {
		clientMW, created := middleware.JwtClient(jwtConfig, opts)
		require.True(t, created, "Client middleware should be created")

		t.Run("WithValidContext", func(t *testing.T) {
			// Create a context with required claims
			now := time.Now()
			claims := &jwt.RegisteredClaims{
				Subject:   "test-user",
				Issuer:    "test-issuer",
				Audience:  []string{"test-audience"},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			}
			// Use authjwt.NewContext to set the claims in the context
			ctx := context.Background()
			header := headerCarrier{}
			tr := &mockTransport{
				kind:        "test",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   header,
				replyHeader: headerCarrier{},
			}
			ctx = transport.NewClientContext(ctx, tr)
			ctx = authjwt.NewContext(ctx, claims)

			var token string
			genHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				if md, ok := transport.FromClientContext(ctx); ok {
					if auth := md.RequestHeader().Get("Authorization"); auth != "" {
						token = strings.TrimPrefix(auth, "Bearer ")
						return token, nil
					}
				}
				return nil, errors.New("no token generated")
			}

			genMiddleware := clientMW(genHandler)
			tokenSource, err := genMiddleware(ctx, nil)
			assert.NoError(t, err, "Token generation should succeed")
			if v, ok := tokenSource.(string); ok {
				token = v
			} else {
				assert.Fail(t, "Token generation failed")
			}
			assert.NotEmpty(t, token, "Token should be generated")
		})

		t.Run("WithEmptyContext", func(t *testing.T) {
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			}

			middlewareFunc := clientMW(handler)
			_, err := middlewareFunc(context.Background(), nil)
			assert.Error(t, err, "Should fail with empty context")
		})
	})
}

// TestJWTTokenValidation tests server-side token validation
func TestJWTTokenValidation(t *testing.T) {
	jwtConfig := &jwtv1.JWT{
		Config: &jwtv1.AuthConfig{
			SigningMethod:       "HS256",
			SigningKey:          "test-secret-key",
			Issuer:              "test-issuer",
			Audience:            []string{"test-audience"},
			AccessTokenLifetime: 1, // 1 second for testing
		},
		ClaimType: "registered",
	}

	t.Run("ValidateTokens", func(t *testing.T) {
		// Create server middleware
		serverMW, created := middleware.JwtServer(jwtConfig, nil)
		require.True(t, created, "Server middleware should be created")

		// Create client middleware with a custom subject factory
		clientOpts := &middleware.Options{
			SubjectFactory: func() string { return "test-user" },
		}
		clientMW, created := middleware.JwtClient(jwtConfig, clientOpts)
		require.True(t, created, "Client middleware should be created")

		t.Run("ValidToken", func(t *testing.T) {
			// Generate a token using the client middleware
			var token string
			genHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				if tr, ok := transport.FromClientContext(ctx); ok {
					if auth := tr.RequestHeader().Get("Authorization"); auth != "" {
						token = strings.TrimPrefix(auth, "Bearer ")
						return token, nil
					}
				}
				return nil, fmt.Errorf("no token generated")
			}

			// Create context with claims for token generation
			claims := &jwt.RegisteredClaims{
				Subject:  "test-user",
				Issuer:   "test-issuer",
				Audience: []string{"test-audience"},
			}
			header := make(headerCarrier)
			header.Set("Authorization", "Bearer "+token)
			ctx := transport.NewClientContext(context.Background(), &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   header,
				replyHeader: make(headerCarrier),
			})
			ctx = authjwt.NewContext(ctx, claims)
			// Generate token
			genMiddleware := clientMW(genHandler)
			tokenSource, err := genMiddleware(ctx, nil)
			require.NoError(t, err, "Token generation should succeed")
			if v, ok := tokenSource.(string); ok {
				token = v
			} else {
				t.Fatal("Token generation failed")
			}
			require.NotEmpty(t, token, "Token should be generated")

			ctx = transport.NewServerContext(context.Background(), &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   header,
				replyHeader: make(headerCarrier),
			})
			ctx = authjwt.NewContext(ctx, claims)
			// Now validate the token
			validateHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				// Extract claims from context using authjwt.FromContext
				claims, ok := authjwt.FromContext(ctx)
				if !ok {
					t.Fatal("Failed to extract claims from context")
				}
				registeredClaims, ok := claims.(*jwt.RegisteredClaims)
				if !ok {
					t.Fatal("Failed to convert claims to RegisteredClaims")
				}
				assert.Equal(t, "test-user", registeredClaims.Subject)
				// Verify the claims are in the context
				assert.Equal(t, "test-user", registeredClaims.Subject, "Subject should match")
				assert.Equal(t, "test-issuer", registeredClaims.Issuer, "Issuer should match")
				assert.Contains(t, registeredClaims.Audience, "test-audience", "Audience should contain test-audience")
				return "validated", nil
			}

			serverCtx := transport.NewServerContext(context.Background(), &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   header,
				replyHeader: make(headerCarrier),
			})
			// Validate token
			validateMiddleware := serverMW(validateHandler)
			result, err := validateMiddleware(serverCtx, nil)

			assert.NoError(t, err, "Token validation should succeed")
			assert.Equal(t, "validated", result, "Should return validated result")
		})

		t.Run("ExpiredToken", func(t *testing.T) {
			// Create an expired token
			claims := &jwt.RegisteredClaims{
				Subject:   "test-user",
				Issuer:    "test-issuer",
				Audience:  []string{"test-audience"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString([]byte("test-secret-key"))
			require.NoError(t, err, "Failed to sign token")

			// Try to validate the expired token
			validateHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "should not reach here", nil
			}
			header := make(headerCarrier)
			header.Set("Authorization", "Bearer "+signedToken)
			serverCtx := transport.NewServerContext(context.Background(), &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   header,
				replyHeader: make(headerCarrier),
			})

			validateMiddleware := serverMW(validateHandler)
			_, err = validateMiddleware(serverCtx, nil)

			assert.Error(t, err, "Should reject expired token")
			assert.Contains(t, err.Error(), "token has expired", "Error should indicate token is expired")
		})

		t.Run("InvalidSignature", func(t *testing.T) {
			// Create a token with invalid signature
			claims := &jwt.RegisteredClaims{
				Subject:   "test-user",
				Issuer:    "test-issuer",
				Audience:  []string{"test-audience"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedToken, err := token.SignedString([]byte("wrong-secret-key"))
			require.NoError(t, err, "Failed to sign token with wrong key")

			// Try to validate the token with invalid signature
			validateHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "should not reach here", nil
			}

			md := metadata.Pairs("authorization", "Bearer "+signedToken)
			serverCtx := metadata.NewIncomingContext(context.Background(), md)

			validateMiddleware := serverMW(validateHandler)
			_, err = validateMiddleware(serverCtx, nil)

			assert.Error(t, err, "Should reject token with invalid signature")
		})

		t.Run("NoToken", func(t *testing.T) {
			// Try to validate without a token
			validateHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "should not reach here", nil
			}

			// Start with some initial metadata to avoid empty map issue
			md := metadata.Pairs("test-key", "test-value")
			serverCtx := metadata.NewIncomingContext(context.Background(), md)

			validateMiddleware := serverMW(validateHandler)
			_, err := validateMiddleware(serverCtx, nil)

			assert.Error(t, err, "Should reject request without token")
		})
	})
}

// TestJWTWithConfigFile tests using configuration from YAML file
func TestJWTWithConfigFile(t *testing.T) {
	t.Run("UsingYAMLConfig", func(t *testing.T) {
		// Create JWT config matching the YAML file
		jwtConfig := &jwtv1.JWT{
			Config: &jwtv1.AuthConfig{
				SigningMethod:       "HS256",
				SigningKey:          "test-secret-key-for-hs256",
				Issuer:              "test-issuer",
				Audience:            []string{"test-audience"},
				AccessTokenLifetime: 3600,
			},
			ClaimType: "registered",
		}

		// Test that we can create middleware with this config
		serverMW, created := middleware.JwtServer(jwtConfig, nil)
		require.True(t, created, "Server middleware should be created with YAML config")
		require.NotNil(t, serverMW, "Server middleware should not be nil")

		clientMW, created := middleware.JwtClient(jwtConfig, nil)
		require.True(t, created, "Client middleware should be created with YAML config")
		require.NotNil(t, clientMW, "Client middleware should not be nil")
	})
}

// TestJWTClientServerFlow tests the complete flow: client generates token, server validates it
func TestJWTClientServerFlow(t *testing.T) {
	// Create JWT configuration
	jwtConfig := &jwtv1.JWT{
		Config: &jwtv1.AuthConfig{
			SigningMethod:       "HS256",
			SigningKey:          "test-secret-key",
			Issuer:              "test-issuer",
			Audience:            []string{"test-audience"},
			AccessTokenLifetime: 3600,
		},
		ClaimType: "registered",
	}

	t.Run("CompleteFlow", func(t *testing.T) {
		// Create server middleware for token validation
		serverMW, created := middleware.JwtServer(jwtConfig, nil)
		require.True(t, created, "Server middleware should be created")
		require.NotNil(t, serverMW, "Server middleware should not be nil")

		// Create client middleware with a custom subject factory
		clientOpts := &middleware.Options{
			SubjectFactory: func() string { return "test-user" },
		}
		clientMW, created := middleware.JwtClient(jwtConfig, clientOpts)
		require.True(t, created, "Client middleware should be created")
		require.NotNil(t, clientMW, "Client middleware should not be nil")

		// Test token generation and validation in sequence
		t.Run("GenerateThenValidate", func(t *testing.T) {
			// Step 1: Generate token using client middleware
			var generatedToken string
			tokenGenHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				// Extract token from transport context instead of metadata
				if tr, ok := transport.FromClientContext(ctx); ok {
					auth := tr.RequestHeader().Get("Authorization")
					if auth != "" {
						generatedToken = strings.TrimPrefix(auth, "Bearer ")
						return generatedToken, nil
					}
				}
				return nil, fmt.Errorf("no token generated")
			}

			// Create context with claims for token generation
			claims := &jwt.RegisteredClaims{
				Subject:  "test-user",
				Issuer:   "test-issuer",
				Audience: []string{"test-audience"},
			}

			// Create a client transport context instead of just authjwt context
			header := make(headerCarrier)
			ctx := authjwt.NewContext(context.Background(), claims)
			clientCtx := transport.NewClientContext(ctx, &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   header,
				replyHeader: make(headerCarrier),
			})

			// Apply client middleware
			clientMiddlewareFunc := clientMW(tokenGenHandler)
			token, err := clientMiddlewareFunc(clientCtx, nil)
			require.NoError(t, err, "Token generation should succeed")
			require.NotEmpty(t, token, "Token should be generated")
			generatedToken = token.(string)

			// Step 2: Validate token using server middleware
			protectedHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				// Extract claims from context using authjwt.FromContext
				claims, ok := authjwt.FromContext(ctx)
				if !ok {
					t.Fatal("Failed to extract claims from context")
				}
				registeredClaims, ok := claims.(*jwt.RegisteredClaims)
				if !ok {
					t.Fatal("Failed to convert claims to RegisteredClaims")
				}
				assert.Equal(t, "test-user", registeredClaims.Subject)
				// Verify the claims are in the context
				assert.Equal(t, "test-user", registeredClaims.Subject, "Subject should match")
				assert.Equal(t, "test-issuer", registeredClaims.Issuer, "Issuer should match")
				assert.Contains(t, registeredClaims.Audience, "test-audience", "Audience should contain test-audience")
				return "protected-resource", nil
			}

			// Create transport server context with token instead of metadata context
			serverHeader := make(headerCarrier)
			serverHeader.Set("Authorization", "Bearer "+generatedToken)
			serverCtx := transport.NewServerContext(context.Background(), &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   serverHeader,
				replyHeader: make(headerCarrier),
			})

			// Apply server middleware
			serverMiddlewareFunc := serverMW(protectedHandler)
			result, err := serverMiddlewareFunc(serverCtx, nil)

			assert.NoError(t, err, "Generated token should be validated successfully")
			assert.Equal(t, "protected-resource", result, "Should return protected resource")
		})

		// Test with YAML config
		t.Run("WithYAMLConfig", func(t *testing.T) {
			// Load YAML config
			yamlConfig := &jwtv1.JWT{
				Config: &jwtv1.AuthConfig{
					SigningMethod:       "HS256",
					SigningKey:          "test-secret-key-for-hs256",
					Issuer:              "test-issuer",
					Audience:            []string{"test-audience"},
					AccessTokenLifetime: 3600,
				},
				ClaimType: "registered",
			}

			// Create server middleware with YAML config
			serverMW, created := middleware.JwtServer(yamlConfig, nil)
			require.True(t, created, "Server middleware should be created with YAML config")

			// Create client middleware with YAML config
			clientMW, created := middleware.JwtClient(yamlConfig, clientOpts)
			require.True(t, created, "Client middleware should be created with YAML config")

			// Test token generation and validation
			var generatedToken string
			tokenGenHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				// Extract token from transport context instead of metadata
				if tr, ok := transport.FromClientContext(ctx); ok {
					auth := tr.RequestHeader().Get("Authorization")
					if auth != "" {
						generatedToken = strings.TrimPrefix(auth, "Bearer ")
						return generatedToken, nil
					}
				}
				return nil, fmt.Errorf("no token generated")
			}

			// Create context with claims for token generation
			claims := &jwt.RegisteredClaims{
				Subject:  "test-user",
				Issuer:   "test-issuer",
				Audience: []string{"test-audience"},
			}

			// Create transport client context
			header := make(headerCarrier)
			ctx := authjwt.NewContext(context.Background(), claims)
			clientCtx := transport.NewClientContext(ctx, &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   header,
				replyHeader: make(headerCarrier),
			})

			// Generate token
			clientMiddlewareFunc := clientMW(tokenGenHandler)
			token, err := clientMiddlewareFunc(clientCtx, nil)
			require.NoError(t, err, "Token generation with YAML config should succeed")
			require.NotEmpty(t, token, "Token should be generated with YAML config")
			generatedToken = token.(string)

			// Validate token
			validateHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
				// If we get here, the token was validated successfully
				return "validated-with-yaml", nil
			}

			// Create transport server context with token
			serverHeader := make(headerCarrier)
			serverHeader.Set("Authorization", "Bearer "+generatedToken)
			serverCtx := transport.NewServerContext(context.Background(), &mockTransport{
				kind:        "http",
				endpoint:    "test",
				operation:   "test",
				reqHeader:   serverHeader,
				replyHeader: make(headerCarrier),
			})

			validateMiddleware := serverMW(validateHandler)
			result, err := validateMiddleware(serverCtx, nil)

			assert.NoError(t, err, "Token validation with YAML config should succeed")
			assert.Equal(t, "validated-with-yaml", result, "Should return validated result with YAML config")
		})
	})
}
