/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package interfaces

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

// Mailer is the interface that wraps the basic Send method.
// Implementations of this interface are responsible for sending emails.
type Mailer interface {
	Send(msg *Message) error
}
