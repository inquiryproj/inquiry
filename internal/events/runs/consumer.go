// Package runs is the event stream for runs.
package runs

// Consumer consumes messages from the queue.
type Consumer interface {
	// Consume consumes a message from the queue.
	Consume() error
	// Close closes the consumer.
	Close() error
}
