package http

import (
	"context"
	"fmt"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	transportv1 "github.com/origadmin/runtime/api/gen/go/transport/v1"
)

// NewHTTPClient creates a new concrete HTTP client connection based on the provided configuration.
// It returns *transhttp.Client, not the generic interfaces.Client.
func NewHTTPClient(ctx context.Context, httpConfig *transportv1.HttpClientConfig, clientOpts *ClientOptions) (*transhttp.Client, error) {
	// Initialize the Kratos HTTP client options using the adapter function.
	kratosOpts, err := initHttpClientOptions(ctx, httpConfig, clientOpts)
	if err != nil {
		return nil, err
	}

	// Create the Kratos HTTP client.
	client, err := transhttp.NewClient(ctx, kratosOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	return client, nil
}
