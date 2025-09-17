package service

import (
	"github.com/origadmin/runtime/optionutil"
)

type serviceOptions struct {
}

// Options contains the options for creating service components.
// It embeds interfaces.OptionValue for common context handling.
type Options = optionutil.Options[serviceOptions]

// Option is a function that configures service.Options.
type Option func(*Options)
