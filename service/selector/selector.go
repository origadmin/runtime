package selector

import (
	"sync"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/selector/wrr"

	discoveryv1 "github.com/origadmin/runtime/api/gen/go/discovery/v1" // Changed import
	"github.com/origadmin/runtime/log"
	"github.com/origadmin/toolkits/errors"
)

const (
	Random = "random"
	WRR    = "wrr"
	P2C    = "p2c"
)

var (
	once    sync.Once
	builder selector.Builder
)

// NewFilter creates a node filter based on the provided selector configuration.
// It currently supports version-based filtering.
func NewFilter(cfg *discoveryv1.Selector) (selector.NodeFilter, error) {
	if cfg == nil {
		return nil, errors.New("selector configuration is nil")
	}
	// Check if the version is specified in the configuration
	if cfg.GetVersion() != "" {
		// Return the version filter and no error
		return filter.Version(cfg.Version), nil
	}
	// If no version is specified, and no other filter type is supported yet, return an error.
	// This is consistent with the original behavior of expecting a filter to be created.
	return nil, errors.New("no valid filter criteria found in selector configuration (e.g., version is not specified)")
}

// SetSelectorGlobalSelector sets the global selector.
func SetSelectorGlobalSelector(selectorType string) {
	if builder != nil {
		return
	}
	var b selector.Builder
	switch selectorType {
	case Random:
		b = random.NewBuilder()
	case WRR:
		b = wrr.NewBuilder()
	case P2C:
		b = p2c.NewBuilder()
	default:
		log.Warnf("selector type %s is not supported", selectorType)
		return
	}
	once.Do(func() {
		if b != nil {
			builder = b
			// Set global selector
			SetGlobalSelector(builder)
		}
	})
}
