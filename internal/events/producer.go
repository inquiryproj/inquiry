package events

import (
	"context"

	"github.com/google/uuid"
)

// Producer produces messages to the queue.
type Producer interface {
	// Produce produces a message to the queue.
	Produce(ctx context.Context, runID uuid.UUID) error
}
