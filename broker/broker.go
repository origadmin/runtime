// Package broker defines a generic interface for a message broker.
package broker

import (
	"context"
)

// Message is a generic message structure.
type Message struct {
	// ID is a unique identifier for the message.
	ID string
	// Metadata contains message metadata.
	Metadata map[string]string
	// Payload is the message content.
	Payload []byte
}

// Publisher defines the interface for publishing messages.
type Publisher interface {
	// Publish sends messages to a topic.
	Publish(ctx context.Context, topic string, messages ...*Message) error
	// Close terminates the publisher's connection.
	Close() error
}

// Subscriber defines the interface for subscribing to messages.
type Subscriber interface {
	// Subscribe returns a channel of messages for a given topic.
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
	// Close terminates the subscriber's connection.
	Close() error
}

// Broker combines the Publisher and Subscriber interfaces.
type Broker interface {
	Publisher
	Subscriber
}
