package interfaces

import (
	"context"
)

// ContextOptions is a foundational struct for functional options.
// It includes a context.Context for passing common, context-bound values.
// Other packages can embed this struct to inherit common context handling.
type ContextOptions struct {
	Context context.Context
}
