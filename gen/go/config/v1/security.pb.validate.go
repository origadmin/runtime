// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: config/v1/security.proto

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

// Validate checks the field values on Security with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Security) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Security with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SecurityMultiError, or nil
// if none found.
func (m *Security) ValidateAll() error {
	return m.validate(true)
}

func (m *Security) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetJwt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SecurityValidationError{
					field:  "Jwt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SecurityValidationError{
					field:  "Jwt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetJwt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SecurityValidationError{
				field:  "Jwt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetCasbin()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SecurityValidationError{
					field:  "Casbin",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SecurityValidationError{
					field:  "Casbin",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCasbin()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SecurityValidationError{
				field:  "Casbin",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return SecurityMultiError(errors)
	}

	return nil
}

// SecurityMultiError is an error wrapping multiple validation errors returned
// by Security.ValidateAll() if the designated constraints aren't met.
type SecurityMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SecurityMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SecurityMultiError) AllErrors() []error { return m }

// SecurityValidationError is the validation error returned by
// Security.Validate if the designated constraints aren't met.
type SecurityValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SecurityValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SecurityValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SecurityValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SecurityValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SecurityValidationError) ErrorName() string { return "SecurityValidationError" }

// Error satisfies the builtin error interface
func (e SecurityValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSecurity.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SecurityValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SecurityValidationError{}

// Validate checks the field values on Security_CasbinConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *Security_CasbinConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Security_CasbinConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// Security_CasbinConfigMultiError, or nil if none found.
func (m *Security_CasbinConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *Security_CasbinConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Disabled

	// no validation rules for PolicyFile

	// no validation rules for ModelFile

	if len(errors) > 0 {
		return Security_CasbinConfigMultiError(errors)
	}

	return nil
}

// Security_CasbinConfigMultiError is an error wrapping multiple validation
// errors returned by Security_CasbinConfig.ValidateAll() if the designated
// constraints aren't met.
type Security_CasbinConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Security_CasbinConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Security_CasbinConfigMultiError) AllErrors() []error { return m }

// Security_CasbinConfigValidationError is the validation error returned by
// Security_CasbinConfig.Validate if the designated constraints aren't met.
type Security_CasbinConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Security_CasbinConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Security_CasbinConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Security_CasbinConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Security_CasbinConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Security_CasbinConfigValidationError) ErrorName() string {
	return "Security_CasbinConfigValidationError"
}

// Error satisfies the builtin error interface
func (e Security_CasbinConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSecurity_CasbinConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Security_CasbinConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Security_CasbinConfigValidationError{}

// Validate checks the field values on Security_JWTConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *Security_JWTConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Security_JWTConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// Security_JWTConfigMultiError, or nil if none found.
func (m *Security_JWTConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *Security_JWTConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Disabled

	// no validation rules for SigningMethod

	// no validation rules for SigningKey

	// no validation rules for OldSigningKey

	if all {
		switch v := interface{}(m.GetExpireTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, Security_JWTConfigValidationError{
					field:  "ExpireTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, Security_JWTConfigValidationError{
					field:  "ExpireTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetExpireTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return Security_JWTConfigValidationError{
				field:  "ExpireTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetRefreshTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, Security_JWTConfigValidationError{
					field:  "RefreshTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, Security_JWTConfigValidationError{
					field:  "RefreshTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRefreshTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return Security_JWTConfigValidationError{
				field:  "RefreshTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for CacheName

	if len(errors) > 0 {
		return Security_JWTConfigMultiError(errors)
	}

	return nil
}

// Security_JWTConfigMultiError is an error wrapping multiple validation errors
// returned by Security_JWTConfig.ValidateAll() if the designated constraints
// aren't met.
type Security_JWTConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m Security_JWTConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m Security_JWTConfigMultiError) AllErrors() []error { return m }

// Security_JWTConfigValidationError is the validation error returned by
// Security_JWTConfig.Validate if the designated constraints aren't met.
type Security_JWTConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Security_JWTConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Security_JWTConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Security_JWTConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Security_JWTConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Security_JWTConfigValidationError) ErrorName() string {
	return "Security_JWTConfigValidationError"
}

// Error satisfies the builtin error interface
func (e Security_JWTConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSecurity_JWTConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Security_JWTConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Security_JWTConfigValidationError{}
