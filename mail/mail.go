/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package mail implements the functions, types, and interfaces for the module.
package mail

import (
	"sync"

	mailv1 "github.com/origadmin/runtime/api/gen/go/mail/v1"
	"github.com/origadmin/runtime/interfaces"
)

// Builder is a function type that takes a Mail configuration and returns a Mailer.
// This allows external packages to provide their specific Mailer implementations.
type Builder = func(cfg *mailv1.Mail) interfaces.Mailer

// initOnce ensures that the mail mailer is initialized only once.
var initOnce sync.Once

// _mailer holds the initialized Mailer instance for the mail package.
var _mailer interfaces.Mailer

// Init initializes the mail mailer.
// This function should be called during the application's bootstrap phase.
// If mailBuilder is nil, the default mailer implementation will be used.
func Init(cfg *mailv1.Mail, mailBuilder Builder) {
	initOnce.Do(func() {
		if mailBuilder == nil {
			// If no builder is provided, use the default mailer implementation
			_mailer = NewDefaultMailer(cfg)
		} else {
			_mailer = mailBuilder(cfg)
		}
	})
}

// Mailer retrieves the mail mailer.
// It panics if the mail mailer has not been initialized.
func Mailer() interfaces.Mailer {
	if _mailer == nil {
		panic("mail: mail mailer not initialized. Call Init() first.")
	}
	return _mailer
}
