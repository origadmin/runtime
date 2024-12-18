/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import (
	"context"
	"fmt"

	"github.com/origadmin/runtime/log"
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
	log.Debugf("Validate called with request: %+v", req)
	switch val := req.(type) {
	case v1Validator:
		log.Debugf("Validating with v1Validator")
		err = val.Validate(!v.failFast)
		log.Debugf("v1Validator returned error: %v", err)
	case v1ValidatorLegacy:
		log.Debugf("Validating with v1ValidatorLegacy")
		err = val.Validate()
		log.Debugf("v1ValidatorLegacy returned error: %v", err)
	case v1ValidatorAll:
		log.Debugf("Validating with v1ValidatorAll")
		err = val.ValidateAll()
		log.Debugf("v1ValidatorAll returned error: %v", err)
	default:
		log.Warnf("Unknown validator type: %+v", req)
		err = fmt.Errorf("unknown validator type: %+v", req)
	}
	if err != nil {
		log.Errorf("Validation failed with error: %v", err)
		if v.callback != nil {
			log.Debugf("Calling callback with error: %v", err)
			v.callback(ctx, err)
		}
		return err
	}
	log.Debugf("Validation successful")
	return nil
}

func NewValidateV1(failFast bool, callback OnValidationErrCallback) Validator {
	return validateV1{failFast: failFast, callback: callback}
}
