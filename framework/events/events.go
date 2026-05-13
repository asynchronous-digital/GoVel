// Events system (Phase 3)
// Provides event dispatching and listener management

package events

import (
	"context"
)

// Event is the base interface for all events
type Event interface {
	EventName() string
}

// Listener handles an event
type Listener func(ctx context.Context, event Event) error

// Dispatcher manages events and listeners
type Dispatcher interface {
	// Register a listener for an event
	Listen(eventName string, listener Listener)

	// Dispatch an event
	Dispatch(ctx context.Context, event Event) error

	// Dispatch multiple events
	DispatchBatch(ctx context.Context, events ...Event) error

	// Forget listeners for an event
	Forget(eventName string)

	// Forget all listeners
	Flush()
}

// Provider manages event services
type Provider interface {
	Boot() error
	Register(dispatcher Dispatcher) error
}

// Event examples:
/*
type UserCreated struct {
    User *models.User
    Timestamp time.Time
}

func (e UserCreated) EventName() string {
    return "user.created"
}

type OrderPlaced struct {
    Order *models.Order
    Timestamp time.Time
}

func (e OrderPlaced) EventName() string {
    return "order.placed"
}
*/

// Usage:
/*
// Register listener
dispatcher.Listen("user.created", func(ctx context.Context, event Event) error {
    user := event.(UserCreated).User
    // Send welcome email, etc.
    return nil
})

// Dispatch event
dispatcher.Dispatch(ctx, UserCreated{User: user, Timestamp: time.Now()})
*/

// Queued events - events that should be processed asynchronously
type QueuedEvent interface {
	Event
	Queue() string
}

// Future implementation with async support
