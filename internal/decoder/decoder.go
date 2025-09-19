package decoder

import (
	"fmt"
	"strings"
	"time"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/mitchellh/mapstructure"

	"github.com/origadmin/runtime/interfaces"
)

// NewDecoder creates a new Decoder instance from a Kratos config.
// This function is the entry point for creating a Decoder instance.
func NewDecoder(config kratosconfig.Config) (interfaces.ConfigDecoder, error) {
	var d decoder
	// Get the root value from the Kratos config and scan it into d.values
	// This ensures d.values contains the actual application config, not the bootstrap structure.
	err := config.Value("").Scan(&d.values)
	if err != nil {
		return nil, fmt.Errorf("failed to scan root config value into decoder: %w", err)
	}
	return &d, nil
}

type decoder struct {
	values map[string]any
}

func (d *decoder) Config() any {
	return d.values
}

// Decode implements the Decoder interface.
func (d *decoder) Decode(key string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	var dataToDecode any
	if key == "" {
		// If key is empty, decode the entire config
		dataToDecode = d.values
	} else {
		// Navigate through the map using the dot-separated key
		var currentValue any = d.values
		keys := strings.Split(key, ".")

		for i, k := range keys {
			currentMap, isMap := currentValue.(map[string]any)
			if !isMap {
				pathSegment := strings.Join(keys[:i], ".")
				return fmt.Errorf("config path '%s' is not a map at segment '%s'", pathSegment, keys[i-1])
			}

			val, ok := currentMap[k]
			if !ok {
				pathSegment := strings.Join(keys[:i+1], ".")
				return fmt.Errorf("config key '%s' not found at path '%s'", k, pathSegment)
			}
			currentValue = val
		}
		dataToDecode = currentValue
	}

	// Configure mapstructure to use "json" tags, allow weakly typed input.
	// RecursiveStructHookFunc is removed for compilation, but may affect nested struct decoding.
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           target,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(time.RFC3339), // Corrected: provide a format
			mapstructure.TextUnmarshallerHookFunc(),
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(dataToDecode)
}

var DefaultDecoder interfaces.ConfigDecoderProvider = interfaces.ConfigDecoderFunc(NewDecoder)
