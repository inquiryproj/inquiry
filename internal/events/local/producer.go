package local

import (
	"context"

	"github.com/google/uuid"
)

// Producer is the local producer implementation.
type Producer struct {
	stream chan uuid.UUID
}

// NewProducer creates a new local producer.
func NewProducer(stream chan uuid.UUID) *Producer {
	return &Producer{
		stream: stream,
	}
}

// Produce produces a new project run.
func (p *Producer) Produce(_ context.Context, projectID uuid.UUID) error {
	p.stream <- projectID
	return nil
}
