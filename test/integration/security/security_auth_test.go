package security_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/stretchr/testify/assert"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb" // For StringValue

	basicv1 "github.com/origadmin/runtime/api/gen/go/config/security/authn/basic/v1"
	securityv1 "github.com/origadmin/runtime/api/gen/go/config/security/v1"
	"github.com/origadmin/runtime/interfaces/security/declarative"
)

// --- Mock Implementations for Testing ---

// MockPrincipal implements declarative.Principal
type MockPrincipal struct {
	ID     string
	Roles  []string
	Claims map[string]*anypb.Any
}

func (mp *MockPrincipal) GetID() string                    { return mp.ID }
func (mp *MockPrincipal) GetRoles() []string               { return mp.Roles }
func (mp *MockPrincipal) GetClaims() map[string]*anypb.Any { return mp.Claims }

// MockCredential implements declarative.Credential
type MockCredential struct {
	RawCredential  string
	CredentialType string
	// Claims now directly holds the claims map, not a single "payload" key
	Claims  map[string]*anypb.Any
	Headers http.Header
}

func (mc *MockCredential) Type() string {
	return mc.CredentialType
}

func (mc *MockCredential) Raw() string {
	return mc.RawCredential
}

func (mc *MockCredential) ParsedPayload(message proto.Message) error {
	// For this mock, we assume the "payload" key in Claims holds the main payload.
	if payloadAny, ok := mc.Claims["payload"]; ok {
		return payloadAny.UnmarshalTo(message)
	}
	return fmt.Errorf("payload claim not found or cannot be unmarshaled")
}

func (mc *MockCredential) Get(key string) (string, bool) {
	val := mc.Headers.Get(key)
	return val, val != ""
}

func (mc *MockCredential) GetAll() map[string]any {
	allClaims := make(map[string]any)
	for k, v := range mc.Claims {
		allClaims[k] = v
	}
	return allClaims
}

// MockCredentialParser implements declarative.CredentialParser
type MockCredentialParser struct {
	ExpectedRawCredential string
	ReturnCredential      declarative.Credential
	ReturnError           error
}

func (mcp *MockCredentialParser) ParseCredential(ctx context.Context, rawCredential string) (declarative.Credential, error) {
	if mcp.ReturnError != nil {
		return nil, mcp.ReturnError
	}
	if rawCredential == mcp.ExpectedRawCredential {
		return mcp.ReturnCredential, nil
	}
	return nil, errors.Unauthorized("INVALID_CREDENTIAL", "credential mismatch")
}

func (mcp *MockCredentialParser) ParseCredentialFrom(ctx context.Context, sourceType string, source any) (declarative.Credential, error) {
	return nil, fmt.Errorf("ParseCredentialFrom not implemented in mock")
}

// --- JwtAuthenticator (Simplified for Test) ---

// jwtAuthenticator implements declarative.Authenticator
type jwtAuthenticator struct {
	parser declarative.CredentialParser
}

// Authenticate now correctly uses cred.ParsedPayload to extract Principal info
func (a *jwtAuthenticator) Authenticate(ctx context.Context, cred declarative.Credential) (declarative.Principal, error) {
	// Use the generated securityv1.Payload to unmarshal the credential's payload.
	payload := &securityv1.Payload{}

	// Use ParsedPayload to get the structured data
	if err := cred.ParsedPayload(payload); err != nil {
		// Format the error message before passing it to errors.Unauthorized
		errMsg := fmt.Sprintf("failed to parse credential payload: %v", err)
		return nil, errors.Unauthorized("UNAUTHORIZED", errMsg)
	}

	// Now, construct the Principal from the parsed payload, checking the optional fields.
	principalID := "default_user"
	roles := []string{"guest"}
	claims := make(map[string]*anypb.Any) // Collect all claims from the payload's specific types

	// Check for specific credential types in the payload
	if token := payload.GetToken(); token != nil {
		principalID = token.GetAccessToken() // Using AccessToken as a placeholder for ID
		roles = []string{"user"}             // Default role for token users
		// You might also extract other claims from the token if it were a full JWT object
	} else if basic := payload.GetBasic(); basic != nil {
		principalID = basic.GetUsername()
		roles = []string{"basic_user"}
	} else if key := payload.GetKey(); key != nil {
		principalID = key.GetKey() // Using key as ID
		roles = []string{"api_key_user"}
	} else if oidc := payload.GetOidc(); oidc != nil {
		// Handle OIDC specific fields if needed
		principalID = oidc.GetIdToken() // Placeholder
		roles = []string{"oidc_user"}
	} else if rawData := payload.GetRawData(); rawData != "" {
		// Handle raw data if needed
		principalID = "raw_data_user" // Placeholder
		roles = []string{"raw_data_user"}
	} else {
		// No specific credential data found, use defaults
	}

	// If a more specific ID was not found, try to get it from a generic claim if available
	// NOTE: This part might be redundant if all identity info is expected in the payload.
	// Keeping it for now as a fallback.
	if principalID == "default_user" {
		if idAny, ok := cred.GetAll()["user_id"].(*anypb.Any); ok {
			var idStr wrapperspb.StringValue // Use StringValue to unmarshal string from Any
			if err := idAny.UnmarshalTo(&idStr); err == nil {
				principalID = idStr.GetValue()
			}
		}
	}

	// Populate claims from credential's meta or other sources if available
	// For this mock, we'll just use the claims from the payload if any were extracted.
	// In a real scenario, you'd populate claims more comprehensively.

	return &MockPrincipal{ID: principalID, Roles: roles, Claims: claims}, nil
}

// --- Integration Test ---

func TestSecurityAuthFlow(t *testing.T) {
	ctx := context.Background()

	// Scenario 1: Successful Authentication with TokenCredential
	t.Run("Successful Authentication with TokenCredential", func(t *testing.T) {
		expectedRawCredential := "valid.jwt.credential"
		tokenCred := &securityv1.TokenCredential{
			AccessToken:  "user123_access_token",
			RefreshToken: "user123_refresh_token",
			ExpiresIn:    3600,
			TokenType:    "Bearer",
		}
		payloadProto := &securityv1.Payload{
			Token: tokenCred, // Directly set the Token field
		}
		payloadAny, err := anypb.New(payloadProto)
		if err != nil {
			t.Fatalf("Failed to create anypb.Any for payloadProto: %v", err)
		}

		headerClaimAny, err := anypb.New(wrapperspb.String("header_val"))
		if err != nil {
			t.Fatalf("Failed to create anypb.Any for header_val: %v", err)
		}

		mockClaims := map[string]*anypb.Any{
			"payload":      payloadAny, // The main payload is wrapped in an Any
			"header_claim": headerClaimAny,
		}

		mockParsedCredential := &MockCredential{
			RawCredential:  expectedRawCredential,
			CredentialType: "jwt",
			Claims:         mockClaims,
			Headers:        http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", expectedRawCredential)}},
		}
		expectedPrincipal := &MockPrincipal{
			ID:     tokenCred.GetAccessToken(), // Expecting ID from AccessToken
			Roles:  []string{"user"},           // Default role for token users
			Claims: map[string]*anypb.Any{},    // Claims are not directly populated in this mock
		}

		authenticator := &jwtAuthenticator{}

		principal, err := authenticator.Authenticate(ctx, mockParsedCredential)

		assert.NoError(t, err)
		assert.NotNil(t, principal)
		assert.Equal(t, expectedPrincipal.GetID(), principal.GetID())
		assert.Equal(t, expectedPrincipal.GetRoles(), principal.GetRoles())
		// assert.Equal(t, expectedPrincipal.GetClaims()["custom_claim"].GetValue(), principal.GetClaims()["custom_claim"].GetValue()) // No custom claims in this mock
	})

	// Scenario 2: Authenticator with CredentialParser (simulating full flow with TokenCredential)
	t.Run("Authenticator with CredentialParser (TokenCredential)", func(t *testing.T) {
		expectedRawCredential := "valid.jwt.credential.from.parser"
		tokenCred := &securityv1.TokenCredential{
			AccessToken:  "parser_user_access_token",
			RefreshToken: "parser_user_refresh_token",
			ExpiresIn:    7200,
			TokenType:    "Bearer",
		}
		payloadProto := &securityv1.Payload{
			Token: tokenCred, // Directly set the Token field
		}
		payloadAny, err := anypb.New(payloadProto)

		if err != nil {
			t.Fatalf("Failed to create anypb.Any for payloadProto: %v", err)
		}

		sourceClaimAny, err := anypb.New(wrapperspb.String("from_parser"))
		if err != nil {
			t.Fatalf("Failed to create anypb.Any for source_claim: %v", err)
		}

		mockClaims := map[string]*anypb.Any{
			"payload":      payloadAny,
			"source_claim": sourceClaimAny,
		}

		mockParsedCredential := &MockCredential{
			RawCredential:  expectedRawCredential,
			CredentialType: "jwt",
			Claims:         mockClaims,
			Headers:        http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", expectedRawCredential)}},
		}

		mockParser := &MockCredentialParser{
			ExpectedRawCredential: expectedRawCredential,
			ReturnCredential:      mockParsedCredential,
			ReturnError:           nil,
		}

		authenticator := &jwtAuthenticator{parser: mockParser}

		reqHeader := make(http.Header)
		reqHeader.Set("Authorization", fmt.Sprintf("Bearer %s", expectedRawCredential))
		rawFromRequest := strings.TrimPrefix(reqHeader.Get("Authorization"), "Bearer ")

		parsedCred, err := authenticator.parser.ParseCredential(ctx, rawFromRequest)
		assert.NoError(t, err)
		assert.NotNil(t, parsedCred)

		principal, err := authenticator.Authenticate(ctx, parsedCred)

		assert.NoError(t, err)
		assert.NotNil(t, principal)
		assert.Equal(t, tokenCred.GetAccessToken(), principal.GetID())

		assert.Contains(t, principal.GetRoles(), "user") // Default role for token users
		// assert.Equal(t, payloadProto.Claims["source_claim"].GetValue(), principal.GetClaims()["source_claim"].GetValue()) // Claims not directly populated in this mock
	})

	// Scenario 3: CredentialParser Returns Error
	t.Run("CredentialParser Returns Error", func(t *testing.T) {
		expectedRawCredential := "invalid.jwt.credential"
		mockParser := &MockCredentialParser{
			ExpectedRawCredential: "some.other.credential",

			ReturnCredential: nil,
			ReturnError:      errors.Unauthorized("BAD_SIGNATURE", "invalid credential signature"),
		}

		authenticator := &jwtAuthenticator{parser: mockParser}

		reqHeader := make(http.Header)
		reqHeader.Set("Authorization", fmt.Sprintf("Bearer %s", expectedRawCredential))
		rawFromRequest := strings.TrimPrefix(reqHeader.Get("Authorization"), "Bearer ")

		parsedCred, err := authenticator.parser.ParseCredential(ctx, rawFromRequest)
		assert.Error(t, err)
		assert.Nil(t, parsedCred)
		assert.True(t, errors.IsUnauthorized(err))
		assert.Contains(t, err.Error(), "invalid credential signature")
	})

	// Scenario 4: Authenticator with no "payload" claim in Credential
	t.Run("Authenticator with no payload claim", func(t *testing.T) {
		mockParsedCredential := &MockCredential{
			RawCredential:  "credential.without.payload",
			CredentialType: "jwt",
			Claims:         map[string]*anypb.Any{}, // No "payload" claim
		}

		authenticator := &jwtAuthenticator{}

		principal, err := authenticator.Authenticate(ctx, mockParsedCredential)

		assert.Error(t, err) // Expecting error because ParsedPayload will fail
		assert.Nil(t, principal)
		assert.True(t, errors.IsUnauthorized(err))
		assert.Contains(t, err.Error(), "failed to parse credential payload")
	})

	// Scenario 5: Authenticator with payload but no specific credential data (e.g., no Token, Basic, Key)
	t.Run("Authenticator with empty payload data", func(t *testing.T) {
		payloadProto := &securityv1.Payload{
			// No specific credential data set
		}
		payloadAny, err := anypb.New(payloadProto)
		if err != nil {
			t.Fatalf("Failed to create anypb.Any for empty payloadProto: %v", err)
		}

		mockParsedCredential := &MockCredential{
			RawCredential:  "credential.with.empty.payload_data",
			CredentialType: "jwt",
			Claims:         map[string]*anypb.Any{"payload": payloadAny},
		}

		authenticator := &jwtAuthenticator{}

		principal, err := authenticator.Authenticate(ctx, mockParsedCredential)

		assert.NoError(t, err) // Should not error, but use defaults
		assert.NotNil(t, principal)
		assert.Equal(t, "default_user", principal.GetID())
		assert.Equal(t, []string{"guest"}, principal.GetRoles())
	})

	// Scenario 6: Successful Authentication with BasicCredential
	t.Run("Successful Authentication with BasicCredential", func(t *testing.T) {
		expectedRawCredential := "valid.basic.credential"
		basicCred := &basicv1.BasicCredential{ // Corrected type to securityv1.BasicCredential
			Username: "basic_user_id",
			Password: "basic_password",
		}
		payloadProto := &securityv1.Payload{
			Basic: basicCred, // Directly set the Basic field
		}
		payloadAny, err := anypb.New(payloadProto)
		if err != nil {
			t.Fatalf("Failed to create anypb.Any for basic payloadProto: %v", err)
		}

		mockClaims := map[string]*anypb.Any{
			"payload": payloadAny,
		}

		mockParsedCredential := &MockCredential{
			RawCredential:  expectedRawCredential,
			CredentialType: "basic",
			Claims:         mockClaims,
			Headers:        http.Header{"Authorization": []string{fmt.Sprintf("Basic %s", expectedRawCredential)}},
		}
		expectedPrincipal := &MockPrincipal{
			ID:     basicCred.GetUsername(),
			Roles:  []string{"basic_user"},
			Claims: map[string]*anypb.Any{},
		}

		authenticator := &jwtAuthenticator{} // Authenticator is generic, not JWT specific

		principal, err := authenticator.Authenticate(ctx, mockParsedCredential)

		assert.NoError(t, err)
		assert.NotNil(t, principal)
		assert.Equal(t, expectedPrincipal.GetID(), principal.GetID())
		assert.Equal(t, expectedPrincipal.GetRoles(), principal.GetRoles())
	})
}
