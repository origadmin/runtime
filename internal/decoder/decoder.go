package decoder

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	kratosconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/mitchellh/mapstructure"

	"google.golang.org/protobuf/types/known/durationpb"

	base "github.com/origadmin/runtime/config/decoder" // Import the new public decoder package with an alias
	"github.com/origadmin/runtime/interfaces"
)

// stringToDurationHookFunc converts a string duration (e.g., "1s", "1m") to a *durationpb.Duration.
func stringToDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		// Only handle string to *durationpb.Duration conversions.
		if f.Kind() != reflect.String || t != reflect.TypeOf(&durationpb.Duration{}) {
			return data, nil
		}

		str, ok := data.(string)
		if !ok {
			return data, nil
		}

		// Parse the string into a time.Duration.
		d, err := time.ParseDuration(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration string '%s': %w", str, err)
		}

		// Convert time.Duration to *durationpb.Duration.
		return durationpb.New(d), nil
	}
}

// defaultDecoder is the default implementation of the interfaces.ConfigDecoder.
// It embeds BaseDecoder to inherit default behaviors and stores the entire
// configuration in a map for generic decoding.
type defaultDecoder struct {
	*base.BaseDecoder // Embed the BaseDecoder from the new public package
	values            map[string]any
}

// NewDecoder creates a new default decoder instance.
// It initializes the BaseDecoder and scans the entire configuration into an internal map
// to support generic decoding lookups.
func NewDecoder(config kratosconfig.Config) (interfaces.ConfigDecoder, error) {
	d := &defaultDecoder{
		BaseDecoder: base.NewBaseDecoder(config), // Initialize the embedded BaseDecoder
		values:      make(map[string]any),
	}

	// Scan the entire config into the internal map for generic decoding.
	if err := config.Scan(&d.values); err != nil {
		return nil, fmt.Errorf("failed to scan config into decoder values: %w", err)
	}

	// Ensure that after scanning, d.values is not empty, indicating successful load.
	if len(d.values) == 0 {
		return nil, fmt.Errorf("decoder values are empty after scanning config")
	}

	return d, nil
}

// getMapstructureDecoderConfig creates a configured mapstructure.DecoderConfig.
// This function is private to the package.
func getMapstructureDecoderConfig(target interface{}) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           target,
		TagName:          "json",
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.TextUnmarshallerHookFunc(),
		),
	}
}

// Decode provides a generic decoding mechanism. It navigates the internal `values` map
// using a dot-separated key and then uses mapstructure to decode the result into the target.
// This method overrides the BaseDecoder's Decode to use the pre-scanned `values` map.
func (d *defaultDecoder) Decode(key string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	var dataToDecode any
	if key == "" {
		// If key is empty, decode the entire config map.
		dataToDecode = d.values
	} else {
		// Navigate through the map using the dot-separated key.
		currentValue, err := d.lookup(key)
		if err != nil {
			return err
		}
		dataToDecode = currentValue
	}

	// Configure mapstructure for flexible decoding.
	config := getMapstructureDecoderConfig(target) // Use the helper function

	msDecoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create mapstructure decoder: %w", err)
	}

	return msDecoder.Decode(dataToDecode)
}

// lookup navigates the nested map `d.values` using a dot-separated key.
func (d *defaultDecoder) lookup(key string) (any, error) {
	var currentValue any = d.values
	keys := strings.Split(key, ".")

	for i, k := range keys {
		currentMap, isMap := currentValue.(map[string]any)
		if !isMap {
			pathSegment := strings.Join(keys[:i], ".")
			return nil, fmt.Errorf("config path '%s' is not a map at segment '%s'", pathSegment, keys[i-1])
		}

		val, ok := currentMap[k]
		if !ok {
			pathSegment := strings.Join(keys[:i+1], ".")
			return nil, fmt.Errorf("config key '%s' not found at path '%s'", k, pathSegment)
		}
		currentValue = val
	}

	return currentValue, nil
}

// DefaultDecoderProvider is the default provider that creates a new decoder instance.
var DefaultDecoderProvider = interfaces.ConfigDecoderFunc(NewDecoder)
