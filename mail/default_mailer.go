/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"

	mailv1 "github.com/origadmin/runtime/api/gen/go/runtime/mail/v1"
	commonv1 "github.com/origadmin/runtime/api/gen/go/runtime/common/v1"
	runtimeerrors "github.com/origadmin/runtime/errors"
	"github.com/origadmin/runtime/interfaces"
)

const Module = "mail.defaultMailer"

// defaultMailer is a default implementation of the interfaces.Mailer.
type defaultMailer struct {
	cfg *mailv1.Mail
}

// NewDefaultMailer creates a new instance of the default Mailer.
func NewDefaultMailer(cfg *mailv1.Mail) interfaces.Mailer {
	return &defaultMailer{
		cfg: cfg,
	}
}

// Send implements the Mailer interface for defaultMailer using net/smtp.
func (m *defaultMailer) Send(msg *interfaces.Message) error {
	smtpConfig := m.cfg.GetSmtpConfig()
	if smtpConfig == nil {
		return runtimeerrors.NewStructured(Module, "SMTP configuration is not provided in the Mail config").WithReason(commonv1.ErrorReason_MISSING_PARAMETER).WithCaller()
	}

	if smtpConfig.Host == "" || smtpConfig.Port == 0 {
		return runtimeerrors.NewStructured(Module, "SMTP configuration (host, port) is missing").WithReason(commonv1.ErrorReason_MISSING_PARAMETER).WithCaller()
	}

	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)

	// Set up authentication information.
	var auth smtp.Auth
	if smtpConfig.Username != "" && smtpConfig.Password != "" {
		auth = smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)
	}

	// Prepare the message
	var b bytes.Buffer
	fmt.Fprintf(&b, "From: %s\r\n", msg.From)
	fmt.Fprintf(&b, "To: %s\r\n", strings.Join(msg.To, ", "))
	if len(msg.Cc) > 0 {
		fmt.Fprintf(&b, "Cc: %s\r\n", strings.Join(msg.Cc, ", "))
	}
	if len(msg.Bcc) > 0 {
		// Bcc is handled by the smtp.SendMail function, not in headers
	}
	fmt.Fprintf(&b, "Subject: %s\r\n", msg.Subject)

	if msg.HTML {
		fmt.Fprintf(&b, "MIME-version: 1.0;\r\n")
		fmt.Fprintf(&b, "Content-Type: text/html; charset=\"UTF-8\";\r\n")
	} else {
		fmt.Fprintf(&b, "Content-Type: text/plain; charset=\"UTF-8\";\r\n")
	}
	fmt.Fprintf(&b, "\r\n%s", msg.Body)

	// Collect all recipients (To, Cc, Bcc)
	recipients := make([]string, 0, len(msg.To)+len(msg.Cc)+len(msg.Bcc))
	recipients = append(recipients, msg.To...)
	recipients = append(recipients, msg.Cc...)
	recipients = append(recipients, msg.Bcc...)

	// Send the email
	err := smtp.SendMail(addr, auth, msg.From, recipients, b.Bytes())
	if err != nil {
		return runtimeerrors.WrapStructured(err, Module, "failed to send email").WithReason(commonv1.ErrorReason_INTERNAL_SERVER_ERROR).WithCaller()
	}

	fmt.Printf("Mail: Successfully sent email from %s to %v with subject '%s' via %s\n", msg.From, msg.To, msg.Subject, addr)
	return nil
}
