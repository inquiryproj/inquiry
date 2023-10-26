package events

// Producer produces messages to the queue.
type Producer[T any] interface {
	// Produce produces a message to the queue.
	Produce(runID T) error
}
