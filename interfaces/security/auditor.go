/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package security implements the functions, types, and interfaces for the module.
package security

import (
	"context"
	"time"
)

type Auditor interface {
	LogAuthEvent(ctx context.Context, event AuditorEvent) error
}

type AuditorEvent struct {
	Timestamp time.Time
	Subject   string
	Action    string
	Object    string
	Success   bool
	ClientIP  string
	UserAgent string
}
