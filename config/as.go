package config

import (
	"reflect"

	"github.com/origadmin/runtime/interfaces"
)

// As finds the first decoder in the chain that matches the type of target.
// If a match is found, it sets target to that decoder value and returns true.
//
// The chain consists of d itself, followed by the sequence of decoders obtained by
// recursively calling As on the preceding decoder. The unwrapping logic is entirely
// controlled by the decoder's implementation of the optional As(any) bool method.
func As(d interfaces.ConfigDecoder, target any) bool {
	if d == nil || target == nil {
		return false
	}
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}

	targetType := val.Type().Elem()

	// Check if the current decoder itself matches the target type.
	if reflect.TypeOf(d).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(d))
		return true
	}

	// Check if the decoder implements the optional As(any) bool method.
	// If it does, we delegate the responsibility of matching and unwrapping to it.
	type assignable interface {
		As(any) bool
	}
	if as, ok := d.(assignable); ok {
		return as.As(target)
	}

	// Recursively call As on the preceding decoder.
	return false
}
