package eventbus

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// ProtoCoder implements eventbus.Coder for Protobuf messages.
type ProtoCoder struct{}

// NewProtoCoder creates a new ProtoCoder.
func NewProtoCoder() *ProtoCoder {
	return &ProtoCoder{}
}

func (c *ProtoCoder) Encode(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("expected proto.Message, got %T", v)
	}
	return proto.Marshal(msg)
}

func (c *ProtoCoder) Decode(data []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("expected proto.Message, got %T", v)
	}
	return proto.Unmarshal(data, msg)
}

// JsonCoder implements eventbus.Coder for JSON.
type JsonCoder struct{}

// NewJsonCoder creates a new JsonCoder.
func NewJsonCoder() *JsonCoder {
	return &JsonCoder{}
}

func (c *JsonCoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JsonCoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

//var _ eventbus.Coder = (*ProtoCoder)(nil)
//var _ eventbus.Coder = (*JsonCoder)(nil)
