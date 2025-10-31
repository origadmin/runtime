package http

import (
	"context"
	"fmt"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	httpv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/http/v1"
)

// NewClient creates a new concrete HTTP client connection based on the provided configuration.
// It returns *transhttp.Client, not the generic interfaces.Client.
func NewClient(ctx context.Context, httpConfig *httpv1.Client, clientOpts *ClientOptions) (*transhttp.Client, error) {
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
