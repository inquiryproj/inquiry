// Package events is the definition for event consumer & producer.
package events

import "context"

// Consumer consumes messages from the queue.
type Consumer interface {
	// Consume consumes a message from the queue.
	Consume() error
	// Shutdown shuts down the consumer.
	Shutdown(ctx context.Context) error
}
