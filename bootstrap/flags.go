/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"crypto/rand"
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

var (
	RandomSuffix = fmt.Sprintf("%08d", time.Now().UnixNano()%(1<<32))
)

// ServiceID returns the ID of the service
func (f Flags) ServiceID() string {
	return f.ServiceName + "." + f.ID
}

// DefaultFlags returns the default flags for the service
func DefaultFlags() Flags {
	return NewFlags(DefaultServiceName, DefaultVersion)
}

// NewFlags returns a new set of flags for the service
func NewFlags(name string, version string) Flags {
	return Flags{
		ID:          RandomID(),
		Version:     version,
		ServiceName: name,
		StartTime:   time.Now(),
		Metadata:    make(map[string]string),
	}
}

func RandomID() string {
	id, err := os.Hostname()
	if err != nil {
		id = "unknown"
	}

	b := make([]byte, 4)
	if _, err := rand.Read(b); err == nil {
		return fmt.Sprintf("%s.%x", id, b)
	}
	return id + "." + RandomSuffix
}
