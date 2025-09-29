package http_test

import (
	"testing"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"

	// Corrected import path for the http package being tested
	"github.com/origadmin/runtime/service/transport/http"
)

func TestWithServerOption(t *testing.T) {
	// Create test server options for Kratos
	kratosServerOpt1 := transhttp.Address(":8080")
	kratosServerOpt2 := transhttp.Timeout(0)

	t.Run("should correctly apply a single server option", func(t *testing.T) {
		// Create an interfaces.Option using WithServerOption
		option := http.WithServerOption(kratosServerOpt1)

		// Retrieve the Kratos server options using FromServerOptions
		// FromServerOptions internally applies the provided interfaces.Option to its own context.
		serverOpts := http.FromServerOptions(option)

		assert.Len(t, serverOpts, 1)
		// In a real scenario, you might want to inspect the content of serverOpts
		// to ensure the correct Kratos option was applied.
		// For simplicity, we only check length here.
	})

	t.Run("should correctly apply multiple server options", func(t *testing.T) {
		// Create interfaces.Option functions for multiple Kratos server options
		option1 := http.WithServerOption(kratosServerOpt1)
		option2 := http.WithServerOption(kratosServerOpt2)

		// Retrieve the Kratos server options
		serverOpts := http.FromServerOptions(option1, option2)

		assert.Len(t, serverOpts, 2)

	})
}

func TestWithClientOption(t *testing.T) {
	// Create test client options for Kratos
	kratosClientOpt1 := transhttp.WithUserAgent("test-agent")
	kratosClientOpt2 := transhttp.WithTimeout(0)

	t.Run("should correctly apply a single client option", func(t *testing.T) {
		// Create an interfaces.Option using WithClientOption
		option := http.WithClientOption(kratosClientOpt1)

		// Retrieve the Kratos client options using FromClientOptions
		clientOpts := http.FromClientOptions(option)

		assert.Len(t, clientOpts, 1)
	})

	t.Run("should correctly apply multiple client options", func(t *testing.T) {
		// Create interfaces.Option functions for multiple Kratos client options
		option1 := http.WithClientOption(kratosClientOpt1)
		option2 := http.WithClientOption(kratosClientOpt2)

		// Retrieve the Kratos client options
		clientOpts := http.FromClientOptions(option1, option2)

		assert.Len(t, clientOpts, 2)
	})
}

func TestFromOptions_Empty(t *testing.T) {
	t.Run("should return empty slice when no options are provided to FromServerOptions", func(t *testing.T) {
		serverOpts := http.FromServerOptions() // No options passed
		assert.Empty(t, serverOpts)
	})

	t.Run("should return empty slice when no options are provided to FromClientOptions", func(t *testing.T) {
		clientOpts := http.FromClientOptions() // No options passed
		assert.Empty(t, clientOpts)
	})

	t.Run("should return empty slice when nil options are provided to FromServerOptions", func(t *testing.T) {
		serverOpts := http.FromServerOptions(nil) // nil option passed
		assert.Empty(t, serverOpts)
	})

	t.Run("should return empty slice when nil options are provided to FromClientOptions", func(t *testing.T) {
		clientOpts := http.FromClientOptions(nil) // nil option passed
		assert.Empty(t, clientOpts)
	})
}

func TestOptionChaining(t *testing.T) {
	// Create Kratos options
	kratosServerOpt := transhttp.Address(":8080")
	kratosClientOpt := transhttp.WithUserAgent("test")

	// Create interfaces.Option functions
	serverOptionFunc := http.WithServerOption(kratosServerOpt)
	clientOptionFunc := http.WithClientOption(kratosClientOpt)

	t.Run("FromServerOptions should only retrieve server options", func(t *testing.T) {
		// Pass both server and client option functions to FromServerOptions
		serverOpts := http.FromServerOptions(serverOptionFunc, clientOptionFunc)
		assert.Len(t, serverOpts, 1)
		// Further checks could ensure the content of the option is correct
	})

	t.Run("FromClientOptions should only retrieve client options", func(t *testing.T) {
		// Pass both server and client option functions to FromClientOptions
		clientOpts := http.FromClientOptions(serverOptionFunc, clientOptionFunc)
		assert.Len(t, clientOpts, 1)
		// Further checks could ensure the content of the option is correct
	})
}
