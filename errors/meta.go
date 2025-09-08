/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package errors

import (
	"errors"
	"fmt"

	kerrors "github.com/go-kratos/kratos/v2/errors"
)

// WithMessage enhances an error with a more specific message.
// It converts the error to a Kratos error first, then sets the message.
func WithMessage(err error, format string, args ...interface{}) *kerrors.Error {
	ke := Convert(err)
	ke.Message = fmt.Sprintf(format, args...)
	return ke
}

// WithMeta adds structured metadata to an error.
// It converts the error to a Kratos error first, then adds the metadata.
func WithMeta(err error, meta map[string]interface{}) *kerrors.Error {
	ke := Convert(err)
	if ke.Metadata == nil {
		ke.Metadata = make(map[string]string)
	}
	for k, v := range meta {
		ke.Metadata[k] = fmt.Sprintf("%v", v) // Convert interface{} to string
	}
	return ke
}

// WithField adds a single key-value pair as metadata to an error.
// It converts the error to a Kratos error first, then adds the field.
func WithField(err error, key string, value interface{}) *kerrors.Error {
	return WithMeta(err, map[string]interface{}{key: value})
}

// LookupMeta retrieves a specific metadata value from a Kratos error.
// It returns the value and a boolean indicating if the key was found.
func LookupMeta(err error, key string) (string, bool) {
	var ke *kerrors.Error
	if !errors.As(err, &ke) || ke.Metadata == nil {
		return "", false
	}
	val, ok := ke.Metadata[key]
	return val, ok
}
