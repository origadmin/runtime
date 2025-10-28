package customize

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	extensionv1 "github.com/origadmin/runtime/api/gen/go/runtime/extension/v1"
)

// GetTypedConfig is a helper function to convert the configuration to a typed protobuf message.
// It's a generic function that works with any protobuf message type.
func GetTypedConfig[T proto.Message](cfg *extensionv1.CustomizeConfig, msg T) (T, error) {
	if cfg == nil {
		return msg, fmt.Errorf("config is nil")
	}

	// Get the config field (previously data)
	configStruct := cfg.GetConfig()
	if configStruct == nil {
		return msg, fmt.Errorf("config field is nil")
	}

	// Convert struct to JSON bytes
	jsonBytes, err := protojson.Marshal(configStruct)
	if err != nil {
		return msg, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Unmarshal JSON to the target message type
	if err := protojson.Unmarshal(jsonBytes, msg); err != nil {
		return msg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return msg, nil
}
