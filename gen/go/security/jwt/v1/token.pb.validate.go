// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: security/jwt/v1/token.proto

package jwtv1

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

// Validate checks the field values on Token with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Token) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Token with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in TokenMultiError, or nil if none found.
func (m *Token) ValidateAll() error {
	return m.validate(true)
}

func (m *Token) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ClientId

	if utf8.RuneCountInString(m.GetUserId()) < 1 {
		err := TokenValidationError{
			field:  "UserId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetAccessToken()) < 1 {
		err := TokenValidationError{
			field:  "AccessToken",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetRefreshToken()) < 1 {
		err := TokenValidationError{
			field:  "RefreshToken",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetExpirationTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TokenValidationError{
					field:  "ExpirationTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TokenValidationError{
					field:  "ExpirationTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetExpirationTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TokenValidationError{
				field:  "ExpirationTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return TokenMultiError(errors)
	}

	return nil
}

// TokenMultiError is an error wrapping multiple validation errors returned by
// Token.ValidateAll() if the designated constraints aren't met.
type TokenMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TokenMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TokenMultiError) AllErrors() []error { return m }

// TokenValidationError is the validation error returned by Token.Validate if
// the designated constraints aren't met.
type TokenValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TokenValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TokenValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TokenValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TokenValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TokenValidationError) ErrorName() string { return "TokenValidationError" }

// Error satisfies the builtin error interface
func (e TokenValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sToken.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TokenValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TokenValidationError{}
