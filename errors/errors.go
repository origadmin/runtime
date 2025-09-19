/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package errors provides a centralized hub for handling, converting, and rendering errors.
// It uses a proto-defined enum for standard error reasons and adapts the Kratos error package.
package errors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"

	commonv1 "github.com/origadmin/runtime/api/gen/go/common/v1"
	"github.com/origadmin/runtime/log"
	tkerrors "github.com/origadmin/toolkits/errors"
)

// RequestTimeout creates a 408 Request Timeout error.
func RequestTimeout(reason, message string) *kerrors.Error {
	return kerrors.New(http.StatusRequestTimeout, reason, message)
}

// MethodNotAllowed creates a 405 Method Not Allowed error.
func MethodNotAllowed(reason, message string) *kerrors.Error {
	return kerrors.New(http.StatusMethodNotAllowed, reason, message)
}

// TooManyRequests creates a 429 Too Many Requests error.
func TooManyRequests(reason, message string) *kerrors.Error {
	return kerrors.New(http.StatusTooManyRequests, reason, message)
}

// TaggedError is an error that carries a specific ErrorReason.
// This allows for explicit mapping of generic errors to predefined reasons.
type TaggedError struct {
	Err    error
	Reason commonv1.ErrorReason
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
func WithReason(err error, reason commonv1.ErrorReason) error {
	if err == nil {
		return nil
	}
	return &TaggedError{Err: err, Reason: reason}
}

// FromReason creates a Kratos error from a predefined error reason from the .proto file.
// This is the primary and consistent way to create standard application errors.
func FromReason(reason commonv1.ErrorReason) *kerrors.Error {
	// The message is a generic default. It can be overridden by setting the Message field directly.
	switch reason {
	// --- General --- 
	case commonv1.ErrorReason_VALIDATION_ERROR:
		return kerrors.BadRequest(reason.String(), "Request validation failed")
	case commonv1.ErrorReason_NOT_FOUND:
		return kerrors.NotFound(reason.String(), "Resource not found")
	case commonv1.ErrorReason_INTERNAL_SERVER_ERROR:
		return kerrors.InternalServer(reason.String(), "Internal server error")
	case commonv1.ErrorReason_METHOD_NOT_ALLOWED:
		return MethodNotAllowed(reason.String(), "Method not allowed")
	case commonv1.ErrorReason_REQUEST_TIMEOUT:
		return RequestTimeout(reason.String(), "Request timeout")
	case commonv1.ErrorReason_CONFLICT:
		return kerrors.Conflict(reason.String(), "Resource conflict")
	case commonv1.ErrorReason_TOO_MANY_REQUESTS:
		return TooManyRequests(reason.String(), "Too many requests")
	case commonv1.ErrorReason_SERVICE_UNAVAILABLE:
		return kerrors.ServiceUnavailable(reason.String(), "Service unavailable")
	case commonv1.ErrorReason_GATEWAY_TIMEOUT:
		return kerrors.GatewayTimeout(reason.String(), "Gateway timeout")

	// --- Auth --- 
	case commonv1.ErrorReason_UNAUTHENTICATED:
		// TODO: These reasons should be defined and handled in the specific security/auth module
		// commonv1.ErrorReason_INVALID_CREDENTIALS,
		// commonv1.ErrorReason_TOKEN_EXPIRED,
		// commonv1.ErrorReason_TOKEN_INVALID,
		// commonv1.ErrorReason_TOKEN_MISSING:
		return kerrors.Unauthorized(reason.String(), "Authentication error")
	case commonv1.ErrorReason_FORBIDDEN:
		return kerrors.Forbidden(reason.String(), "Permission denied")

	// --- Database --- 
	case commonv1.ErrorReason_DATABASE_ERROR:
		return kerrors.InternalServer(reason.String(), "Database error")
	case commonv1.ErrorReason_RECORD_NOT_FOUND:
		return kerrors.NotFound(reason.String(), "Record not found")
	case commonv1.ErrorReason_CONSTRAINT_VIOLATION, commonv1.ErrorReason_DUPLICATE_KEY:
		return kerrors.Conflict(reason.String(), "Database constraint violation")
	case commonv1.ErrorReason_DATABASE_CONNECTION_FAILED:
		return kerrors.ServiceUnavailable(reason.String(), "Database connection failed")

	// --- Business --- 
	case commonv1.ErrorReason_INVALID_STATE, commonv1.ErrorReason_MISSING_PARAMETER, commonv1.ErrorReason_INVALID_PARAMETER:
		return kerrors.BadRequest(reason.String(), "Invalid business parameter")
	case commonv1.ErrorReason_RESOURCE_EXISTS, commonv1.ErrorReason_RESOURCE_IN_USE, commonv1.ErrorReason_ABORTED:
		return kerrors.Conflict(reason.String(), "Business resource conflict")
	case commonv1.ErrorReason_CANCELLED:
		return kerrors.ClientClosed(reason.String(), "Operation was cancelled")
	case commonv1.ErrorReason_OPERATION_NOT_ALLOWED:
		return kerrors.Forbidden(reason.String(), "Operation not allowed")

	// --- Registry Errors (6000-6999) ---
	case commonv1.ErrorReason_REGISTRY_NOT_FOUND:
		return kerrors.NotFound(reason.String(), "Registry entry not found")
	// TODO: These reasons should be defined and handled in the specific registry module
	//case commonv1.ErrorReason_INVALID_REGISTRY_CONFIG:
	//	return kerrors.BadRequest(reason.String(), "Invalid registry configuration")
	//case commonv1.ErrorReason_REGISTRY_CREATION_FAILURE:
	//	return kerrors.InternalServer(reason.String(), "Registry creation failed")

	default:
		return kerrors.InternalServer(commonv1.ErrorReason_UNKNOWN_ERROR.String(), "An unknown error occurred")
	}
}

// NewMessage creates a Kratos error from a predefined error reason, with a formatted message.
func NewMessage(reason commonv1.ErrorReason, format string, a ...interface{}) *kerrors.Error {
	err := FromReason(reason)
	err.Message = fmt.Sprintf(format, a...) // Directly set the message
	return err
}

// NewMessageWithMeta creates a Kratos error from a predefined error reason,
// with a formatted message and specified metadata.
func NewMessageWithMeta(reason commonv1.ErrorReason, metadata map[string]string, format string, a ...interface{}) *kerrors.Error {
	err := FromReason(reason)
	err.Message = fmt.Sprintf(format, a...) // Directly set the message
	if err.Metadata == nil {
		err.Metadata = make(map[string]string)
	}
	for k, v := range metadata {
		err.Metadata[k] = v // Directly set metadata
	}
	return err
}

// WrapAndConvert wraps an original error with a reason, converts it to a Kratos error,
// and sets a formatted message.
func WrapAndConvert(originalErr error, reason commonv1.ErrorReason, format string, a ...interface{}) *kerrors.Error {
	// 1. Wrap the original error with the specified reason
	taggedErr := WithReason(originalErr, reason)

	// 2. Convert the tagged error to a Kratos error
	convertedErr := Convert(taggedErr)

	// 3. Set the formatted message
	convertedErr.Message = fmt.Sprintf(format, a...)

	return convertedErr
}

// Convert takes any standard Go error and converts it into a structured Kratos error.
func Convert(err error) *kerrors.Error {
	if err == nil {
		return nil
	}

	// 1. Check if the error is a TaggedError (explicitly mapped by developer)
	var taggedErr *TaggedError
	if errors.As(err, &taggedErr) {
		// Use the tagged reason, and preserve the original error message
		ke := FromReason(taggedErr.Reason)
		ke.Message = fmt.Sprintf("%s", taggedErr.Error()) // Directly set the message
		return ke
	}

	// 2. Check if the error is already a Kratos error (from Kratos itself or a plugin)
	var ke *kerrors.Error
	if errors.As(err, &ke) {
		// Try to map the existing Kratos error's Reason to our commonv1.ErrorReason
		parsedReason, ok := commonv1.ErrorReason_value[ke.Reason]
		if ok {
			// If the reason matches one of our predefined reasons, use our FromReason
			// to ensure consistency (e.g., default message, code from our proto)
			// and merge metadata.
			newKe := FromReason(commonv1.ErrorReason(parsedReason))
			newKe.Message = fmt.Sprintf("%s", ke.Message) // Directly set the message
			if newKe.Metadata == nil {
				newKe.Metadata = make(map[string]string)
			}
			for k, v := range ke.Metadata {
				newKe.Metadata[k] = v // Directly set metadata
			}
			return newKe
		}
		// If the reason does not match our predefined reasons,
		// it's an external Kratos error with an unknown reason.
		// Convert it to a generic external service error, and put the original reason and code into metadata.
		newKe := FromReason(commonv1.ErrorReason_EXTERNAL_SERVICE_ERROR)
		newKe.Message = fmt.Sprintf("External Kratos error: %s", ke.Message) // Directly set the message
		if newKe.Metadata == nil {
			newKe.Metadata = make(map[string]string)
		}
		newKe.Metadata["original_reason"] = ke.Reason
		newKe.Metadata["original_code"] = fmt.Sprintf("%d", ke.Code)
		return newKe
	}

	// 3. Handle specific standard library errors (implicit mapping)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return FromReason(commonv1.ErrorReason_RECORD_NOT_FOUND)
	case errors.Is(err, context.DeadlineExceeded):
		return FromReason(commonv1.ErrorReason_REQUEST_TIMEOUT)
	case errors.Is(err, context.Canceled):
		return FromReason(commonv1.ErrorReason_CANCELLED)
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
		return FromReason(commonv1.ErrorReason_VALIDATION_ERROR)
	}

	// 4. Default to INTERNAL_SERVER_ERROR for any other unhandled error
	ke = FromReason(commonv1.ErrorReason_INTERNAL_SERVER_ERROR)
	ke.Message = fmt.Sprintf("%s", err.Error()) // Directly set the message
	return ke
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
