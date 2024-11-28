/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

// Validator is an interface for validating protobuf messages.
type Validator interface {
	Validate(message proto.Message) error
}

// validate is a struct that implements the Validator interface.
type validate struct {
	v *protovalidate.Validator
}

// Validate validates a protobuf message.
func (v validate) Validate(message proto.Message) error {
	return v.Validate(message)
}

// NewValidate creates a new Validator.
func NewValidate(opts ...ProtoValidatorOption) (Validator, error) {
	v, err := NewProtoValidate(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	return v, nil
}
