package local

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/inquiryproj/inquiry/internal/app"
)

type mockProcessor struct {
	sync.Mutex
	counter int
}

func (m *mockProcessor) Process(_ uuid.UUID) (*app.ProjectRunOutput, error) {
	m.Lock()
	defer m.Unlock()
	m.counter++

	return nil, nil
}

func TestConsumer(t *testing.T) {
	mockProcessor := &mockProcessor{}

	stream := make(chan uuid.UUID)
	c := NewConsumer(stream, mockProcessor)

	go func() {
		assert.NoError(t, c.Consume())
	}()

	for i := 0; i < 10000; i++ {
		stream <- uuid.New()
	}
	err := c.Close()
	assert.NoError(t, err)
	assert.Equal(t, 10000, mockProcessor.counter)
}
