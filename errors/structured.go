/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Structured represents an error with module context and structured metadata.
// This type is primarily used for internal module error handling.
type Structured struct {
	// Module indicates which module the error occurred in (e.g., "auth", "db")
	Module string `json:"module"`
	// Message is the human-readable error message
	Message string `json:"message"`
	// Op is the operation that caused the error (usually the function name)
	Op string `json:"operation,omitempty"`
	// Metadata contains additional context about the error
	Metadata map[string]string `json:"metadata,omitempty"`
	// Err is the underlying error
	Err error `json:"-"`
}

// NewStructured creates a new structured error
func NewStructured(module, format string, args ...interface{}) *Structured {
	return &Structured{
		Module:  module,
		Message: fmt.Sprintf(format, args...),
	}
}

// WithCaller adds operation context to the error by getting the caller function name
func (e *Structured) WithCaller() *Structured {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		e.Op = "unknown"
	} else {
		funcName := runtime.FuncForPC(pc).Name()
		e.Op = funcName
	}
	return e
}

// WithOperation sets a specific operation name for the error
func (e *Structured) WithOperation(op string) *Structured {
	e.Op = op
	return e
}

// WithMetadata adds metadata to the error
func (e *Structured) WithMetadata(metadata map[string]string) *Structured {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	for k, v := range metadata {
		e.Metadata[k] = v
	}
	return e
}

// WithField adds a single key-value pair to the error metadata
func (e *Structured) WithField(key string, value string) *Structured {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
	return e
}

// Error returns the error message
func (e *Structured) Error() string {
	var sb strings.Builder

	// Add module if available
	if e.Module != "" {
		sb.WriteString("[")
		sb.WriteString(e.Module)
		sb.WriteString("] ")
	}

	// Add operation if available
	if e.Op != "" {
		sb.WriteString(e.Op)
		sb.WriteString(": ")
	}

	// Add error message
	sb.WriteString(e.Message)

	// Add underlying error if available
	if e.Err != nil {
		sb.WriteString(": ")
		sb.WriteString(e.Err.Error())
	}

	return sb.String()
}

// Unwrap returns the underlying error
func (e *Structured) Unwrap() error {
	return e.Err
}

// Is checks if the target error is of type *Structured and matches based on module
func (e *Structured) Is(target error) bool {
	var t *Structured
	ok := errors.As(target, &t)
	if !ok {
		return false
	}

	return e.Module == t.Module
}

// WrapStructured wraps an error with additional context
func WrapStructured(err error, module, format string, args ...interface{}) *Structured {
	if err == nil {
		return nil
	}
	return &Structured{
		Module:  module,
		Message: fmt.Sprintf(format, args...),
		Err:     err,
	}
}
