package service

import (
	"github.com/origadmin/runtime/interfaces"
)

// Options contains the options for creating service components.
// It embeds interfaces.OptionValue for common context handling.
type Options struct {
	interfaces.OptionValue // Updated type
}

// Option is a function that configures service.Options.
type Option func(*Options)

func DefaultServerOptions() *Options {
	return &Options{
		OptionValue: interfaces.DefaultOptions(), // Updated function call
	}
}
