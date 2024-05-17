package event

import (
	"fmt"
	"time"
)

const (
	EVENT_BACK_OFF_DELAY = 1
	EVENT_TIMEOUT        = 10
)

// Emit adds an event to the event channel
// Ensures the event is eventually delivered, keeps retrying
func Emit(event Event) error {
	for {
		select {
		case events <- event:
			return nil
		case <-time.After(time.Second * EVENT_TIMEOUT):
			return fmt.Errorf("event queue full after timeout, event dropped: %v", event)
		default:
			// Channel full, wait and retry
			time.Sleep(time.Second * EVENT_BACK_OFF_DELAY)
		}
	}
}

// LazyEmitLazy adds an event to the event channel
// It throws an error when the channel is full, does not retry
func LazyEmit(event Event) error {
	select {
	case events <- event:
		return nil
	default:
		// Channel full, handle buffer overflow (optional)
		return fmt.Errorf("event queue full, event dropped: %v", event)
	}
}
