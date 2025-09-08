/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package errors

import (
	"errors"
	"fmt"

	kerrors "github.com/go-kratos/kratos/v2/errors"
)

// WithMessage enhances an error with a more specific message.
// It converts the error to a Kratos error first, then sets the message.
func WithMessage(err error, format string, args ...interface{}) *kerrors.Error {
	ke := Convert(err)
	ke.Message = fmt.Sprintf(format, args...)
	return ke
}

// WithMeta adds structured metadata to an error.
// It converts the error to a Kratos error first, then adds the metadata.
func WithMeta(err error, meta map[string]interface{}) *kerrors.Error {
	ke := Convert(err)
	if ke.Metadata == nil {
		ke.Metadata = make(map[string]string)
	}
	for k, v := range meta {
		ke.Metadata[k] = fmt.Sprintf("%v", v) // Convert interface{} to string
	}
	return ke
}

// WithField adds a single key-value pair as metadata to an error.
// It converts the error to a Kratos error first, then adds the field.
func WithField(err error, key string, value interface{}) *kerrors.Error {
	return WithMeta(err, map[string]interface{}{key: value})
}

// LookupMeta retrieves a specific metadata value from a Kratos error.
// It returns the value and a boolean indicating if the key was found.
func LookupMeta(err error, key string) (string, bool) {
	var ke *kerrors.Error
	if !errors.As(err, &ke) || ke.Metadata == nil {
		return "", false
	}
	val, ok := ke.Metadata[key]
	return val, ok
}

// --- Metadata Classification Helpers ---

// WithErrorOrigin adds the 'error_origin' metadata to an error.
// Example: WithErrorOrigin(err, "database"), WithErrorOrigin(err, "network")
func WithErrorOrigin(err error, origin string) *kerrors.Error {
	return WithField(err, "error_origin", origin)
}

// WithRecoverability adds the 'recoverability' metadata to an error.
// Example: WithRecoverability(err, "retriable"), WithRecoverability(err, "non_retriable")
func WithRecoverability(err error, status string) *kerrors.Error {
	return WithField(err, "recoverability", status)
}

// WithImpact adds the 'impact' metadata to an error.
// Example: WithImpact(err, "fatal"), WithImpact(err, "critical"), WithImpact(err, "minor")
func WithImpact(err error, impact string) *kerrors.Error {
	return WithField(err, "impact", impact)
}

// WithAudience adds the 'audience' metadata to an error.
// Example: WithAudience(err, "user"), WithAudience(err, "developer"), WithAudience(err, "operator")
func WithAudience(err error, audience string) *kerrors.Error {
	return WithField(err, "audience", audience)
}

// WithBusinessDomain adds the 'business_domain' metadata to an error.
// Example: WithBusinessDomain(err, "user_management"), WithBusinessDomain(err, "order_processing")
func WithBusinessDomain(err error, domain string) *kerrors.Error {
	return WithField(err, "business_domain", domain)
}
