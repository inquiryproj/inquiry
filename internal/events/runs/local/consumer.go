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
	"github.com/inquiryproj/inquiry/internal/events/runs"
)

// ErrCloseTimeout is returned when the consumer close times out.
var ErrCloseTimeout = fmt.Errorf("consumer close timed out")

// Consumer is the local consumer implementation.
type Consumer struct {
	stream    chan uuid.UUID
	closeChan chan (struct{})
	doneChan  chan (struct{})

	workerPool workers.Pool[uuid.UUID, *app.ProjectRunOutput]
}

// NewConsumer creates a new local consumer.
// the local consumer has the limitation that it cannot guranatee that all runs are processed.
// In case of a shutdown, the consumer will try to process all runs, within the given timeout.
// If the timeout is reached, the consumer will stop processing runs and set all non finished
// runs to canceled.
func NewConsumer(stream chan (uuid.UUID), processor runs.Processor) *Consumer {
	c := &Consumer{
		stream:    stream,
		closeChan: make(chan (struct{})),
		doneChan:  make(chan (struct{})),
	}
	workerPool := workers.NewUnBufferedPool(context.Background(),
		processor.Process,
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

// Close closes the consumer.
func (c *Consumer) Close() error {
	c.closeChan <- struct{}{}
	select {
	case <-c.doneChan:
		return nil
	case <-time.After(time.Second * 10):
		// FIXME: If timeout set all running and created runs to canceled.
		// We might even want to do this on startup
		// Local consumers are not able to guarantee that all runs are processed.
		return ErrCloseTimeout
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
