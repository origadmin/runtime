package extension

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	extensionv1 "github.com/origadmin/runtime/api/gen/go/runtime/extension/v1"
)

// GetTypedConfig retrieves a configuration of the specified type and name from the ExtensionMap
// and deserializes it into the result object.
// This function uses generics, requiring the generic type T to implement the proto.Message interface.
//
// Parameters:
//   - m: A pointer to extensionv1.ExtensionMap containing configuration information.
//   - configType: The type of the configuration.
//   - name: The name of the configuration.
//   - result: An object used to receive the deserialized configuration.
//
// Returns:
//   - bool: Indicates whether the configuration was successfully found and deserialized.
//   - error: Represents any errors that occurred during the operation.
func GetTypedConfig[T proto.Message](m *extensionv1.ExtensionMap, configType, name string, result T) (bool, error) {
	// Find the configuration set of the specified type from the ExtensionMap
	configs := ConfigsFromType(m, configType)
	if configs == nil {
		// If the configuration set is not found, return false and nil error
		return false, nil
	}

	for _, cfg := range configs {
		// Check if the configuration is enabled
		if !cfg.Enabled {
			// If the configuration is not enabled, return false and nil error
			continue
		}
		if cfg.GetName() != name {
			continue
		}
		// Deserialize the configuration value into the result object
		if err := cfg.GetValue().UnmarshalTo(result); err != nil {
			// If deserialization fails, return false and an error with detailed information
			return false, fmt.Errorf("unmarshal %s/%s config: %w", configType, name, err)
		}
		// Configuration was successfully found and deserialized, return true and nil error
		return true, nil
	}
	// Configuration was not found, return false and nil error
	return false, nil
}

// ConfigFromName returns the config with the given name.
func ConfigFromName(cc *extensionv1.Extension, name string) *extensionv1.Extension_Config {
	configs := cc.GetConfigs()
	if configs == nil {
		return nil
	}

	for _, config := range configs {
		if config.GetName() == name {
			return config
		}
	}
	return nil
}

// ConfigsFromType returns all configs with the given type.
func ConfigsFromType(cm *extensionv1.ExtensionMap, typo string) []*extensionv1.Extension_Config {
	types := cm.GetTypes()
	if types == nil {
		return nil
	}

	if ret, ok := types[typo]; ok {
		return ret.GetConfigs()
	}
	return nil
}
