package service

import (
	"github.com/origadmin/framework/runtime/interfaces"
)

// Options contains the options for creating service components.
// It embeds interfaces.ContextOptions for common context handling.
type Options struct {
	interfaces.ContextOptions // Anonymous embedding
	// Add service-specific common options here if needed.
}

// Option is a function that configures service.Options.
type Option func(*Options)
