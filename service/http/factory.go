package http

import (
	"time"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/context"

	"github.com/origadmin/runtime/interfaces"
	"github.com/origadmin/runtime/service"
)

const (
	defaultTimeout = 5 * time.Second
)

// httpProtocolFactory implements service.ProtocolFactory for HTTP.
type httpProtocolFactory struct{}

// NewClient creates a new HTTP client instance by delegating to the direct implementation.
func (f *httpProtocolFactory) NewClient(ctx context.Context, cfg *configv1.Service, opts ...service.Option) (interfaces.Client, error) {
	return NewClient(ctx, cfg, opts...)
}

// NewServer creates a new HTTP server instance by delegating to the direct implementation.
func (f *httpProtocolFactory) NewServer(cfg *configv1.Service, opts ...service.Option) (interfaces.Server, error) {
	return NewServer(cfg, opts...)
}

func init() {
	service.RegisterProtocol("http", &httpProtocolFactory{})
}
