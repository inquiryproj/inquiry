package runs

import "github.com/google/uuid"

// Producer produces messages to the queue.
type Producer interface {
	// Produce produces a message to the queue.
	Produce(projectID uuid.UUID) error
}
