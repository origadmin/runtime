package runtime

import (
	"github.com/origadmin/runtime/bootstrap"
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/interfaces/options"
)

// appOptions holds the configurable settings for a App.
type appOptions struct {
	appInfo *appInfo // Use the concrete struct type for internal operations.

	// Other options
	bootstrapOpts   []options.Option
	containerOpts   []options.Option
	kratosAppOpts   []options.Option
	structuredCfg   interfaces.StructuredConfig
	config          interfaces.Config
	bootstrapResult bootstrap.Result
}

type Option = options.Option
