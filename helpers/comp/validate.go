/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package comp

import (
	"github.com/origadmin/runtime/contracts/component"
)

// IsReserved checks if the metadata string is system-reserved.
func IsReserved(s string) bool {
	return len(s) > 0 && s[0] == '_'
}

// IsValidIdentifier checks if the string is a valid identifier for Category, Scope or Tag.
func IsValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	// Forbidden characters that might conflict with internal key generation
	for _, r := range s {
		if r == ' ' || r == ':' || r == '@' || r == ',' {
			return false
		}
	}
	return true
}

// ValidateScope checks if the scope is valid for the given locator.
func ValidateScope(l component.Locator, s component.Scope) bool {
	known := l.Scopes()
	for _, k := range known {
		if k == s {
			return true
		}
	}
	return false
}
