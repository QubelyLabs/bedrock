package event

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	BUFFER_LIMIT = 1000
	RETRY_COUNT  = 3               // Number of retries for failed listener calls
	RETRY_DELAY  = 5 * time.Second // Delay between retries
)

var (
	events    = make(chan Event, BUFFER_LIMIT)
	listeners = map[string]Listener{}
	mutex     = &sync.Mutex{}
)

// Start start a for ever loop that listen to Event via the events channel
// It uses the Event.Name property to determine the handler to call
func StartListener() error {
	for {
		// Receive an event from the channel
		event, ok := <-events
		if !ok {
			return nil
		}

		// Check if a listener exists for the event name
		listener, ok := listeners[event.Name]
		if !ok {
			fmt.Printf("No listener found for event: %s\n", event.Name)
			continue
		}

		var err error
		for i := 0; i < RETRY_COUNT; i++ {
			err = listener(event.Payload...)
			if err == nil {
				break
			}
			log.Printf("Error processing event %s (attempt %d): %v", event.Name, i+1, err)
			time.Sleep(RETRY_DELAY)
		}

		if err != nil {
			log.Printf("Failed to process event %s after %d retries: %v", event.Name, RETRY_COUNT, err)
		}
	}
}

// RegisterListener registers a listener function for a specific event name
func RegisterListener(name string, listener Listener) error {
	mutex.Lock()
	defer mutex.Unlock()

	listeners[name] = listener
	return nil
}

// UnregisterListener removes a listener function for a specific event name
func UnregisterListener(name string) error {
	mutex.Lock()
	defer mutex.Unlock()

	delete(listeners, name)
	return nil
}
