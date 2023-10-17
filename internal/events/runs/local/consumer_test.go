package local

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/inquiryproj/inquiry/internal/app"
)

type mockProcessor struct {
	sync.Mutex
	counter         int
	processDuration time.Duration
}

func (m *mockProcessor) Process(_ uuid.UUID) (*app.ProjectRunOutput, error) {
	time.Sleep(m.processDuration)
	m.Lock()
	defer m.Unlock()
	m.counter++
	return nil, nil
}

func TestConsumerCloseGraceful(t *testing.T) {
	mockProcessor := &mockProcessor{}

	stream := make(chan uuid.UUID)
	c := NewConsumer(stream, mockProcessor)

	go func() {
		assert.NoError(t, c.Consume())
	}()

	for i := 0; i < 10000; i++ {
		stream <- uuid.New()
	}
	err := c.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 10000, mockProcessor.counter)
}

func TestConsumerCloseTimeout(t *testing.T) {
	mockProcessor := &mockProcessor{
		processDuration: 10 * time.Second,
	}

	stream := make(chan uuid.UUID)
	c := NewConsumer(stream, mockProcessor, WithCloseTimeout(1*time.Millisecond))

	go func() {
		assert.NoError(t, c.Consume())
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			stream <- uuid.New()
		}
	}()
	time.Sleep(time.Millisecond * 100)
	err := c.Shutdown(context.Background())
	assert.ErrorIs(t, err, ErrCloseTimeout)
}
