/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package errors

import (
	"fmt"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	tkerrors "github.com/origadmin/toolkits/errors"
)

// Meta defines the type for error metadata
type Meta map[string]interface{}

// WithMeta adds metadata to an error
// If the error is of type *kerrors.Error, it adds the metadata directly
// Otherwise, it converts the error to *kerrors.Error and then adds the metadata
func WithMeta(err error, meta Meta) *kerrors.Error {
	if err == nil {
		return nil
	}

	// Use FromError to handle wrapped errors properly
	kerr := kerrors.FromError(err)

	// Add metadata
	if kerr.Metadata == nil {
		kerr.Metadata = make(map[string]string)
	}
	for k, v := range meta {
		kerr.Metadata[k] = fmt.Sprintf("%v", v)
	}

	return kerr
}

// WithField adds a single key-value pair as metadata to an error
func WithField(err error, key string, value interface{}) *kerrors.Error {
	return WithMeta(err, Meta{key: value})
}

// NewWithMeta creates a new error with metadata
func NewWithMeta(code int, reason, message string, meta Meta) *kerrors.Error {
	return WithMeta(kerrors.New(code, reason, message), meta)
}

// ErrorfWithMeta creates a formatted error with metadata
func ErrorfWithMeta(code int, reason string, meta Meta, format string, a ...interface{}) *kerrors.Error {
	return NewWithMeta(code, reason, fmt.Sprintf(format, a...), meta)
}

// FromErrorWithMeta converts a standard error to a Kratos error with metadata
func FromErrorWithMeta(err error, meta Meta) *kerrors.Error {
	if err == nil {
		return nil
	}

	// Handle errors that implement the ErrorWithCode interface
	if tkerrors.IsType[tkerrors.ErrorWithCode](err) {
		if codeErr, ok := err.(interface{ Code() int }); ok {
			return NewWithMeta(codeErr.Code(), "CUSTOM_ERROR", err.Error(), meta)
		}
	}

	return WithMeta(err, meta)
}
