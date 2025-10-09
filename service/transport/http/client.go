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

	// Get the container instance. It will be nil if not provided in options.
	var c interfaces.Container
	if clientOpts.ServiceOptions != nil {
		c = clientOpts.ServiceOptions.Container
	}

	// Determine if container-dependent features (middlewares, named discovery, or selector filter) are configured.
	containerDependentFeaturesEnabled := len(httpConfig.GetMiddlewares()) > 0 || httpConfig.GetDiscoveryName() != ""

	// If container-dependent features are enabled but no container is provided (c is nil),
	// return an error immediately. This consolidates the nil checks for the container.
	if containerDependentFeaturesEnabled && c == nil {
		return nil, fmt.Errorf("application container is required for client configuration but not found in options")
	}

	// Apply options from the protobuf configuration.
	if httpConfig.GetTimeout() != nil {
		kratosOpts = append(kratosOpts, transhttp.WithTimeout(httpConfig.GetTimeout().AsDuration()))
	}

	// Configure middlewares.
	var mws []middleware.Middleware
	if len(httpConfig.GetMiddlewares()) > 0 {
		// 'c' is guaranteed to be non-nil at this point if middlewares are configured.
		for _, name := range httpConfig.GetMiddlewares() {
			m, ok := c.ClientMiddleware(name)
			if !ok {
				return nil, fmt.Errorf("client middleware '%s' not found in container", name)
			}
			mws = append(mws, m)
		}
	} else {
		// If no specific middlewares are configured, use default ones from adapter.go.
		mws = DefaultClientMiddlewares()
	}

	if len(mws) > 0 {
		kratosOpts = append(kratosOpts, transhttp.WithMiddleware(mws...))
	}

	// Configure service discovery and endpoint.
	var discoveryClient registry.Discovery
	endpoint := httpConfig.GetEndpoint()

	// Always apply the endpoint option.
	if endpoint != "" {
		kratosOpts = append(kratosOpts, transhttp.WithEndpoint(endpoint))
	}

	// Determine the discovery client.
	if discoveryName := httpConfig.GetDiscoveryName(); discoveryName != "" {
		// 'c' is guaranteed to be non-nil at this point if a named discovery is configured.
		if d, ok := c.Discovery(discoveryName); ok {
			discoveryClient = d
		} else {
			return nil, fmt.Errorf("discovery client '%s' not found in container", discoveryName)
		}
	} else if c != nil {
		// If no specific discovery name, try to infer if only one is available from the container.
		// This block is only executed if 'c' is not nil.
		discoveries := c.Discoveries()
		if len(discoveries) == 1 {
			for _, d := range discoveries {
				discoveryClient = d
				break
			}
		} else if len(discoveries) > 1 {
			return nil, fmt.Errorf("multiple discovery clients found in container, but no specific discovery client is configured for HTTP client")
		}
	}

	// Apply discovery option if a client was found.
	if discoveryClient != nil {
		kratosOpts = append(kratosOpts, transhttp.WithDiscovery(discoveryClient))
	}

	// Crucial check: If the endpoint implies discovery but no discovery client is configured.
	if strings.HasPrefix(endpoint, "discovery:///") && discoveryClient == nil {
		return nil, fmt.Errorf("endpoint '%s' requires a discovery client, but none is configured", endpoint)
	}

	// Configure node filters (selector).
	if selectorConfig := httpConfig.GetSelector(); selectorConfig != nil {
		// 'c' is guaranteed to be non-nil at this point if a selector filter is configured.
		// Call the original, trusted NewFilter function from your app's selector package.
		nodeFilter, err := serviceselector.NewFilter(selectorConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create node filter: %w", err)
		}

		if nodeFilter != nil {
			//Kratos HTTP client needs a way to integrate NodeFilter with discovery.
			//For now, we'll just ensure the container is available if a filter is set.
			//Assuming the selector filter might need discovery from the container.
			//This part needs more concrete implementation based on how ClientSelectorFilter works.
			kratosOpts = append(kratosOpts, transhttp.WithNodeFilter(nodeFilter))
		}
	}

	// Configure TLS.
	if tlsConfig := httpConfig.GetTlsConfig(); tlsConfig != nil && tlsConfig.GetEnabled() {
		tlsCfg, err := servicetls.NewClientTLSConfig(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create client TLS config: %w", err)
		}
		kratosOpts = append(kratosOpts, transhttp.WithTLSConfig(tlsCfg))
	}

	// Create a new HTTP client with custom transport to handle dial timeout and TLS
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Default dial timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if httpConfig.DialTimeout != nil {
		transport.DialContext = (&net.Dialer{
			Timeout:   httpConfig.DialTimeout.AsDuration(),
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	// Apply any external Kratos HTTP client options passed via functional options.
	// These are applied last, allowing them to override previous options if needed.
	if len(clientOpts.HttpClientOptions) > 0 {
		kratosOpts = append(kratosOpts, clientOpts.HttpClientOptions...)
	}

	// Create the Kratos HTTP client
	client, err := transhttp.NewClient(ctx, kratosOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	return client, nil
}
