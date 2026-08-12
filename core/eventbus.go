package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// EventBus is a typed pub/sub contract. Implementations may back it with
// in-memory storage, Redis, NATS or similar; the interface abstracts the
// transport so core code stays transport-agnostic.
type EventBus[T any] interface {
	// Publish delivers event to every subscribed handler. Implementations
	// must respect ctx cancellation.
	Publish(ctx context.Context, event T) error
	// Subscribe registers handler for all future events. Implementations
	// may ignore ctx. The returned Subscription is an opaque handle for
	// Unsubscribe.
	Subscribe(ctx context.Context, handler func(T)) (Subscription, error)
	// Unsubscribe removes a previously registered subscription.
	Unsubscribe(sub Subscription)
}

// Subscription is an opaque handle identifying a registered handler.
type Subscription interface {
	ID() uint64
}

type subscription struct {
	id uint64
}

func (s subscription) ID() uint64 { return s.id }

// NewInMemoryBus returns an EventBus backed by process-local memory. It is
// safe for concurrent use. Handlers run synchronously inside Publish; a
// panicking handler is recovered and the panic value is returned as error.
func NewInMemoryBus[T any]() EventBus[T] {
	return &inMemoryBus[T]{
		subs: make(map[uint64]func(T)),
	}
}

type inMemoryBus[T any] struct {
	mu   sync.RWMutex
	subs map[uint64]func(T)
	next uint64
}

func (b *inMemoryBus[T]) Publish(ctx context.Context, event T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.mu.RLock()
	handlers := make([]func(T), 0, len(b.subs))
	for _, h := range b.subs {
		handlers = append(handlers, h)
	}
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := b.deliver(ctx, h, event); err != nil {
			return err
		}
	}
	return nil
}

func (b *inMemoryBus[T]) deliver(ctx context.Context, h func(T), event T) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = v
			default:
				err = &panicError{value: r}
			}
		}
	}()
	h(event)
	return nil
}

func (b *inMemoryBus[T]) Subscribe(_ context.Context, handler func(T)) (Subscription, error) {
	if handler == nil {
		return nil, errors.New("eventbus: nil handler")
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.next++
	sub := subscription{id: b.next}
	b.subs[sub.id] = handler
	return sub, nil
}

func (b *inMemoryBus[T]) Unsubscribe(sub Subscription) {
	if sub == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subs, sub.ID())
}

// panicError wraps a non-error panic value so it satisfies error.
type panicError struct {
	value any
}

func (e *panicError) Error() string {
	return fmt.Sprintf("%v", e.value)
}
