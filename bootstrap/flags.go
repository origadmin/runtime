/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"fmt"
	"os"
	"time"
)

const (
	DefaultServiceName = "origadmin.service.v1"
	DefaultVersion     = "v1.0.0"
)

// Flags is a struct that holds the flags for the service
type Flags struct {
	ID          string
	Version     string
	ServiceName string
	StartTime   time.Time
	Metadata    map[string]string
}

// ServiceID returns the ID of the service
func (f Flags) ServiceID() string {
	return f.ID + "." + f.ServiceName
}

// DefaultFlags returns the default flags for the service
func DefaultFlags() Flags {
	return NewFlags(DefaultServiceName, DefaultVersion)
}

// NewFlags returns a new set of flags for the service
func NewFlags(name string, version string) Flags {
	return Flags{
		ID:          RandomTimeID(),
		Version:     version,
		ServiceName: name,
		StartTime:   time.Now(),
		Metadata:    make(map[string]string),
	}
}

func RandomTimeID() string {
	id, err := os.Hostname()
	if err != nil {
		id = "unknown"
	}
	return id + "." + fmt.Sprintf("%08d", time.Now().Unix()%(1<<32))
}
