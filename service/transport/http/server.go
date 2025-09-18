package http

import (
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/configure"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/runtime/service"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewServer creates a new HTTP server with the given configuration and options.
// It is the recommended way to create a server when the protocol is known in advance.
func NewServer(cfg *configv1.Service, opts ...service.Option) (*transhttp.Server, error) {
	ll := log.NewHelper(log.With(log.GetLogger(), "module", "service/http"))
	ll.Debugf("Creating new HTTP server instance with config: %+v", cfg)

	// Get base configuration from the service config
	serverOpts, err := adaptServerConfig(cfg)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to adapt server config for HTTP server creation")
	}

	// Apply any additional options
	svcOpts := configure.Apply(service.DefaultServerOptions(), opts)

	// Apply any options from context
	serverOptsFromContext := FromServerOptions(svcOpts)
	serverOpts = append(serverOpts, serverOptsFromContext...)

	// Create and return the server
	return transhttp.NewServer(serverOpts...), nil
}
