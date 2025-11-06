/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package errors

import (
	"net/http"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	commonv1 "github.com/origadmin/runtime/api/gen/go/config/common/v1"
)

// ToKratos converts an internal Structured error to a Kratos error for API responses.
// This function resides in a dedicated API adapter to keep internal errors decoupled from API logic.
func ToKratos(e *Structured, reason commonv1.ErrorReason) *kerrors.Error {
	if e == nil {
		return nil
	}

	// Prefer explicit reason when provided
	if reason != commonv1.ErrorReason_UNKNOWN_ERROR {
		err := NewMessage(reason, "%s", e.Message) // Fixed: Pass e.Message as an argument to a format string
		// Attach metadata hints
		if err.Metadata == nil {
			err.Metadata = make(map[string]string)
		}
		if e.Module != "" {
			err.Metadata["module"] = e.Module
		}
		if e.Op != "" {
			err.Metadata["operation"] = e.Op
		}
		for k, v := range e.Metadata {
			err.Metadata[k] = v
		}
		return err
	}

	// Fallback: default to internal server error without relying on internal codes
	statusCode := http.StatusInternalServerError

	reasonStr := e.Module
	if reasonStr == "" {
		reasonStr = "INTERNAL_ERROR"
	}

	kratosErr := kerrors.New(statusCode, reasonStr, e.Message)
	if kratosErr.Metadata == nil {
		kratosErr.Metadata = make(map[string]string)
	}
	if e.Module != "" {
		kratosErr.Metadata["module"] = e.Module
	}
	if e.Op != "" {
		kratosErr.Metadata["operation"] = e.Op
	}
	for k, v := range e.Metadata {
		kratosErr.Metadata[k] = v
	}
	return kratosErr
}
