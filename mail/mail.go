/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package mail implements the functions, types, and interfaces for the module.
package mail

import (
	"sync"

	"github.com/origadmin/toolkits/mail"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
)

type (
	Builder = func(cfg *configv1.Mail) mail.Sender
	Sender  = mail.Sender
)

var (
	builder Builder
	sender  mail.Sender
	once    = &sync.Once{}
)

// Register registers a mail sender builder.
func Register(b Builder) {
	if builder != nil {
		panic("mail: Register called twice")
	}
	builder = b
}

// New returns a new mail sender.
func New(cfg *configv1.Mail) mail.Sender {
	once.Do(func() {
		if sender == nil {
			sender = builder(cfg)
		}
	})
	return sender
}
