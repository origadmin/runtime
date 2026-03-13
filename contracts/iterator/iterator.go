/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package iterator

// Iterator is a common interface for iterating over collections.
type Iterator interface {
	// Next moves the cursor to the next element and returns true if successful.
	Next() bool
	// Value returns the current element's key and value.
	Value() (string, any)
	// Err returns any error encountered during iteration.
	Err() error
}
