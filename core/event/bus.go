package event

import (
	"context"
	"sync"
)

type Handler func(context.Context, Event)

type Bus struct {
	mu          sync.RWMutex
	subscribers map[Type]map[uint64]Handler
	nextID      uint64
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[Type]map[uint64]Handler),
	}
}

func (b *Bus) Subscribe(eventType Type, handler Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	id := b.nextID

	if _, exists := b.subscribers[eventType]; !exists {
		b.subscribers[eventType] = make(map[uint64]Handler)
	}

	b.subscribers[eventType][id] = handler

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		handlers, exists := b.subscribers[eventType]
		if !exists {
			return
		}

		delete(handlers, id)

		if len(handlers) == 0 {
			delete(b.subscribers, eventType)
		}
	}
}

func (b *Bus) Publish(ctx context.Context, evt Event) {
	b.mu.RLock()

	handlersMap := b.subscribers[evt.EventType]

	handlers := make([]Handler, 0, len(handlersMap))

	for _, handler := range handlersMap {
		handlers = append(handlers, handler)
	}

	b.mu.RUnlock()

	for _, handler := range handlers {
		handler(ctx, evt)
	}
}
