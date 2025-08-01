/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package queue

import (
	"context"
	"time"
)

// Iterator interface
type Iterator[E any] interface {
	// Next advances the iterator and returns true if there is a next element
	Next() bool
	// Value returns the current element
	Value() E
}

// Queue interface
type Queue[E any] interface {
	// Offer adds an element to the queue if possible, returning true on success
	Offer(E) bool
	// Poll retrieves and removes the head of the queue, or returns false if empty
	Poll() (E, bool)
	// Peek retrieves but does not remove the head of the queue, or returns false if empty
	Peek() (E, bool)
	// Size returns the number of elements in the queue
	Size() int64
	// IsEmpty returns true if the queue contains no elements
	IsEmpty() bool
	// Clear removes all elements from the queue
	Clear()
	// ToSlice returns a slice containing all the elements in the queue
	ToSlice() []E
	// Iterator returns an Iterator over the elements in this queue
	Iterator() Iterator[E]
}

// Provider interface
type Provider[E any] interface {
	Queue(topic string) Queue[E] // Get a queue by topic
	Publish(ctx context.Context, topic string, payload E) error
	Subscribe(ctx context.Context, topic string, handler func(E)) error
	Close() error
}

// Message represent a message
type Message struct {
	ID        string
	Topic     string
	Timestamp time.Time
	Payload   string
}

// DeliverMessage represents a delivered message interface
type DeliverMessage interface {
	// Message return to original message
	Message() *Message
	// Ack confirm that the message was successfully processed
	Ack(ctx context.Context) error
	// Nack rejecting a message may cause the message to rejoin the team or be discarded
	Nack(ctx context.Context) error
}

// MessageQueue is a message Queue interface that extends the basic queue interface
type MessageQueue interface {
	// Publish a message to the queue
	Publish(ctx context.Context, payload string) (string, error)
	// Subscribe to the queue and process the messages
	Subscribe(ctx context.Context, handler func(DeliverMessage)) error
	// Size returns the number of messages in the queue
	Size() int64
}

// MessageQueueProvider provides the ability to get MessageQueue by subject
type MessageQueueProvider interface {
	// MessageQueue obtain the corresponding MessageQueue based on subject
	MessageQueue(subject string) (MessageQueue, error)
	// PublishToTopic publishes a message to the specified topic
	PublishToTopic(ctx context.Context, topic string, payload string) error
	// SubscribeToTopic subscribe to messages for a specified topic
	SubscribeToTopic(ctx context.Context, topic string, handler func(DeliverMessage)) (Subscription, error)
	// RequestReply send a request and wait for a response
	RequestReply(ctx context.Context, topic string, payload string) (*Message, error)
	// QueueSubscribe create a queue subscription
	QueueSubscribe(ctx context.Context, topic string, handler func(DeliverMessage)) (Subscription, error)
	// Close closes all MessageQueue connections
	Close() error
}

// Subscription indicates a subscription
type Subscription interface {
	// Unsubscribe from the subscription
	Unsubscribe() error
	// Topic returns the subscribed topic
	Topic() string
	// IsQueue returns whether the subscription is a queue subscription
	IsQueue() bool
	// Queue return the queue name (if it is a queue subscription)
	Queue() (string, bool)
	// IsActive returns whether the subscription is active
	IsActive() bool
	// Pause suspend receiving message
	Pause() error
	// Resume resume receiving message
	Resume() error
	// HandlerMessage set up a new message handler
	HandlerMessage(handler func(*Message)) error
	// Pending returns the number of messages to be processed
	Pending() (int, error)
	// Delivered returns the number of messages that have been delivered
	Delivered() (uint64, error)
	// Dropped returns the number of messages discarded due to client timeouts or disconnections
	Dropped() (uint64, error)
	// ClearMaxPending clear statistics about the maximum number of messages to be processed
	ClearMaxPending() error
	// MaxPending returns the maximum number of pending messages
	MaxPending() (int, int, error)
}
