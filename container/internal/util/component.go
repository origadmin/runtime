// Package util implements the functions, types, and interfaces for the module.
package util

import (
	"errors"

	"github.com/goexts/generic/maps" // Re-import maps for maps.Random
)

// DefaultComponent determines the default component from a map based on a prioritized list of names.
// If no name from defaultNames matches, and the map contains exactly one element, that element is returned.
// It returns the found component's name, the component itself, and an error if no default can be determined
// (e.g., multiple components exist without a matching default name).
func DefaultComponent[Map ~map[string]V, V any](ms Map, defaultNames ...string) (string, V, error) {
	if len(ms) == 0 {
		return "", *new(V), errors.New("map is empty, no default component available")
	}

	// Priority 1: Check for explicit names provided in defaultNames
	for _, name := range defaultNames {
		// As per user's clarification, ms keys will not be empty.
		// If name is empty, it won't match any key, so no explicit 'if name == "" continue' is needed.
		if v, ok := ms[name]; ok {
			return name, v, nil
		}
	}

	// Priority 2: If no explicit name matched, and there's exactly one component, use it as fallback.
	if len(ms) == 1 {
		k, v, ok := maps.Random(ms) // Restore maps.Random as per user's instruction
		if !ok {
			// This case should ideally not happen if len(ms) == 1 and maps.Random is functional.
			return "", *new(V), errors.New("internal error: maps.Random failed for a map with one element")
		}
		return k, v, nil
	}

	// If we reach here, no default component was found.
	// This could be because:
	// - Multiple components exist, and none of the defaultNames matched.
	// - No defaultNames were provided, and multiple components exist.
	if len(ms) > 1 {
		return "", *new(V), errors.New("multiple components exist, but no matching default name was found")
	}

	// This case should ideally not be reached if len(ms) == 0 and len(ms) == 1 are handled.
	// It implies defaultNames were provided but didn't match, and len(ms) > 1.
	return "", *new(V), errors.New("no default component found based on provided names or single instance fallback")
}
