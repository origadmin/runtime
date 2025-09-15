package decoder

import (
	"fmt"
	"path"
	"strings"

	"github.com/mitchellh/mapstructure"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/interfaces"
)

// NewDecoder creates a new Decoder instance from a Kratos config.
// This function is the entry point for creating a Decoder instance.
func NewDecoder(config kratosconfig.Config) (interfaces.ConfigDecoder, error) {
	var d decoder
	err := config.Scan(&d.values) // Scan the entire config into the internal map
	if err != nil {
		return nil, err
	}
	return &d, nil
}

type decoder struct {
	values map[string]any
}

func (d *decoder) GetService(name string) *configv1.Service {
	var service configv1.Service
	if err := d.Decode(path.Join("services", name), &service); err != nil {
		return nil
	}
	return &service
}

func (d *decoder) GetServices() map[string]*configv1.Service {
	var services map[string]*configv1.Service
	if err := d.Decode("services", &services); err != nil {
		return nil
	}
	return services
}

// Decode implements the Decoder interface.
func (d *decoder) Decode(key string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}
	if key == "" {
		// If key is empty, decode the entire config
		return mapstructure.Decode(d.values, target)
	}

	// Navigate through the map using the dot-separated key
	var currentValue any = d.values
	keys := strings.Split(key, ".")

	for i, k := range keys {
		currentMap, isMap := currentValue.(map[string]any)
		if !isMap {
			// This means a path like "a.b.c" was requested, but "a" or "a.b" is not a map.
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

	// After the loop, currentValue holds the final value to be decoded.
	return mapstructure.Decode(currentValue, target)
}

var DefaultDecoder interfaces.ConfigDecoderProvider = interfaces.ConfigDecoderFunc(NewDecoder)
