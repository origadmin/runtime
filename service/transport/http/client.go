package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
	"github.com/origadmin/runtime/interfaces"
	serviceselector "github.com/origadmin/runtime/service/selector"
	servicetls "github.com/origadmin/runtime/service/tls"
)

// NewHTTPClient creates a new concrete HTTP client connection based on the provided configuration.
// It returns *transhttp.Client, not the generic interfaces.Client.
func NewHTTPClient(ctx context.Context, httpConfig *transportv1.HttpClientConfig, clientOpts *ClientOptions) (*transhttp.Client, error) {
	// Prepare the Kratos HTTP client options.
	var kratosOpts []transhttp.ClientOption

	// Get the container instance.
	var c interfaces.Container
	if clientOpts.ServiceOptions != nil {
		c = clientOpts.ServiceOptions.Container
	}

	// Centralized check for container dependency.
	if (len(httpConfig.GetMiddlewares()) > 0 || httpConfig.GetDiscoveryName() != "" || httpConfig.GetSelector() != nil) && c == nil {
		return nil, fmt.Errorf("application container is required for client configuration (middlewares, discovery, or selector) but not found in options")
	}

	// Apply options from the protobuf configuration.
	if httpConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transhttp.WithTimeout(httpConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if len(httpConfig.GetMiddlewares()) > 0 {
		for _, name := range httpConfig.GetMiddlewares() {
			m, ok := c.ClientMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("client middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		mws = DefaultClientMiddlewares()
	}
	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transhttp.WithMiddleware(mws...))
	}

	// Configure service discovery and endpoint.
	var discoveryClient registry.Discovery
	endpoint := httpConfig.GetEndpoint()
	if endpoint != "" {
		kratosOpts = append(kratosOpts, transhttp.WithEndpoint(endpoint))
	}

	if discoveryName := httpConfig.GetDiscoveryName(); discoveryName != "" {
		if d, ok := c.Discovery(discoveryName); ok {
			discoveryClient = d
		} else {
			return nil, fmt.Errorf("discovery client '%s' not found in container", discoveryName)
		}
	} else if c != nil {
		discoveries := c.Discoveries()
		if len(discoveries) == 1 {
			for _, d := range discoveries {
				discoveryClient = d
				break
			}
		}
	}

	if discoveryClient != nil {
		kratosOpts = append(kratosOpts, transhttp.WithDiscovery(discoveryClient))
	}

	if strings.HasPrefix(endpoint, "discovery:///") && discoveryClient == nil {
		return nil, fmt.Errorf("endpoint '%s' requires a discovery client, but none is configured", endpoint)
	}

	// Configure node filters (selector).
	if selectorConfig := httpConfig.GetSelector(); selectorConfig != nil {
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create node filter: %w", err)
		}
		if nodeFilter != nil {
			kratosOpts = append(kratosOpts, transhttp.WithNodeFilter(nodeFilter))
		}
	}

	// Create and configure the HTTP transport.
	dialer := &net.Dialer{
		Timeout:   30 * time.Second, // Default dial timeout
		KeepAlive: 30 * time.Second,
	}
	if httpConfig.DialTimeout != nil {
		dialer.Timeout = httpConfig.DialTimeout.AsDuration()
	}

	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Configure TLS and apply it to the transport.
	if tlsConfig := httpConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create client TLS config: %w", err)
		}
		transport.TLSClientConfig = tlsCfg
	}

	// IMPORTANT: Pass the configured transport to the Kratos client.
	kratosOpts = append(kratosOpts, transhttp.WithTransport(transport))

	// Apply any external Kratos HTTP client options passed via functional options.
	if len(clientOpts.HttpClientOptions) > 0 {
		kratosOpts = append(kratosOpts, clientOpts.HttpClientOptions...)
	}

	// Create the Kratos HTTP client.
	client, err := transhttp.NewClient(ctx, kratosOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	return client, nil
}
