package runs

import (
	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/events/runs/local"
	"github.com/inquiryproj/inquiry/internal/events/runs/run"
)

// ConsumerType represents the consumer type.
type ConsumerType string

// Different consumer types.
const (
	ConsumerTypeLocal ConsumerType = "local"
)

// Options represents the options.
type Options struct {
	ConsumerType ConsumerType
}

func defaultOptions() *Options {
	return &Options{
		ConsumerType: ConsumerTypeLocal,
	}
}

// Opts represents a function that modifies the options.
type Opts func(*Options)

// NewProducerConsumer creates a new producer and consumer.
func NewProducerConsumer(runProcessor run.Processor, opts ...Opts) (Producer, Consumer, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	switch options.ConsumerType {
	case ConsumerTypeLocal:
		stream := make(chan uuid.UUID)
		return local.NewProducer(stream), local.NewConsumer(stream, runProcessor), nil
	default:
		return nil, nil, ErrUnknownConsumerType
	}
}
