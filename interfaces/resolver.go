package interfaces

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
)

// NewResolver creates a new Resolved instance from a Kratos config.
// This function is the entry point for creating a Resolved instance.
func NewResolver(config kratosconfig.Config) (Resolved, error) {
	var r resolver
	err := config.Scan(&r.values) // Scan the entire config into the internal map
	if err != nil {
		return nil, err
	}
	return &r, nil
}

type resolver struct {
	values map[string]any
}

// Decode implements the Resolved interface.
func (r *resolver) Decode(key string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	if key == "" {
		// If key is empty, decode the entire config
		return mapstructure.Decode(r.values, target)
	}

	// Navigate through the map using the dot-separated key
	currentValue := r.values
	keys := strings.Split(key, ".")
	for i, k := range keys {
		if val, ok := currentValue[k]; ok {
			if i == len(keys)-1 {
				// Last key, decode the value
				return mapstructure.Decode(val, target)
			}
			// Not the last key, continue navigating
			if nextMap, isMap := val.(map[string]any); isMap {
				currentValue = nextMap
			} else {
				return fmt.Errorf("config path '%s' is not a map at key '%s'", key, k)
			}
		} else {
			return fmt.Errorf("config key '%s' not found at path '%s'", k, key)
		}
	}
	return nil // Should not be reached if key is not empty
}
