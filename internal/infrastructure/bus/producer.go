package bus

import (
	"fmt"
	"sync"
	"time"
)

type MessageBus[T any] struct {
	subscribers   map[string]chan T
	mu            sync.RWMutex
	closeNotifier chan struct{}
}

func NewMessageBus[T any]() *MessageBus[T] {
	p := &MessageBus[T]{
		subscribers:   make(map[string]chan T),
		mu:            sync.RWMutex{},
		closeNotifier: make(chan struct{}),
	}

	return p
}

func (bus *MessageBus[T]) Subscribe() (<-chan T, string, error) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	subscriberID := fmt.Sprintf("%d", time.Now().UnixNano())

	ch := make(chan T)
	bus.subscribers[subscriberID] = ch

	return ch, subscriberID, nil
}

func (bus *MessageBus[T]) Unsubscribe(subscriberID string) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	ch, exists := bus.subscribers[subscriberID]
	if !exists {
		return fmt.Errorf("subscriber with ID %s does not exist", subscriberID)
	}

	close(ch)
	delete(bus.subscribers, subscriberID)
	return nil
}

func (bus *MessageBus[T]) Publish(message T) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	for _, ch := range bus.subscribers {
		select {
		case ch <- message:
		default:
		}
	}
}

func (bus *MessageBus[T]) Close() {
	// Signal that we're closing
	close(bus.closeNotifier)

	// Close all subscriber channels
	bus.mu.Lock()
	defer bus.mu.Unlock()

	for _, ch := range bus.subscribers {
		close(ch)
	}

	bus.subscribers = make(map[string]chan T)
}
