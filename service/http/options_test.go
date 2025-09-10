package http_test

import (
	"testing"

	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/stretchr/testify/assert"

	"github.com/origadmin/runtime/service"
	"github.com/origadmin/runtime/service/http"
)

func TestWithServerOption(t *testing.T) {
	// Create test server options
	serverOpt1 := transhttp.Address(":8080")
	serverOpt2 := transhttp.Timeout(0)

	// Test with single option
	opt1 := service.DefaultServerOptions()
	http.WithServerOption(serverOpt1)(opt1.(*service.Options))

	// Test with multiple options
	opt2 := service.DefaultServerOptions()
	http.WithServerOption(serverOpt1, serverOpt2)(opt2.(*service.Options))

	// Verify options
	opts1 := http.FromServerOptions(opt1.(*service.Options))
	assert.Len(t, opts1, 1)

	opts2 := http.FromServerOptions(opt2.(*service.Options))
	assert.Len(t, opts2, 2)
}

func TestWithClientOption(t *testing.T) {
	// Create test client options
	clientOpt1 := transhttp.WithUserAgent("test-agent")
	clientOpt2 := transhttp.WithTimeout(0)

	// Test with single option
	opt1 := service.DefaultServerOptions()
	http.WithClientOption(clientOpt1)(opt1.(*service.Options))

	// Test with multiple options
	opt2 := service.DefaultServerOptions()
	http.WithClientOption(clientOpt1, clientOpt2)(opt2.(*service.Options))

	// Verify options
	opts1 := http.FromClientOptions(opt1.(*service.Options))
	assert.Len(t, opts1, 1)

	opts2 := http.FromClientOptions(opt2.(*service.Options))
	assert.Len(t, opts2, 2)
}

func TestFromOptions_Empty(t *testing.T) {
	// Test with nil options
	t.Run("nil options", func(t *testing.T) {
		assert.Empty(t, http.FromServerOptions(nil))
		assert.Empty(t, http.FromClientOptions(nil))
	})

	// Test with empty options
	t.Run("empty options", func(t *testing.T) {
		opts := &service.Options{}
		assert.Empty(t, http.FromServerOptions(opts))
		assert.Empty(t, http.FromClientOptions(opts))
	})
}

func TestOptionChaining(t *testing.T) {
	// Test chaining multiple option functions
	opts := service.DefaultServerOptions()

	http.WithServerOption(transhttp.Address(":8080"))(opts.(*service.Options))
	http.WithClientOption(transhttp.WithUserAgent("test"))(opts.(*service.Options))

	serverOpts := http.FromServerOptions(opts.(*service.Options))
	clientOpts := http.FromClientOptions(opts.(*service.Options))

	assert.Len(t, serverOpts, 1)
	assert.Len(t, clientOpts, 1)
}
