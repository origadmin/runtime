// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pwt/v1/pwt.proto

package pwtv1

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

// Validate checks the field values on PWT with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *PWT) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PWT with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PWTMultiError, or nil if none found.
func (m *PWT) ValidateAll() error {
	return m.validate(true)
}

func (m *PWT) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetExpirationTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PWTValidationError{
					field:  "ExpirationTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PWTValidationError{
					field:  "ExpirationTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetExpirationTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PWTValidationError{
				field:  "ExpirationTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetIssuedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PWTValidationError{
					field:  "IssuedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PWTValidationError{
					field:  "IssuedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetIssuedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PWTValidationError{
				field:  "IssuedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetNotBefore()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, PWTValidationError{
					field:  "NotBefore",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, PWTValidationError{
					field:  "NotBefore",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetNotBefore()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return PWTValidationError{
				field:  "NotBefore",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Issuer

	// no validation rules for Subject

	// no validation rules for JwtId

	// no validation rules for ClientId

	// no validation rules for ClientSecret

	if len(errors) > 0 {
		return PWTMultiError(errors)
	}

	return nil
}

// PWTMultiError is an error wrapping multiple validation errors returned by
// PWT.ValidateAll() if the designated constraints aren't met.
type PWTMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PWTMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PWTMultiError) AllErrors() []error { return m }

// PWTValidationError is the validation error returned by PWT.Validate if the
// designated constraints aren't met.
type PWTValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PWTValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PWTValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PWTValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PWTValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PWTValidationError) ErrorName() string { return "PWTValidationError" }

// Error satisfies the builtin error interface
func (e PWTValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPWT.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PWTValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PWTValidationError{}
