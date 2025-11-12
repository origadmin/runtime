// Package declarative provides a generic way to access key-value pairs from a request source,
// such as HTTP headers or gRPC metadata. This avoids using 'any' and ensures type safety.
package declarative

// ValueProvider provides a generic way to access key-value pairs from a request source,
// such as HTTP headers or gRPC metadata. This avoids using 'any' and ensures type safety.
type ValueProvider interface {
	// Values returns the values associated with the given key.
	// It returns a slice of strings because sources like HTTP headers can have
	// multiple values for the same key.
	Values(key string) []string
	// Get returns the first value associated with the given key.
	// If the key is not found, it returns an empty string.
	Get(key string) string
	// GetAll returns all key-value pairs from the source.
	GetAll() map[string][]string
}
