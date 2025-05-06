/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type (
	// ValidationError is an error that occurs during validation.
	ValidationError = protovalidate.ValidationError
	// CompilationError is an error that occurs during compilation.
	CompilationError = protovalidate.CompilationError
	// RuntimeError is an error that occurs during runtime.
	RuntimeError = protovalidate.RuntimeError
	// ProtoValidator is a v1Validator for protobuf messages.
	ProtoValidator = protovalidate.Validator
	// ProtoValidatorOption is an option for the ProtoValidator.
	ProtoValidatorOption = protovalidate.ValidatorOption
)

// NewProtoValidate creates a new ProtoValidator.
func NewProtoValidate(opts ...ProtoValidatorOption) (protovalidate.Validator, error) {
	return protovalidate.New(opts...)
}

// ProtoValidate validates a protobuf message.
func ProtoValidate(message proto.Message) error {
	return protovalidate.Validate(message)
}

// WithMessages sets the messages to validate.
func WithMessages(messages ...proto.Message) ProtoValidatorOption {
	return protovalidate.WithMessages(messages...)
}

// WithDisableLazy disables lazy validation.
func WithDisableLazy() ProtoValidatorOption {
	return protovalidate.WithDisableLazy()
}

// WithExtensionTypeResolver sets the extension type resolver.
func WithExtensionTypeResolver(extensionTypeResolver protoregistry.ExtensionTypeResolver) ProtoValidatorOption {
	return protovalidate.WithExtensionTypeResolver(extensionTypeResolver)
}

// WithAllowUnknownFields allows unknown fields.
func WithAllowUnknownFields() ProtoValidatorOption {
	return protovalidate.WithAllowUnknownFields()
}
