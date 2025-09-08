/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package errors provides a centralized hub for handling, converting, and rendering errors.
// It uses a proto-defined enum for standard error reasons and adapts the Kratos error package.
package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/origadmin/framework/runtime/api/gen/go/apierrors"
	"github.com/origadmin/runtime/context"
	"github.com/origadmin/runtime/log"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// RequestTimeout creates a 408 Request Timeout error.
func RequestTimeout(reason, message string) *kerrors.Error {
	return New(http.StatusRequestTimeout, reason, message)
}

// MethodNotAllowed creates a 405 Method Not Allowed error.
func MethodNotAllowed(reason, message string) *kerrors.Error {
	return New(http.StatusMethodNotAllowed, reason, message)
}

// TooManyRequests creates a 429 Too Many Requests error.
func TooManyRequests(reason, message string) *kerrors.Error {
	return New(http.StatusTooManyRequests, reason, message)
}

// TaggedError is an error that carries a specific ErrorReason.
// This allows for explicit mapping of generic errors to predefined reasons.
type TaggedError struct {
	Err    error
	Reason apierrors.ErrorReason
}

func (e *TaggedError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Reason.String() // Fallback if no wrapped error
}

func (e *TaggedError) Unwrap() error {
	return e.Err
}

// WithReason tags a generic error with a specific ErrorReason.
// This allows Convert to map it to a specific Kratos error type.
func WithReason(err error, reason apierrors.ErrorReason) error {
	if err == nil {
		return nil
	}
	return &TaggedError{Err: err, Reason: reason}
}

// FromReason creates a Kratos error from a predefined error reason from the .proto file.
// This is the primary and consistent way to create standard application errors.
func FromReason(reason apierrors.ErrorReason) *kerrors.Error {
	// The message is a generic default. It can be overridden by using WithMessage().
	switch reason {
	// --- General --- 
	case apierrors.ErrorReason_VALIDATION_ERROR: 
		return BadRequest(reason.String(), "Request validation failed")
	case apierrors.ErrorReason_NOT_FOUND:
		return NotFound(reason.String(), "Resource not found")
	case apierrors.ErrorReason_INTERNAL_SERVER_ERROR:
		return InternalServer(reason.String(), "Internal server error")
	case apierrors.ErrorReason_METHOD_NOT_ALLOWED:
		return MethodNotAllowed(reason.String(), "Method not allowed")
	case apierrors.ErrorReason_REQUEST_TIMEOUT:
		return RequestTimeout(reason.String(), "Request timeout")
	case apierrors.ErrorReason_CONFLICT:
		return Conflict(reason.String(), "Resource conflict")
	case apierrors.ErrorReason_TOO_MANY_REQUESTS:
		return TooManyRequests(reason.String(), "Too many requests")
	case apierrors.ErrorReason_SERVICE_UNAVAILABLE:
		return ServiceUnavailable(reason.String(), "Service unavailable")
	case apierrors.ErrorReason_GATEWAY_TIMEOUT:
		return GatewayTimeout(reason.String(), "Gateway timeout")

	// --- Auth --- 
	case apierrors.ErrorReason_UNAUTHENTICATED, apierrors.ErrorReason_INVALID_CREDENTIALS, apierrors.ErrorReason_TOKEN_EXPIRED, apierrors.ErrorReason_TOKEN_INVALID, apierrors.ErrorReason_TOKEN_MISSING:
		return Unauthorized(reason.String(), "Authentication error")
	case apierrors.ErrorReason_FORBIDDEN:
		return Forbidden(reason.String(), "Permission denied")

	// --- Database --- 
	case apierrors.ErrorReason_DATABASE_ERROR:
		return InternalServer(reason.String(), "Database error")
	case apierrors.ErrorReason_RECORD_NOT_FOUND:
		return NotFound(reason.String(), "Record not found")
	case apierrors.ErrorReason_CONSTRAINT_VIOLATION, apierrors.ErrorReason_DUPLICATE_KEY:
		return Conflict(reason.String(), "Database constraint violation")
	case apierrors.ErrorReason_DATABASE_CONNECTION_FAILED:
		return ServiceUnavailable(reason.String(), "Database connection failed")

	// --- Business --- 
	case apierrors.ErrorReason_INVALID_STATE, apierrors.ErrorReason_MISSING_PARAMETER, apierrors.ErrorReason_INVALID_PARAMETER:
		return BadRequest(reason.String(), "Invalid business parameter")
	case apierrors.ErrorReason_RESOURCE_EXISTS, apierrors.ErrorReason_RESOURCE_IN_USE, apierrors.ErrorReason_ABORTED:
		return Conflict(reason.String(), "Business resource conflict")
	case apierrors.ErrorReason_CANCELLED:
		return ClientClosed(reason.String(), "Operation was cancelled")
	case apierrors.ErrorReason_OPERATION_NOT_ALLOWED:
		return Forbidden(reason.String(), "Operation not allowed")

	default:
		return InternalServer(apierrors.ErrorReason_UNKNOWN_ERROR.String(), "An unknown error occurred")
	}
}

// Convert takes any standard Go error and converts it into a structured Kratos error.
func Convert(err error) *kerrors.Error {
	if err == nil {
		return nil
	}

	var ke *kerrors.Error
	if errors.As(err, &ke) {
		return ke
	}

	// Check if the error is a TaggedError, allowing explicit mapping.
	var taggedErr *TaggedError
	if errors.As(err, &taggedErr) {
		return WithMessage(FromReason(taggedErr.Reason), taggedErr.Error())
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return FromReason(apierrors.ErrorReason_RECORD_NOT_FOUND)
	case errors.Is(err, context.DeadlineExceeded):
		return FromReason(apierrors.ErrorReason_REQUEST_TIMEOUT)
	case errors.Is(err, context.Canceled):
		return FromReason(apierrors.ErrorReason_CANCELLED)
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
		return FromReason(apierrors.ErrorReason_VALIDATION_ERROR)
	}

	// For unknown errors, create a standard internal error but preserve the original message.
	return WithMessage(FromReason(apierrors.ErrorReason_INTERNAL_SERVER_ERROR), err.Error())
}

// NewErrorEncoder returns a new transhttp.EncodeErrorFunc for centralized logging and error conversion.
func NewErrorEncoder() transhttp.EncodeErrorFunc {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		if err == nil {
			return
		}

		ke := Convert(err)

		if tr, ok := transport.FromServerContext(r.Context()); ok {
			fields := []interface{}{
				"kind", "server", "operation", tr.Operation(),
				"code", ke.Code, "reason", ke.Reason, "error", ke.Message,
			}
			if len(ke.Metadata) > 0 {
				fields = append(fields, "metadata", ke.Metadata)
			}

			type stackTracer interface {
				StackTrace() tkerrors.StackTrace
			}
			var st stackTracer
			if errors.As(err, &st) {
				// Limit stack trace to first 2 frames for brevity in logs
				fields = append(fields, "stack", fmt.Sprintf("%+v", st.StackTrace()[0:2]))
			}

			log.Context(r.Context()).Errorw(fields...)
		}

		transhttp.DefaultErrorEncoder(w, r, ke)
	}
}
