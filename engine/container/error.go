/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package container

import (
	"fmt"
	"strings"

	"github.com/origadmin/runtime/contracts/component"
	"github.com/origadmin/toolkits/errors"
)

// EngineError represents a structured error from the engine container.
type EngineError struct {
	Op       string
	Category component.Category
	Scope    component.Scope
	Name     string
	Tags     []string
	Message  string
	Err      error
}

// Error implements the error interface.
func (e *EngineError) Error() string {
	var sb strings.Builder
	sb.WriteString("engine: ")
	if e.Op != "" {
		sb.WriteString(e.Op)
		sb.WriteString(" failed")
	}
	if e.Message != "" {
		sb.WriteString(": ")
		sb.WriteString(e.Message)
	}
	sb.WriteString(" [")
	fields := []string{}
	if e.Category != "" {
		fields = append(fields, fmt.Sprintf("category=%s", e.Category))
	}
	if e.Scope != "" {
		fields = append(fields, fmt.Sprintf("scope=%s", e.Scope))
	}
	if e.Name != "" {
		fields = append(fields, fmt.Sprintf("name=%s", e.Name))
	}
	if len(e.Tags) > 0 {
		fields = append(fields, fmt.Sprintf("tags=%v", e.Tags))
	}
	sb.WriteString(strings.Join(fields, ", "))
	sb.WriteString("]")
	if e.Err != nil {
		sb.WriteString(" | cause: ")
		sb.WriteString(e.Err.Error())
	}
	return sb.String()
}

// Unwrap returns the underlying error.
func (e *EngineError) Unwrap() error {
	return e.Err
}

// newErrorf creates a new EngineError with stack trace.
func newErrorf(op string, cat component.Category, scope component.Scope, name string, tags []string, format string, a ...interface{}) error {
	e := &EngineError{
		Op:       op,
		Category: cat,
		Scope:    scope,
		Name:     name,
		Tags:     tags,
		Message:  fmt.Sprintf(format, a...),
	}
	return errors.WithStack(e)
}

// wrapErrorf wraps an existing error with EngineError and stack trace.
func wrapErrorf(err error, op string, cat component.Category, scope component.Scope, name string, tags []string, format string, a ...interface{}) error {
	e := &EngineError{
		Op:       op,
		Category: cat,
		Scope:    scope,
		Name:     name,
		Tags:     tags,
		Message:  fmt.Sprintf(format, a...),
		Err:      err,
	}
	return errors.WithStack(e)
}
