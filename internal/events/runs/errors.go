package runs

import "fmt"

// ErrUnknownConsumerType is returned when the consumer type is unknown.
var ErrUnknownConsumerType = fmt.Errorf("unknown consumer type")
