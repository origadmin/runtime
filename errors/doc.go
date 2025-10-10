/*
Package errors provides enhanced error handling with module support, error codes, and metadata.

This package extends the standard error handling with additional features and is divided into two main parts:

1. API interface unified error handling (integration with Kratos errors)
2. Internal module error creation and handling

The package also provides utilities for working with error metadata and compatibility with the standard errors package.

Basic Usage:

	// Create a new structured error for internal module use
	err := errors.New("auth", "INVALID_TOKEN", "invalid authentication token")

	// Add operation context
	err = err.WithCaller() // Automatically gets the caller function name
	// or
	err = err.WithOperation("CheckToken") // Manually set operation name

	// Add metadata
	err = err.WithField("user_id", 123)

	// Wrap an existing error
	err = errors.Wrap(err, "db", "QUERY_FAILED", "failed to query user")

	// Convert to Kratos error for API responses
	kratosErr := errors.ToKratos(err, commonv1.ErrorReason_UNAUTHENTICATED)

	// Check error type
	if errors.Is(err, &errors.Structured{Module: "auth", Code: "INVALID_TOKEN"}) {
		// Handle invalid token
	}

	// Working with error metadata
	if val, ok := errors.LookupMeta(kratosErr, "user_id"); ok {
		// Use the metadata
	}
*/
package errors