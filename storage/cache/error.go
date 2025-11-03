/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package cache

import (
	"errors"
)

type cacheError struct {
	msg string
	err error
}

func (c *cacheError) Error() string {
	if c.err == nil {
		return c.msg
	}
	return c.msg + ": " + c.err.Error()
}

func (c *cacheError) Unwrap() error {
	return c.err
}

func (c *cacheError) Is(err error) bool {
	if err == nil {
		return c == nil
	}

	var cacheError *cacheError
	ok := errors.As(err, &cacheError)
	if !ok {
		return false
	}
	if cacheError.msg != c.msg {
		return errors.Is(err, c.err)
	}
	return ok
}

func NewError(msg string) error {
	return &cacheError{msg: msg}
}

func NewErrorWith(msg string, err error) error {
	return &cacheError{msg: msg, err: err}
}
