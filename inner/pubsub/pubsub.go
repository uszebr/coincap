package pubsub

import (
	"log/slog"
	"sync"
)

// Refactored to have  Pub/Sub build into the DataStore
//
// Publisher/Subscriber
type PubSub[T any] struct {
	subscribers map[chan T]struct{}
	mu          sync.RWMutex
}

// creates a new PubSub system for type T
func NewPubSub[T any]() *PubSub[T] {
	slog.Debug("NewPubSub")
	return &PubSub[T]{
		subscribers: make(map[chan T]struct{}),
	}
}

// adds a subscriber to the PubSub and returns a channel to watch updates
func (ps *PubSub[T]) Subscribe() chan T {
	slog.Debug("Subscribing new channel")
	ch := make(chan T)
	ps.mu.Lock()
	ps.subscribers[ch] = struct{}{}
	ps.mu.Unlock()
	return ch
}

// removes a subscriber from the PubSub; channel is closed
func (ps *PubSub[T]) Unsubscribe(ch chan T) {
	slog.Debug("Unsubscribing channel")
	ps.mu.Lock()
	delete(ps.subscribers, ch)
	close(ch)
	ps.mu.Unlock()
}

// Publishing data to all subscribers
func (ps *PubSub[T]) Publish(data T) {
	slog.Debug("Publishing data to subscribers")
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for ch := range ps.subscribers {
		ch <- data
	}
}
