package decoder

import (
	kratosconfig "github.com/go-kratos/kratos/v2/config"

	"github.com/origadmin/runtime/interfaces"
)

// DefaultDecoderProvider is the default implementation of the ConfigDecoderProvider interface.
// It creates a new Decoder instance.
var DefaultDecoderProvider interfaces.ConfigDecoderProvider = interfaces.ConfigDecoderFunc(newDefaultDecoder)

func newDefaultDecoder(config kratosconfig.Config) (interfaces.ConfigDecoder, error) {
	return NewDecoder(config), nil
}
