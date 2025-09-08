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
	// We can import toolkits/errors in the future to add more specific error inspection.
	// toolkits "github.com/origadmin/toolkits/errors"
)

//go:generate adptool .
//go:adapter:package github.com/go-kratos/kratos/v2/errors kerrors

// Convert takes any standard Go error and converts it into a structured Kratos error.
// This function is the core of the error handling package, acting as the bridge
// between internal business logic errors and standardized API errors.
func Convert(err error) *kerrors.Error {
	if err == nil {
		return nil
	}
	// Use Kratos's built-in converter as a robust baseline.
	// It correctly handles errors that are already *kerrors.Error.
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
		// Log the original error with request context.
		if tr, ok := transport.FromServerContext(r.Context()); ok {
			log.Context(r.Context()).Errorw(
				"kind", "server",
				"path", tr.Operation(),
				"error", err.Error(),
			)
		}

		// Convert the error to our standardized format.
		ke := Convert(err)

		// Delegate to the default encoder to write the final response.
		transhttp.DefaultErrorEncoder(w, r, ke)
	}
}
