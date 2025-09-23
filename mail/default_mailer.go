/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package mail

import (
	"fmt"

	mailv1 "github.com/origadmin/runtime/api/gen/go/mail/v1"
	"github.com/origadmin/runtime/interfaces"
)

// defaultMailer is a default implementation of the interfaces.Mailer.
type defaultMailer struct {
	// Add fields for configuration if needed, e.g., SMTP server details
	cfg *mailv1.Mail
}

// NewDefaultMailer creates a new instance of the default Mailer.
func NewDefaultMailer(cfg *mailv1.Mail) interfaces.Mailer {
	return &defaultMailer{
		cfg: cfg,
	}
}

// Send implements the Mailer interface for defaultMailer.
func (m *defaultMailer) Send(msg *interfaces.Message) error {
	// This is a placeholder implementation.
	// In a real application, this would contain logic to send email,
	// e.g., using an SMTP client based on m.cfg.
	fmt.Printf("Default Mailer: Sending email from %s to %v with subject '%s'\n", msg.From, msg.To, msg.Subject)
	fmt.Printf("Body: %s (HTML: %t)\n", msg.Body, msg.HTML)
	if m.cfg != nil {
		fmt.Printf("Mailer configured with: %+v\n", m.cfg)
	}
	return nil
}
