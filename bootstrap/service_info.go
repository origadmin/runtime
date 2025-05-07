/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package bootstrap

import (
	"fmt"
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
