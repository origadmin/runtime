// Package declarative defines the contract for extracting a credential from a request
// by using a ValueProvider to access request data in a transport-agnostic way.
package declarative

// CredentialExtractor defines the contract for extracting a credential from a request
// and wrapping it into a standardized Credential object.
type CredentialExtractor interface {
	// Extract finds a credential within the provided source, wraps it in a standard
	// Credential object, and returns it.
	Extract(provider ValueProvider) (Credential, error)
}
