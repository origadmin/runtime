package interfaces

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"

	kratosconfig "github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/log"
)

type ResolveFunc func(config kratosconfig.Config) (Resolved, error)

func (r ResolveFunc) Resolve(config kratosconfig.Config) (Resolved, error) {
	return r(config)
}

type resolver struct {
	values map[string]any
}

// All methods that previously implemented Resolved interface are removed
// as Resolved is now an empty interface.

func (r *resolver) WithDecode(name string, v any, decode func([]byte, any) error) error {
	if v == nil {
		return fmt.Errorf("value %s is nil", name)
	}
	data, err := r.Value(name)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("value %s is nil", name)
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return decode(marshal, v)
}

func (r *resolver) Value(name string) (any, error) {
	v, ok := r.values[name]
	if !ok {
		return nil, fmt.Errorf("value %s not found", name)
	}
	return v, nil
}

func (r *resolver) decodeConfig(key string, target interface{}) bool {
	v, ok := r.values[key]
	if !ok {
		return false
	}
	if err := mapstructure.Decode(v, target); err != nil {
		log.Errorf("Failed to decode config key '%s': %v", key, err)
		return false
	}
	return true
}

var DefaultResolver Resolved = ResolveFunc(func(config kratosconfig.Config) (Resolved, error) {
	var r resolver
	err := config.Scan(&r.values)
	if err != nil {
		return nil, err
	}
	return &r, nil // Return pointer to resolver
})

// All adapter structs are removed as they are no longer needed for an empty Resolved interface.
