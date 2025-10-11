package helper

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/genelet/determined/dethcl"
	"github.com/go-viper/encoding/dotenv"
	"github.com/go-viper/encoding/ini"
	"github.com/go-viper/encoding/javaproperties"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var codec = viper.NewCodecRegistry()

type hclCodec struct {
}

func (h hclCodec) Encode(v map[string]any) ([]byte, error) {
	return dethcl.Marshal(v)
}

func (h hclCodec) Decode(b []byte, v map[string]any) error {
	return dethcl.Unmarshal(b, v)
}

func init() {
	err := codec.RegisterCodec("dotenv", &dotenv.Codec{})
	if err != nil {
		panic(err)
	}
	err = codec.RegisterCodec("hcl", &hclCodec{})
	if err != nil {
		panic(err)
	}
	err = codec.RegisterCodec("ini", &ini.Codec{})
	if err != nil {
		panic(err)
	}
	err = codec.RegisterCodec("properties", &javaproperties.Codec{})
	if err != nil {
		panic(err)
	}
}

// SaveConfigToFile saves a protobuf message to a file in the specified format.
func SaveConfigToFile(t *testing.T, msg proto.Message, path string, formatName string) {
	t.Helper()

	// Convert format name to lowercase and handle special cases
	format := strings.ToLower(formatName)
	if format == "prototext" {
		format = "text"
	}

	// Ensure path has the correct extension
	if !strings.HasSuffix(strings.ToLower(path), "."+format) {
		path = path + "." + format
	}

	// Create a new Viper instance
	v := viper.NewWithOptions(viper.WithCodecRegistry(codec))
	v.SetTypeByDefaultValue(true)

	// Use protojson MarshalOptions to maintain original field names
	marshaler := &protojson.MarshalOptions{
		UseProtoNames:   true,  // Use field names as defined in proto
		EmitUnpopulated: false, // Skip empty fields
		Indent:          "  ",  // Indentation
	}

	// Marshal to JSON
	data, err := marshaler.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Unmarshal to map
	var configMap map[string]interface{}
	if err := json.Unmarshal(data, &configMap); err != nil {
		t.Fatalf("Failed to unmarshal config to map: %v", err)
	}

	// Set config values
	for key, value := range configMap {
		v.Set(key, value)
	}

	// Set the config type based on the format
	v.SetConfigType(format)

	// Write to file
	if err := v.WriteConfigAs(path); err != nil {
		t.Fatalf("Failed to write config file %s (format: %s): %v", path, format, err)
	}
}

// LoadConfigFromFile loads a config file into a protobuf message.
func LoadConfigFromFile(t *testing.T, path string, msg proto.Message) {
	t.Helper()

	// Clean and get file extension
	cleanPath := filepath.Clean(path)
	ext := filepath.Ext(cleanPath)
	if ext != "" {
		ext = ext[1:] // Remove the leading dot
	}
	ext = strings.ToLower(ext)

	// Default path: use Viper for supported decoders (yaml/json/toml/ini/env/properties/hcl)
	v := viper.NewWithOptions(viper.WithCodecRegistry(codec))
	v.SetConfigFile(cleanPath)
	v.SetConfigType(ext)

	if err := v.ReadInConfig(); err != nil {
		if strings.Contains(err.Error(), "decoder not found") {
			t.Skipf("Skip: %s decoder not enabled in current environment: %v", ext, err)
			return
		}
		t.Fatalf("Failed to read config file %s: %v", cleanPath, err)
	}

	var configMap map[string]any
	if ext == "env" || ext == "properties" || ext == "ini" {
		flat := map[string]any{}
		for _, k := range v.AllKeys() {
			flat[k] = v.Get(k)
		}
		configMap = expandFlatKeys(flat)
	} else {
		configMap = v.AllSettings()
	}
	jsonData, err := json.Marshal(configMap)
	if err != nil {
		t.Fatalf("Failed to marshal config to JSON: %v", err)
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	if err := unmarshaler.Unmarshal(jsonData, msg); err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}
}

// CloneProtoMessage creates a deep copy of a protobuf message.
func CloneProtoMessage(t *testing.T, src proto.Message) proto.Message {
	t.Helper()
	dst := src.ProtoReflect().New().Interface()
	proto.Merge(dst, src)
	return dst
}

// expandFlatKeys expands dot-delimited keys with numeric indices into nested maps and slices.
// Example: {
//   "a.b.0.c": 1,
//   "a.b.1.c": 2,
//   "x.y": "z"
// } => {"a":{"b":[{"c":1},{"c":2}]}, "x":{"y":"z"}}
func expandFlatKeys(flat map[string]any) map[string]any {
	root := map[string]any{}
	for key, val := range flat {
		parts := strings.Split(key, ".")
		insertInto(root, parts, val)
	}
	return root
}

func insertInto(cur any, parts []string, val any) {
	if len(parts) == 0 {
		return
	}
	isIndex := func(s string) (int, bool) {
		i, err := strconv.Atoi(s)
		if err != nil || i < 0 {
			return 0, false
		}
		return i, true
	}
	key := parts[0]
	idx, ok := isIndex(key)
	if ok {
		// current should be slice
		slice, _ := cur.([]any)
		// ensure capacity
		for len(slice) <= idx {
			slice = append(slice, nil)
		}
		if len(parts) == 1 {
			slice[idx] = val
		} else {
			next := slice[idx]
			if next == nil {
				// decide next container by looking ahead: if next part is index -> slice, else map
				if _, ok2 := isIndex(parts[1]); ok2 {
					next = []any{}
				} else {
					next = map[string]any{}
				}
			}
			slice[idx] = next
			insertInto(slice[idx], parts[1:], val)
		}
		// write back when cur is the root reference type
		switch c := cur.(type) {
		case map[string]any:
			// cannot index map by integer; this branch won't be used when cur is map
			_ = c
		case []any:
			for i := range c {
				c[i] = slice[i]
			}
		}
		return
	}
	// key is string
	nextMap, okm := cur.(map[string]any)
	if !okm {
		// if cur is nil slice element, convert to map
		if cur == nil {
			nextMap = map[string]any{}
		} else {
			return
		}
	}
	next, exists := nextMap[key]
	if len(parts) == 1 {
		nextMap[key] = val
		return
	}
	if !exists || next == nil {
		// decide container by lookahead
		if _, ok2 := isIndex(parts[1]); ok2 {
			next = []any{}
		} else {
			next = map[string]any{}
		}
		nextMap[key] = next
	}
	insertInto(next, parts[1:], val)
}

func isStringIndex(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	if err != nil || i < 0 {
		return 0, false
	}
	return i, true
}