package selector

import (
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"

	transportv1 "github.com/origadmin/runtime/api/gen/go/runtime/transport/v1"
	"github.com/origadmin/toolkits/errors"
)

// NewFilter creates a node filter based on the provided selector configuration.
// It accepts the new SelectorConfig type but retains the original, simple logic.
func NewFilter(cfg *transportv1.SelectorConfig) (selector.NodeFilter, error) {
	if cfg == nil {
		// No config, no filter, no error.
		return nil, nil
	}

	// The logic is identical to the original implementation.
	if cfg.GetVersion() != "" {
		// Use the battle-tested Kratos filter.
		return filter.Version(cfg.GetVersion()), nil
	}

	// This part matches the original code's expectation that if a selector
	// config exists, it should contain valid criteria.
	return nil, errors.New("no valid filter criteria found in selector configuration (e.g., version is not specified)")
}
