/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package application implements the functions, types, and interfaces for the module.
package application

import (
	"fmt"
	"os"
	"time"
)

var (
	timeSuffix = fmt.Sprintf("%08d", time.Now().UnixNano()%(1<<32))
)

type Application struct {
	ID        string
	Name      string
	Version   string
	StartTime time.Time
	Metadata  map[string]string
	genID     func() string
}

func (obj *Application) Init(name string, version string) {
	obj.ID = obj.getID()
	obj.Name = name
	obj.Version = version
	obj.StartTime = time.Now()
	obj.Metadata = make(map[string]string)
}

func (obj *Application) SetGenID(genID func() string) {
	obj.genID = genID
}

func (obj *Application) getID() string {
	if obj.genID != nil {
		obj.ID = obj.genID()
	}
	return RandomID()
}

func RandomID() string {
	id, err := os.Hostname()
	if err != nil {
		id = "unknown"
	}
	return id + "." + timeSuffix
}
