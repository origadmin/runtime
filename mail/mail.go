/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package mail implements the functions, types, and interfaces for the module.
package mail

import (
	"fmt"
	"sync"

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
	"github.com/origadmin/runtime/container" // Import the new container package
)

// Message represents a generic email message.
type Message struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
	HTML    bool
	// Add other common fields like attachments, headers if needed
}

// Sender is the interface that wraps the basic Send method.
// Implementations of this interface are responsible for sending emails.
type Sender interface {
	Send(msg *Message) error
}

type (
	// Builder is a function type that takes a Mail configuration and returns a Sender.
	// This allows external packages to provide their specific Sender implementations.
	Builder = func(cfg *configv1.Mail) Sender
)

// mailInitOnce ensures that the mail sender is initialized only once.
var mailInitOnce sync.Once

// Init initializes the mail sender and registers it with the global container.
// This function should be called during the application's bootstrap phase.
func Init(cfg *configv1.Mail, mailBuilder Builder) {
	mailInitOnce.Do(func() {
		if mailBuilder == nil {
			panic("mail: mailBuilder cannot be nil during initialization")
		}
		senderInstance := mailBuilder(cfg)
		container.GlobalContainer.Register("mail.Sender", senderInstance)
	})
}

// GetSender retrieves the mail sender from the global container.
// It panics if the mail sender has not been initialized or registered.
func GetSender() Sender {
	cap, ok := container.GlobalContainer.Get("mail.Sender")
	if !ok {
		panic("mail: mail.Sender not initialized or registered in container")
	}
	sender, ok := cap.(Sender)
	if !ok {
		panic(fmt.Sprintf("mail: registered capability 'mail.Sender' is not of type Sender, got %T", cap))
	}
	return sender
}
