// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: config/v1/task.proto

package configv1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Task with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Task) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Task with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in TaskMultiError, or nil if none found.
func (m *Task) ValidateAll() error {
	return m.validate(true)
}

func (m *Task) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Type

	// no validation rules for Name

	if all {
		switch v := interface{}(m.GetAsynq()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TaskValidationError{
					field:  "Asynq",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TaskValidationError{
					field:  "Asynq",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAsynq()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TaskValidationError{
				field:  "Asynq",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMachinery()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TaskValidationError{
					field:  "Machinery",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TaskValidationError{
					field:  "Machinery",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMachinery()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TaskValidationError{
				field:  "Machinery",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetCron()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TaskValidationError{
					field:  "Cron",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TaskValidationError{
					field:  "Cron",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCron()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TaskValidationError{
				field:  "Cron",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return TaskMultiError(errors)
	}

	return nil
}

// TaskMultiError is an error wrapping multiple validation errors returned by
// Task.ValidateAll() if the designated constraints aren't met.
type TaskMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TaskMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TaskMultiError) AllErrors() []error { return m }

// TaskValidationError is the validation error returned by Task.Validate if the
// designated constraints aren't met.
type TaskValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TaskValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TaskValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TaskValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TaskValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TaskValidationError) ErrorName() string { return "TaskValidationError" }

// Error satisfies the builtin error interface
func (e TaskValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTask.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TaskValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TaskValidationError{}

// Validate checks the field values on Task_Asynq with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Task_Asynq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Task_Asynq with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in Task_AsynqMultiError, or
// nil if none found.
func (m *Task_Asynq) ValidateAll() error {
	return m.validate(true)
}

func (m *Task_Asynq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Endpoint

	// no validation rules for Password

	// no validation rules for Db

	// no validation rules for Location

	if len(errors) > 0 {
		return Task_AsynqMultiError(errors)
	}

	return nil
}

// Task_AsynqMultiError is an error wrapping multiple validation errors
// returned by Task_Asynq.ValidateAll() if the designated constraints aren't met.
type Task_AsynqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Task_AsynqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Task_AsynqMultiError) AllErrors() []error { return m }

// Task_AsynqValidationError is the validation error returned by
// Task_Asynq.Validate if the designated constraints aren't met.
type Task_AsynqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Task_AsynqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Task_AsynqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Task_AsynqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Task_AsynqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Task_AsynqValidationError) ErrorName() string { return "Task_AsynqValidationError" }

// Error satisfies the builtin error interface
func (e Task_AsynqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTask_Asynq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Task_AsynqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Task_AsynqValidationError{}

// Validate checks the field values on Task_Machinery with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Task_Machinery) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Task_Machinery with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in Task_MachineryMultiError,
// or nil if none found.
func (m *Task_Machinery) ValidateAll() error {
	return m.validate(true)
}

func (m *Task_Machinery) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return Task_MachineryMultiError(errors)
	}

	return nil
}

// Task_MachineryMultiError is an error wrapping multiple validation errors
// returned by Task_Machinery.ValidateAll() if the designated constraints
// aren't met.
type Task_MachineryMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Task_MachineryMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Task_MachineryMultiError) AllErrors() []error { return m }

// Task_MachineryValidationError is the validation error returned by
// Task_Machinery.Validate if the designated constraints aren't met.
type Task_MachineryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Task_MachineryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Task_MachineryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Task_MachineryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Task_MachineryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Task_MachineryValidationError) ErrorName() string { return "Task_MachineryValidationError" }

// Error satisfies the builtin error interface
func (e Task_MachineryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTask_Machinery.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Task_MachineryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Task_MachineryValidationError{}

// Validate checks the field values on Task_Cron with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Task_Cron) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Task_Cron with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in Task_CronMultiError, or nil
// if none found.
func (m *Task_Cron) ValidateAll() error {
	return m.validate(true)
}

func (m *Task_Cron) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Addr

	if len(errors) > 0 {
		return Task_CronMultiError(errors)
	}

	return nil
}

// Task_CronMultiError is an error wrapping multiple validation errors returned
// by Task_Cron.ValidateAll() if the designated constraints aren't met.
type Task_CronMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Task_CronMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Task_CronMultiError) AllErrors() []error { return m }

// Task_CronValidationError is the validation error returned by
// Task_Cron.Validate if the designated constraints aren't met.
type Task_CronValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Task_CronValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Task_CronValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Task_CronValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Task_CronValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Task_CronValidationError) ErrorName() string { return "Task_CronValidationError" }

// Error satisfies the builtin error interface
func (e Task_CronValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTask_Cron.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Task_CronValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Task_CronValidationError{}
