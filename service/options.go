package service

import (
	"github.com/origadmin/framework/runtime/interfaces"
)

// Options contains the options for creating service components.
// It embeds interfaces.ContextOptions for common context handling.
type Options struct {
	interfaces.OptionValue
}

// Option is a function that configures service.Options.
type Option func(*Options)

func DefaultServerOptions() interfaces.OptionValue {
	return interfaces.DefaultOptions()
}
