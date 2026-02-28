package configutil

import (
	"errors"
	"fmt"
)

// Identifiable represents a configuration item that can be uniquely identified by a name.
// This interface is a constraint for the generic Normalize function.
// We include comparable to allow checking if an item is the zero value (nil for pointers).
type Identifiable interface {
	comparable
	GetName() string
}

// Normalize applies a set of heuristic rules to a configuration structure
// to ensure consistency and provide sensible defaults. It is a generic function
// that operates on any configuration slice whose elements implement the Identifiable interface.
//
// The function does not modify the input slices directly but returns the normalized
// default item and a new list of configuration items.
//
// Normalization Logic:
// 1. Builds a map of all provided configurations for efficient lookup.
// 2. If a `defaultItem` is provided but is not present in the `configs` slice (by name), it is added.
// 3. If the `configs` slice is empty after the above step, an error is returned.
// 4. Determines the definitive default item based on the following priority:
//    a. The item specified by `activeName`, if provided and exists. (Highest Priority)
//    b. The initially provided `defaultItem`, if provided.
//    c. The first item in the `configs` slice. (Fallback)
// 5. Ensures that the returned default item is always the instance from the `configs` slice (object identity consistency).
//
// Type Parameters:
//   - T: The type of the configuration item, which must implement the Identifiable interface.
//        T must be comparable (pointers satisfy this).
//
// Parameters:
//   - activeName: The name of the configuration item that should be active. Can be an empty string.
//   - defaultItem: The default configuration item. Can be nil.
//   - configs: A slice of configuration items.
//
// Returns:
//   - The normalized default configuration item.
//   - The normalized slice of configuration items.
//   - An error if the configuration is invalid (e.g., no configs available) or inconsistent (e.g., activeName not found).
func Normalize[T Identifiable](activeName string, defaultItem T, configs []T) (T, []T, error) {
	var zero T // The zero value for type T (nil if T is a pointer).

	// Create a mutable copy of the configs slice to avoid modifying the original input.
	normalizedConfigs := make([]T, len(configs))
	copy(normalizedConfigs, configs)

	configMap := make(map[string]T)
	for _, c := range normalizedConfigs {
		configMap[c.GetName()] = c
	}

	// Rule 2: If a defaultItem is provided but not in configs (check by Name), add it.
	if defaultItem != zero {
		if _, exists := configMap[defaultItem.GetName()]; !exists {
			normalizedConfigs = append(normalizedConfigs, defaultItem)
			configMap[defaultItem.GetName()] = defaultItem
		}
	}

	// Rule 3: After potential addition, if configs is empty, it's an error.
	if len(normalizedConfigs) == 0 {
		return zero, nil, errors.New("no configurations provided")
	}

	// Rule 4: Determine the definitive default item.
	var definitiveDefault T

	// Priority 1: Active Name (Override)
	if activeName != "" {
		if activeConf, exists := configMap[activeName]; exists {
			definitiveDefault = activeConf
		} else {
			return zero, nil, fmt.Errorf("active configuration '%s' not found in the provided configs", activeName)
		}
	}

	// Priority 2: Provided Default (if Active didn't set it)
	if definitiveDefault == zero && defaultItem != zero {
		// CRITICAL: Always use the instance from the map to ensure object identity consistency.
		// Even if defaultItem was just added in Rule 2, getting it from the map is safer and consistent.
		if val, ok := configMap[defaultItem.GetName()]; ok {
			definitiveDefault = val
		} else {
			// This branch should be unreachable because Rule 2 ensures defaultItem is in the map.
			definitiveDefault = defaultItem
		}
	}

	// Priority 3: Fallback to first item (covers "only one" and "multiple but no default" cases)
	if definitiveDefault == zero && len(normalizedConfigs) > 0 {
		definitiveDefault = normalizedConfigs[0]
	}

	// Rule 5: Final validation.
	if definitiveDefault == zero {
		// This should be unreachable if normalizedConfigs is not empty.
		return zero, nil, errors.New("could not determine a default configuration")
	}

	return definitiveDefault, normalizedConfigs, nil
}
