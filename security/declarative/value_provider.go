package declarative

import (
	"net/http"

	"google.golang.org/grpc/metadata"
)

// HTTPHeaderProvider adapts an http.Header to the ValueProvider interface.
type HTTPHeaderProvider struct {
	Header http.Header
}

// Get retrieves values from the http.Header.
func (p HTTPHeaderProvider) Get(key string) []string {
	return p.Header.Values(key)
}

// GRPCCredentialsProvider adapts a grpc.Metadata to the ValueProvider interface.
type GRPCCredentialsProvider struct {
	MD metadata.MD
}

// Get retrieves values from the grpc.Metadata.
func (p GRPCCredentialsProvider) Get(key string) []string {
	return p.MD.Get(key)
}
