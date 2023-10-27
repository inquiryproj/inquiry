package events

import (
	"context"
)

// Producer produces messages to the queue.
type Producer[T any] interface {
	// Produce produces a message to the queue.
	Produce(ctx context.Context, runID T) error
}
