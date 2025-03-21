/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import "context"

type Option struct {
	version          Version
	failFast         bool
	callback         OnValidationErrCallback
	validatorOptions []ProtoValidatorOption
}
type OptionSetting = func(*Option)

// OnValidationErrCallback is a function that will be invoked on validation error(s).
// It returns true if the error is handled and should be ignored, false otherwise.
type OnValidationErrCallback func(ctx context.Context, err error) bool

// WithOnValidationErrCallback registers function that will be invoked on validation error(s).
func WithOnValidationErrCallback(onValidationErrCallback OnValidationErrCallback) OptionSetting {
	return func(o *Option) {
		o.callback = onValidationErrCallback
	}
}

// WithFailFast tells v1Validator to immediately stop doing further validation after first validation error.
// This option is ignored if message is only supporting v1Validator.v1ValidatorLegacy interface.
func WithFailFast(failFast bool) OptionSetting {
	return func(o *Option) {
		o.failFast = failFast
	}
}

// WithV2ProtoValidatorOptions registers options for Validator with version 2.
func WithV2ProtoValidatorOptions(opts ...ProtoValidatorOption) OptionSetting {
	return func(o *Option) {
		o.version = V2
		o.validatorOptions = opts
	}
}
