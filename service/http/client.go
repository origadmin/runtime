package http

import (
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goexts/generic/configure"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context" // Use project's context
	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// NewClient creates a new HTTP client.
// It is the recommended way to create a client when the protocol is known in advance.
func NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	clientOpts, err := adaptClientConfig(cfg)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to adapt client config for http client creation")
	}

	svcOpts := configure.Apply(&service.Options{}, opts)
	if clientOptsFromCtx := FromClientOptions(svcOpts); len(clientOptsFromCtx) > 0 {
		clientOpts = append(clientOpts, clientOptsFromCtx...)
	}

	client, err := transhttp.NewClient(ctx, clientOpts...)
	if err != nil {
		return nil, tkerrors.Wrapf(err, "failed to create http client")
	}
	return client, nil
}
