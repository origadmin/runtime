/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package configutil

import (
	"github.com/origadmin/runtime/contracts"
)

// ExtractName attempts to get a name from an item using common interfaces.
// It prioritizes the Named interface and falls back to ExtractType.
func ExtractName(item any) string {
	if item == nil {
		return ""
	}
	// Priority 1: Named interface
	if n, ok := item.(contracts.Named); ok {
		if name := n.GetName(); name != "" {
			return name
		}
	}
	// Priority 2: Fallback to Type extraction
	return ExtractType(item)
}

// ExtractType attempts to get a type/identity from an item using common interfaces.
// It aggregates Typed, Dialectal, and Driver interfaces.
func ExtractType(item any) string {
	if item == nil {
		return ""
	}
	if t, ok := item.(contracts.Typed); ok {
		if name := t.GetType(); name != "" {
			return name
		}
	}
	if d, ok := item.(contracts.Dialectal); ok {
		if name := d.GetDialect(); name != "" {
			return name
		}
	}
	if d, ok := item.(contracts.Driver); ok {
		if name := d.GetDriver(); name != "" {
			return name
		}
	}
	return ""
}
