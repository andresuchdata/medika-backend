package events

import (
	"context"
)

// Event represents a domain event
type Event interface {
	EventType() string
	EventData() map[string]interface{}
}

// Bus represents an event bus
type Bus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(ctx context.Context, eventType string, handler Handler) error
}

// Handler represents an event handler
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Simple in-memory event bus implementation
type InMemoryBus struct {
	handlers map[string][]Handler
}

func NewInMemoryBus() Bus {
	return &InMemoryBus{
		handlers: make(map[string][]Handler),
	}
}

func (b *InMemoryBus) Publish(ctx context.Context, event Event) error {
	handlers, exists := b.handlers[event.EventType()]
	if !exists {
		return nil // No handlers for this event type
	}

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			// Log error but don't fail the entire event publishing
			// In a real implementation, you might want to retry or use a dead letter queue
			continue
		}
	}

	return nil
}

func (b *InMemoryBus) Subscribe(ctx context.Context, eventType string, handler Handler) error {
	if _, exists := b.handlers[eventType]; !exists {
		b.handlers[eventType] = make([]Handler, 0)
	}

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}
