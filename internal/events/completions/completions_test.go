package completions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProducerConsumerLocal(t *testing.T) {
	producer, consumer, err := NewProducerConsumer(&processor{})
	assert.NoError(t, err)
	assert.NotNil(t, producer)
	assert.NotNil(t, consumer)
}

func TestNewProducerConsumerUnknown(t *testing.T) {
	producer, consumer, err := NewProducerConsumer(&processor{},
		WithConsumerType(ConsumerType("unknown")),
	)
	assert.Error(t, err)
	assert.Nil(t, producer)
	assert.Nil(t, consumer)
}
