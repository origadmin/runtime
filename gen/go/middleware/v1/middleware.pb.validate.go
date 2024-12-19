// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: middleware/v1/middleware.proto

package middlewarev1

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

// Validate checks the field values on Middleware with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Middleware) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Middleware with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MiddlewareMultiError, or
// nil if none found.
func (m *Middleware) ValidateAll() error {
	return m.validate(true)
}

func (m *Middleware) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Logging

	// no validation rules for Recovery

	// no validation rules for Tracing

	// no validation rules for CircuitBreaker

	if all {
		switch v := interface{}(m.GetMetadata()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewareValidationError{
				field:  "Metadata",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetRateLimiter()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "RateLimiter",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "RateLimiter",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRateLimiter()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewareValidationError{
				field:  "RateLimiter",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMetrics()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Metrics",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Metrics",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMetrics()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewareValidationError{
				field:  "Metrics",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetValidator()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Validator",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Validator",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetValidator()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewareValidationError{
				field:  "Validator",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetJwt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Jwt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Jwt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetJwt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewareValidationError{
				field:  "Jwt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetSelector()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Selector",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MiddlewareValidationError{
					field:  "Selector",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetSelector()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MiddlewareValidationError{
				field:  "Selector",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return MiddlewareMultiError(errors)
	}

	return nil
}

// MiddlewareMultiError is an error wrapping multiple validation errors
// returned by Middleware.ValidateAll() if the designated constraints aren't met.
type MiddlewareMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MiddlewareMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MiddlewareMultiError) AllErrors() []error { return m }

// MiddlewareValidationError is the validation error returned by
// Middleware.Validate if the designated constraints aren't met.
type MiddlewareValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MiddlewareValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MiddlewareValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MiddlewareValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MiddlewareValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MiddlewareValidationError) ErrorName() string { return "MiddlewareValidationError" }

// Error satisfies the builtin error interface
func (e MiddlewareValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMiddleware.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MiddlewareValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MiddlewareValidationError{}

// Validate checks the field values on Middleware_Metadata with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *Middleware_Metadata) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Middleware_Metadata with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// Middleware_MetadataMultiError, or nil if none found.
func (m *Middleware_Metadata) ValidateAll() error {
	return m.validate(true)
}

func (m *Middleware_Metadata) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Enabled

	// no validation rules for Prefix

	// no validation rules for Data

	if len(errors) > 0 {
		return Middleware_MetadataMultiError(errors)
	}

	return nil
}

// Middleware_MetadataMultiError is an error wrapping multiple validation
// errors returned by Middleware_Metadata.ValidateAll() if the designated
// constraints aren't met.
type Middleware_MetadataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Middleware_MetadataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Middleware_MetadataMultiError) AllErrors() []error { return m }

// Middleware_MetadataValidationError is the validation error returned by
// Middleware_Metadata.Validate if the designated constraints aren't met.
type Middleware_MetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Middleware_MetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Middleware_MetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Middleware_MetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Middleware_MetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Middleware_MetadataValidationError) ErrorName() string {
	return "Middleware_MetadataValidationError"
}

// Error satisfies the builtin error interface
func (e Middleware_MetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMiddleware_Metadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Middleware_MetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Middleware_MetadataValidationError{}
