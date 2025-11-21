package interfaces

import (
	"github.com/go-kratos/kratos/v2/transport"
)

// Encoder defines a function type for encoding a value into a byte slice.
type Encoder func(v any) ([]byte, error)

// GlobalDefaultKey is a constant string representing the global default key.
const GlobalDefaultKey = "default"

// Component is a generic runtime component.
type Component interface{}

// Server is the top-level abstraction for all service types within our framework.
// It is translated through the inline transport. Server, ensuring that any type that implements our Server interface,
// At the same time, it also automatically meets the requirements of Kratos App for transport. Server interface.
type Server interface {
	transport.Server
}

// Client is a tagged interface that represents an instance of a client connection, such as *grpc. ClientConn。
// Since clients with different protocols (e.g. gRPC, HTTP) do not have a unified interface,
// We use an empty interface to provide flexibility, and the caller needs to make type assertions.
type Client interface{}
