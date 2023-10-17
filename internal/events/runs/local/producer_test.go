package local

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	mockProcessor := &mockProcessor{}

	stream := make(chan uuid.UUID)

	p := NewProducer(stream)

	c := NewConsumer(stream, mockProcessor)
	go func() {
		assert.NoError(t, c.Consume())
	}()
	assert.NoError(t, p.Produce(uuid.New()))
	assert.NoError(t, c.Shutdown(context.Background()))
	assert.Equal(t, 1, mockProcessor.counter)
}
