/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"context"
)

// The v1ValidatorAll interface at protoc-gen-validate main branch.
// See https://github.com/envoyproxy/protoc-gen-validate/pull/468.
type v1ValidatorAll interface {
	ValidateAll() error
}

// The validate interface starting with protoc-gen-validate v0.6.0.
// See https://github.com/envoyproxy/protoc-gen-validate/pull/455.
type v1Validator interface {
	Validate(all bool) error
}

// The validate interface prior to protoc-gen-validate v0.6.0.
type v1ValidatorLegacy interface {
	Validate() error
}

type validateV1 struct {
	failFast bool
	callback OnValidationErrCallback
}

func (v validateV1) Validate(ctx context.Context, req interface{}) (err error) {
	switch val := req.(type) {
	case v1Validator:
		err = val.Validate(!v.failFast)
	case v1ValidatorLegacy:
		err = val.Validate()
	case v1ValidatorAll:
		err = val.ValidateAll()
	}
	if err != nil {
		if v.callback != nil {
			v.callback(ctx, err)
		}
		return err
	}

	return nil
}

func NewValidateV1(failFast bool, callback OnValidationErrCallback) Validator {
	return validateV1{failFast: failFast, callback: callback}
}
