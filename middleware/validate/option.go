/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import "context"

type Options struct {
	version          Version
	failFast         bool
	callback         OnValidationErrCallback
	validatorOptions []ProtoValidatorOption
}
type Option = func(*Options)

// OnValidationErrCallback is a function that will be invoked on validation error(s).
// It returns true if the error is handled and should be ignored, false otherwise.
type OnValidationErrCallback func(ctx context.Context, err error) bool

// WithOnValidationErrCallback registers function that will be invoked on validation error(s).
func WithOnValidationErrCallback(onValidationErrCallback OnValidationErrCallback) Option {
	return func(o *Options) {
		o.callback = onValidationErrCallback
	}
}

// WithFailFast tells v1Validator to immediately stop doing further validation after first validation error.
// This option is ignored if message is only supporting v1Validator.v1ValidatorLegacy interface.
func WithFailFast(failFast bool) Option {
	return func(o *Options) {
		o.failFast = failFast
	}
}

// WithV2ProtoValidatorOptions registers options for Validator with version 2.
func WithV2ProtoValidatorOptions(opts ...ProtoValidatorOption) Option {
	return func(o *Options) {
		o.version = V2
		o.validatorOptions = opts
	}
}
