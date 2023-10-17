// Package local implements the local event stream for runs.
package local

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wimspaargaren/workers"

	"github.com/inquiryproj/inquiry/internal/app"
	"github.com/inquiryproj/inquiry/internal/events/runs/run"
)

// Options represents the options for the consumer.
type Options struct {
	// CloseTimeout is the timeout for closing the consumer.
	CloseTimeout time.Duration
	// ParallelProcessors is the number of parallel processors.
	ParallelProcessors int
}

// Opts represents a function that modifies the options.
type Opts func(*Options)

// WithCloseTimeout sets the close timeout.
func WithCloseTimeout(timeout time.Duration) Opts {
	return func(o *Options) {
		o.CloseTimeout = timeout
	}
}

func defaultOptions() *Options {
	return &Options{
		CloseTimeout:       time.Second * 10,
		ParallelProcessors: 25,
	}
}

// ErrCloseTimeout is returned when the consumer close times out.
var ErrCloseTimeout = fmt.Errorf("consumer close timed out")

// Consumer is the local consumer implementation.
type Consumer struct {
	stream    chan uuid.UUID
	closeChan chan (struct{})
	doneChan  chan (struct{})

	closeTimeout time.Duration

	workerPool workers.Pool[uuid.UUID, *app.ProjectRunOutput]
}

// NewConsumer creates a new local consumer.
// the local consumer has the limitation that it cannot guranatee that all runs are processed.
// In case of a shutdown, the consumer will try to process all runs, within the given timeout.
// If the timeout is reached, the consumer will stop processing runs and set all non finished
// runs to canceled.
func NewConsumer(stream chan (uuid.UUID), runProcessor run.Processor, opts ...Opts) *Consumer {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	c := &Consumer{
		stream:       stream,
		closeChan:    make(chan (struct{})),
		doneChan:     make(chan (struct{})),
		closeTimeout: options.CloseTimeout,
	}
	workerPool := workers.NewUnBufferedPool(context.Background(),
		runProcessor.Process,
		workers.WithWorkers(options.ParallelProcessors),
	)
	c.workerPool = workerPool
	return c
}

// Consume consumes the stream.
func (c *Consumer) Consume() error {
	go c.processResults()
	for {
		select {
		case id := <-c.stream:
			err := c.workerPool.AddJob(id)
			if err != nil {
				return err
			}
		case <-c.closeChan:
			c.workerPool.Done()
			return nil
		}
	}
}

func (c *Consumer) processResults() {
	x, y := c.workerPool.ResultChannels()
	for {
		select {
		case _, ok := <-x:
			if !ok {
				c.doneChan <- struct{}{}
			}
		case err, ok := <-y:
			if !ok {
				c.doneChan <- struct{}{}
			}
			// FIXME ADD LOGGER, error should already be handled in processor
			log.Default().Println(err)
		}
	}
}

// Shutdown shuts the consumer down.
func (c *Consumer) Shutdown(ctx context.Context) error {
	go func() {
		c.closeChan <- struct{}{}
	}()
	select {
	case <-c.doneChan:
		return nil
	case <-time.After(c.closeTimeout):
		// FIXME: If timeout set all running and created runs to canceled.
		// We might even want to do this on startup
		// Local consumers are not able to guarantee that all runs are processed.
		return ErrCloseTimeout
	case <-ctx.Done():
		return ctx.Err()
	}
}

// FIXME: Fix the following:
// endpoint:
//   create run
//
// Happy flow
// consumer {
//   update db started run
//   update db finished run
// }
//
// Shutdown flow
// consumer {
//   update db started run
//   on shutdown: canceled run
// }
//
// Timeout flow
// consumer {
//   update db started run
//   on timeout: timedout run
// }
