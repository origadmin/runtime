/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"

	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/log"
)

// v2Validator is an interface for validating protobuf messages.
type v2Validator interface {
	Validate(message proto.Message) error
}

// validate is a struct that implements the v2Validator interface.
type validateV2 struct {
	v protovalidate.Validator
}

// ValidateV2 validates a protobuf message.
func (v validateV2) ValidateV2(message proto.Message) error {
	return v.ValidateV2(message)
}

func (v validateV2) Validate(ctx context.Context, req any) error {
	log.Debugf("validateV2 Validate called with request: %+v", req)
	if message, ok := req.(proto.Message); ok {
		log.Debugf("validateV2 Validate: request is a proto.Message")
		err := v.v.Validate(message)
		if err != nil {
			log.Warnf("validateV2 Validate: validation failed: %v", err)
		}
		return err
	}
	log.Debugf("validateV2 Validate: request is not a proto.Message")
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
