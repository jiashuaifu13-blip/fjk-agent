package event

import (
	"context"
	"sync"
)

type Handler func(context.Context, Event)

type Bus struct {
	mu       sync.RWMutex
	handlers map[Type][]Handler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[Type][]Handler),
	}
}

func (b *Bus) Subscribe(eventType Type, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(
		b.handlers[eventType],
		handler,
	)
}

func (b *Bus) Publish(ctx context.Context, evt Event) {
	b.mu.RLock()

	handlers := append(
		[]Handler(nil),
		b.handlers[evt.EventType]...,
	)

	b.mu.RUnlock()

	for _, handler := range handlers {
		handler(ctx, evt)
	}
}
