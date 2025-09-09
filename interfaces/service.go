package interfaces

import (
	"github.com/go-kratos/kratos/v2/transport"
)

// Server is the top-level abstraction for all service types within our framework.
// It is translated through the inline transport. Server, ensuring that any type that implements our Server interface,
// At the same time, it also automatically meets the requirements of Kratos App for transport. Server interface.
type Server interface {
	transport.Server // <-- Core: Ensure full compatibility with Kratos App
}

// Client is the top-level abstraction for all service types within our framework.
// It is translated through the inline transport. Client, ensuring that any type that implements our Client interface,
// At the same time, it also automatically meets the requirements of Kratos App for transport. Client interface.
type Client interface {
}
