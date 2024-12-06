/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime/context"
)

// v2Validator is an interface for validating protobuf messages.
type v2Validator interface {
	Validate(message proto.Message) error
}

// validate is a struct that implements the v2Validator interface.
type validateV2 struct {
	v *protovalidate.Validator
}

// ValidateV2 validates a protobuf message.
func (v validateV2) ValidateV2(message proto.Message) error {
	return v.ValidateV2(message)
}

func (v validateV2) Validate(ctx context.Context, req any) error {
	if message, ok := req.(proto.Message); ok {
		return v.v.Validate(message)
	}
	return nil
}

// NewValidateV2 creates a new v2Validator.
func NewValidateV2(opts ...ProtoValidatorOption) (Validator, error) {
	v, err := NewProtoValidate(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize v1Validator: %w", err)
	}

	return &validateV2{v: v}, nil
}
