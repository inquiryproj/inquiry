// Package runs is the event stream for runs.
package runs

import "context"

// Consumer consumes messages from the queue.
type Consumer interface {
	// Consume consumes a message from the queue.
	Consume() error
	// Shutdown shuts down the consumer.
	Shutdown(ctx context.Context) error
}
