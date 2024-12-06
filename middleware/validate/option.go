/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package validate implements the functions, types, and interfaces for the module.
package validate

import "context"

type Config struct {
	version                 Version
	failFast                bool
	onValidationErrCallback OnValidationErrCallback
	protoValidatorOptions   []ProtoValidatorOption
}
type ConfigSetting = func(*Config)

type OnValidationErrCallback func(ctx context.Context, err error)

// WithOnValidationErrCallback registers function that will be invoked on validation error(s).
func WithOnValidationErrCallback(onValidationErrCallback OnValidationErrCallback) ConfigSetting {
	return func(o *Config) {
		o.onValidationErrCallback = onValidationErrCallback
	}
}

// WithFailFast tells v1Validator to immediately stop doing further validation after first validation error.
// This option is ignored if message is only supporting v1Validator.v1ValidatorLegacy interface.
func WithFailFast(failFast bool) ConfigSetting {
	return func(o *Config) {
		o.failFast = failFast
	}
}

// WithV2ProtoValidatorOptions registers options for Validator with version 2.
func WithV2ProtoValidatorOptions(opts ...ProtoValidatorOption) ConfigSetting {
	return func(o *Config) {
		o.version = V2
		o.protoValidatorOptions = opts
	}
}
