/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package mail implements the functions, types, and interfaces for the module.
package mail

import (
	"sync"

	// 移除对 github.com/origadmin/toolkits/mail 的直接导入，因为 Sender 接口将在此处定义
	// "github.com/origadmin/toolkits/mail" 

	configv1 "github.com/origadmin/runtime/api/gen/go/config/v1"
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
	Builder = func(cfg *configv1.Mail) Sender // Builder 现在返回此包定义的 Sender 接口
)

var (
	builder Builder
	sender  Sender // sender 现在是此包定义的 Sender 接口类型
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
func New(cfg *configv1.Mail) Sender { // New 现在返回此包定义的 Sender 接口
	once.Do(func() {
		if sender == nil {
			sender = builder(cfg)
		}
	})
	return sender
}
