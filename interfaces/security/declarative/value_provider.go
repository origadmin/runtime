// Package declarative provides a generic way to access key-value pairs from a request source,
// such as HTTP headers or gRPC metadata. This avoids using 'any' and ensures type safety.
package declarative

// ValueProvider provides a generic way to access key-value pairs from a request source,
// such as HTTP headers or gRPC metadata. This avoids using 'any' and ensures type safety.
type ValueProvider interface {
	// Get returns the values associated with the given key.
	// It returns a slice of strings because sources like HTTP headers can have
	// multiple values for the same key.
	Get(key string) []string
}
