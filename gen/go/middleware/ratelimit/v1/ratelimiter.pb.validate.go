// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: middleware/ratelimit/v1/ratelimiter.proto

package ratelimitv1

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

// Validate checks the field values on RateLimiter with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *RateLimiter) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in RateLimiterMultiError, or
// nil if none found.
func (m *RateLimiter) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Enabled

	if _, ok := _RateLimiter_Name_InLookup[m.GetName()]; !ok {
		err := RateLimiterValidationError{
			field:  "Name",
			reason: "value must be in list [bbr memory redis]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Period

	// no validation rules for XRatelimitLimit

	// no validation rules for XRatelimitRemaining

	// no validation rules for XRatelimitReset

	// no validation rules for RetryAfter

	if all {
		switch v := interface{}(m.GetMemory()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "Memory",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "Memory",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMemory()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiterValidationError{
				field:  "Memory",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetRedis()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "Redis",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiterValidationError{
					field:  "Redis",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRedis()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiterValidationError{
				field:  "Redis",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return RateLimiterMultiError(errors)
	}

	return nil
}

// RateLimiterMultiError is an error wrapping multiple validation errors
// returned by RateLimiter.ValidateAll() if the designated constraints aren't met.
type RateLimiterMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiterMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiterMultiError) AllErrors() []error { return m }

// RateLimiterValidationError is the validation error returned by
// RateLimiter.Validate if the designated constraints aren't met.
type RateLimiterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiterValidationError) ErrorName() string { return "RateLimiterValidationError" }

// Error satisfies the builtin error interface
func (e RateLimiterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiterValidationError{}

var _RateLimiter_Name_InLookup = map[string]struct{}{
	"bbr":    {},
	"memory": {},
	"redis":  {},
}

// Validate checks the field values on RateLimiter_Redis with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RateLimiter_Redis) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter_Redis with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RateLimiter_RedisMultiError, or nil if none found.
func (m *RateLimiter_Redis) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter_Redis) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Addr

	// no validation rules for Username

	// no validation rules for Password

	// no validation rules for Db

	if len(errors) > 0 {
		return RateLimiter_RedisMultiError(errors)
	}

	return nil
}

// RateLimiter_RedisMultiError is an error wrapping multiple validation errors
// returned by RateLimiter_Redis.ValidateAll() if the designated constraints
// aren't met.
type RateLimiter_RedisMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiter_RedisMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiter_RedisMultiError) AllErrors() []error { return m }

// RateLimiter_RedisValidationError is the validation error returned by
// RateLimiter_Redis.Validate if the designated constraints aren't met.
type RateLimiter_RedisValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiter_RedisValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiter_RedisValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiter_RedisValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiter_RedisValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiter_RedisValidationError) ErrorName() string {
	return "RateLimiter_RedisValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimiter_RedisValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter_Redis.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiter_RedisValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiter_RedisValidationError{}

// Validate checks the field values on RateLimiter_Memory with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RateLimiter_Memory) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimiter_Memory with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RateLimiter_MemoryMultiError, or nil if none found.
func (m *RateLimiter_Memory) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimiter_Memory) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetExpiration()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_MemoryValidationError{
					field:  "Expiration",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_MemoryValidationError{
					field:  "Expiration",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetExpiration()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_MemoryValidationError{
				field:  "Expiration",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetCleanupInterval()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimiter_MemoryValidationError{
					field:  "CleanupInterval",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimiter_MemoryValidationError{
					field:  "CleanupInterval",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCleanupInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimiter_MemoryValidationError{
				field:  "CleanupInterval",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return RateLimiter_MemoryMultiError(errors)
	}

	return nil
}

// RateLimiter_MemoryMultiError is an error wrapping multiple validation errors
// returned by RateLimiter_Memory.ValidateAll() if the designated constraints
// aren't met.
type RateLimiter_MemoryMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimiter_MemoryMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimiter_MemoryMultiError) AllErrors() []error { return m }

// RateLimiter_MemoryValidationError is the validation error returned by
// RateLimiter_Memory.Validate if the designated constraints aren't met.
type RateLimiter_MemoryValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimiter_MemoryValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimiter_MemoryValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimiter_MemoryValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimiter_MemoryValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimiter_MemoryValidationError) ErrorName() string {
	return "RateLimiter_MemoryValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimiter_MemoryValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimiter_Memory.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimiter_MemoryValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimiter_MemoryValidationError{}