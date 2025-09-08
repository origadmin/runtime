/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package errors provides a centralized hub for handling, converting, and rendering errors.
// It provides custom error handlers/encoders for HTTP and gRPC that integrate with the
// Kratos ecosystem while providing centralized logging and error conversion.
package errors

import (
	"net/http"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	apiErrors "github.com/origadmin/framework/runtime/api/gen/go/apierrors"
	"github.com/origadmin/runtime/log"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// Convert takes any standard Go error and converts it into a structured Kratos error.
// This function is the core of the error handling package, acting as the bridge
// between internal business logic errors and standardized API errors.
// It uses toolkits/errors for enhanced error handling capabilities.
func Convert(err error) *kerrors.Error {
	if err == nil {
		return nil
	}

	// If it's already a Kratos error, return it directly
	if ke, ok := err.(*kerrors.Error); ok {
		return ke
	}

	// Handle error chains using toolkits/errors
	if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
		if unwrapped := unwrapper.Unwrap(); unwrapped != nil {
			return Convert(unwrapped)
		}
	}

	// Handle error with code from toolkits/errors
	if tkerrors.IsType[tkerrors.ErrorWithCode](err) {
		if codeErr, ok := err.(interface{ Code() int }); ok {
			return kerrors.New(
				codeErr.Code(),
				"CUSTOM_ERROR",
				err.Error(),
			)
		}
	}

	// Convert other error types to Kratos error
	ke := kerrors.FromError(err)

	// If the conversion results in a generic error with a 500 status code,
	// we can provide a more specific, standardized internal error type.
	if ke.Code == 500 && ke.Reason == kerrors.UnknownReason {
		return apiErrors.ErrorInternalServerError(ke.Message)
	}

	return ke
}

// NewErrorEncoder returns a new transhttp.EncodeErrorFunc that provides centralized
// logging and error conversion for all HTTP responses.
// It wraps the default Kratos error encoder, enhancing it with our custom logic.
func NewErrorEncoder() transhttp.EncodeErrorFunc {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			return
		}

		// Convert the error to Kratos error
		ke := Convert(err)

		// Log the error with request context and metadata
		if tr, ok := transport.FromServerContext(r.Context()); ok {
			fields := []interface{}{
				"kind", "server",
				"path", tr.Operation(),
				"error", ke.Message,
				"code", ke.Code,
				"reason", ke.Reason,
			}

			// Add metadata to log if present
			if len(ke.Metadata) > 0 {
				fields = append(fields, "metadata", ke.Metadata)
			}

			// Add error chain information if available
			if cause := tkerrors.Unwrap(err); cause != nil {
				fields = append(fields, "cause", cause.Error())
			}

			log.Context(r.Context()).Errorw(
				fields...,
			)
		}

		// Delegate to the default encoder to write the final response
		transhttp.DefaultErrorEncoder(w, r, ke)
	}
}

// ConvertError converts an error to a Kratos error with additional metadata.
// It uses toolkits/errors for enhanced error handling capabilities.
func ConvertError(err error, meta map[string]string) *kerrors.Error {
	if err == nil {
		return nil
	}

	// Convert the error first
	ke := Convert(err)
	if ke == nil {
		return nil
	}

	// Add metadata if provided
	if len(meta) > 0 {
		if ke.Metadata == nil {
			ke.Metadata = make(map[string]string)
		}
		for k, v := range meta {
			ke.Metadata[k] = v
		}

		// Add metadata to the error chain if using toolkits/errors
		if tkerrors.IsType[tkerrors.ErrorWithCode](err) {
			// If the error supports metadata, add it directly
			if metaSetter, ok := err.(interface{ SetMetadata(map[string]string) }); ok {
				metaSetter.SetMetadata(meta)
			}
		}
	}

	return ke
}

// Wrap wraps an error with a message, preserving the error type.
// It uses toolkits/errors for enhanced error wrapping.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	// Use toolkits/errors for enhanced error wrapping
	return tkerrors.Wrap(err, message)
}

// Wrapf wraps an error with a formatted message, preserving the error type.
// It uses toolkits/errors for enhanced error wrapping.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// Use toolkits/errors for enhanced error wrapping
	return tkerrors.Wrapf(err, format, args...)
}

// IsError checks if any error in the chain matches the target.
// It's a wrapper around toolkits/errors.Is for better error chain handling.
func IsError(err, target error) bool {
	return tkerrors.Is(err, target)
}

// AsError finds the first error in the chain that matches the target type.
// It's a wrapper around toolkits/errors.As for better error chain handling.
func AsError(err error, target interface{}) bool {
	return tkerrors.As(err, target)
}

// UnwrapError returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, UnwrapError returns nil.
// It's a wrapper around toolkits/errors.Unwrap for better error chain handling.
func UnwrapError(err error) error {
	return tkerrors.Unwrap(err)
}
