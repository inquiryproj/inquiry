// Package completions initialises the completions event stream consumer and producer.
package completions

import (
	"github.com/google/uuid"

	"github.com/inquiryproj/inquiry/internal/events"
	"github.com/inquiryproj/inquiry/internal/events/local"
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
func NewProducerConsumer(completionsProcessor Processor, opts ...Opts) (events.Producer, events.Consumer, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	switch options.ConsumerType {
	case ConsumerTypeLocal:
		stream := make(chan uuid.UUID)
		return local.NewProducer(stream), local.NewConsumer(stream, completionsProcessor.Process), nil
	default:
		return nil, nil, events.ErrUnknownConsumerType
	}
}
