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

// ServiceInfo is a struct that holds the flags for the service
type ServiceInfo struct {
	ID        string
	Name      string
	Version   string
	StartTime time.Time
	Metadata  map[string]string
}

var (
	RandomSuffix = fmt.Sprintf("%08d", time.Now().UnixNano()%(1<<32))
)

// ServiceID returns the ID of the service
func (si ServiceInfo) ServiceID() string {
	return si.Name + "." + si.ID
}

func NewServiceInfo(name, version string) ServiceInfo {
	return ServiceInfo{
		ID:      RandomID(),
		Name:    name,
		Version: version,
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
