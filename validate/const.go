/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type (
	// ValidationError is an error that occurs during validation.
	ValidationError = protovalidate.ValidationError
	// CompilationError is an error that occurs during compilation.
	CompilationError = protovalidate.CompilationError
	// RuntimeError is an error that occurs during runtime.
	RuntimeError = protovalidate.RuntimeError
	// ProtoValidator is a validator for protobuf messages.
	ProtoValidator = protovalidate.Validator
	// ProtoValidatorOption is an option for the ProtoValidator.
	ProtoValidatorOption = protovalidate.ValidatorOption
	// StandardConstraintResolver is a constraint resolver for the ProtoValidator.
	StandardConstraintResolver = protovalidate.StandardConstraintResolver
	// StandardConstraintInterceptor is a constraint interceptor for the ProtoValidator.
	StandardConstraintInterceptor = protovalidate.StandardConstraintInterceptor
)

// NewProtoValidate creates a new ProtoValidator.
func NewProtoValidate(opts ...ProtoValidatorOption) (*protovalidate.Validator, error) {
	return protovalidate.New(opts...)
}

// ProtoValidate validates a protobuf message.
func ProtoValidate(message proto.Message) error {
	return protovalidate.Validate(message)
}

// WithUTC sets the time zone to UTC.
func WithUTC(use bool) ProtoValidatorOption {
	return protovalidate.WithUTC(use)
}

// WithFailFast sets the fail fast option.
func WithFailFast(failFast bool) ProtoValidatorOption {
	return protovalidate.WithFailFast(failFast)
}

// WithMessages sets the messages to validate.
func WithMessages(messages ...proto.Message) ProtoValidatorOption {
	return protovalidate.WithMessages(messages...)
}

// WithDescriptors sets the descriptors to validate.
func WithDescriptors(descriptors ...protoreflect.MessageDescriptor) ProtoValidatorOption {
	return protovalidate.WithDescriptors(descriptors...)
}

// WithDisableLazy disables lazy validation.
func WithDisableLazy(disable bool) ProtoValidatorOption {
	return protovalidate.WithDisableLazy(disable)
}

// WithStandardConstraintInterceptor adds a standard constraint interceptor.
func WithStandardConstraintInterceptor(interceptor StandardConstraintInterceptor) ProtoValidatorOption {
	return protovalidate.WithStandardConstraintInterceptor(interceptor)
}

// WithExtensionTypeResolver sets the extension type resolver.
func WithExtensionTypeResolver(extensionTypeResolver protoregistry.ExtensionTypeResolver) ProtoValidatorOption {
	return protovalidate.WithExtensionTypeResolver(extensionTypeResolver)
}

// WithAllowUnknownFields allows unknown fields.
func WithAllowUnknownFields(allow bool) ProtoValidatorOption {
	return protovalidate.WithAllowUnknownFields(allow)
}
