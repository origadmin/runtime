/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package application implements the functions, types, and interfaces for the module.
package application

import (
	"time"
)

type Application struct {
	ID          string
	Version     string
	ServiceName string
	StartTime   time.Time
	Metadata    map[string]string
}

func New(id, version, serviceName string, startTime time.Time, metadata map[string]string) *Application {
	return &Application{
		ID:          id,
		Version:     version,
		ServiceName: serviceName,
		StartTime:   startTime,
		Metadata:    metadata,
	}
}
