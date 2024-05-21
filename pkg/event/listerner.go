package event

import (
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
	listeners = []Listener{}
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

		var wg sync.WaitGroup
		wg.Add(len(listeners))
		for _, listener := range listeners {
			if listener.Name == event.Name {
				go func(l Listener, e Event) {
					defer wg.Done()
					var err error
					for i := 0; i < RETRY_COUNT; i++ {
						err = l.Handler(e.Payload...)
						if err == nil {
							break
						}
						log.Printf("Error processing event %s (listener %s, attempt %d): %v", event.Name, l.Name, i+1, err)
						time.Sleep(RETRY_DELAY)
					}
					if err != nil {
						log.Printf("Failed to process event %s for listener %s after %d retries: %v", event.Name, l.Name, RETRY_COUNT, err)
					}
				}(listener, event)
			}
		}
		wg.Wait()
	}
}

// RegisterListener registers a listener function for a specific event name
func RegisterListener(listener Listener) error {
	mutex.Lock()
	defer mutex.Unlock()

	listeners = append(listeners, listener)
	return nil
}

// UnregisterListener removes a listener function based on its name
func UnregisterListener(name string) error {
	mutex.Lock()
	defer mutex.Unlock()

	newList := make([]Listener, 0, len(listeners))
	for _, listener := range listeners {
		if listener.Name != name {
			newList = append(newList, listener)
		}
	}
	listeners = newList
	return nil
}
