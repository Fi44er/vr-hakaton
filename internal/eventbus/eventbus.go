package eventbus

import (
	"sync"
)

// EventBus is a simple in-memory event bus.
type EventBus struct {
	subscribers map[string][]chan interface{}
	mu          sync.RWMutex
}

// New creates a new instance of EventBus.
func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Subscribe allows a subscriber to listen for events of a specific type.
func (eb *EventBus) Subscribe(eventType string, ch chan interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
}

// Publish sends an event to all subscribers of the given event type.
func (eb *EventBus) Publish(eventType string, event interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if chans, found := eb.subscribers[eventType]; found {
		for _, ch := range chans {
			ch <- event
		}
	}
}

