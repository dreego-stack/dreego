package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestEventBusPublishDeliversToSubscribedHandler(t *testing.T) {
	bus := NewInMemoryBus[string]()
	got := make(chan string, 1)
	_, err := bus.Subscribe(context.Background(), func(e string) {
		got <- e
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.Publish(context.Background(), "hello"); err != nil {
		t.Fatalf("publish: %v", err)
	}

	select {
	case v := <-got:
		if v != "hello" {
			t.Fatalf("received %q, want %q", v, "hello")
		}
	default:
		t.Fatal("subscribed handler did not receive published event")
	}
}

func TestEventBusMultipleSubscribersAllReceive(t *testing.T) {
	bus := NewInMemoryBus[int]()
	const n = 3
	var mu sync.Mutex
	received := make([]int, 0, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		_, err := bus.Subscribe(context.Background(), func(e int) {
			defer wg.Done()
			mu.Lock()
			received = append(received, e)
			mu.Unlock()
		})
		if err != nil {
			t.Fatalf("subscribe %d: %v", i, err)
		}
	}

	if err := bus.Publish(context.Background(), 42); err != nil {
		t.Fatalf("publish: %v", err)
	}
	wg.Wait()

	if len(received) != n {
		t.Fatalf("received %d deliveries, want %d", len(received), n)
	}
	for _, v := range received {
		if v != 42 {
			t.Fatalf("delivered %d, want %d", v, 42)
		}
	}
}

func TestEventBusUnsubscribeStopsDelivery(t *testing.T) {
	bus := NewInMemoryBus[string]()
	got := make(chan string, 4)
	sub, err := bus.Subscribe(context.Background(), func(e string) {
		got <- e
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	bus.Publish(context.Background(), "first")
	bus.Unsubscribe(sub)
	bus.Publish(context.Background(), "second")

	if len(got) != 1 {
		t.Fatalf("handler received %d events after unsubscribe, want 1", len(got))
	}
	if v := <-got; v != "first" {
		t.Fatalf("received %q, want %q", v, "first")
	}
}

func TestEventBusTypedBusesAreIndependent(t *testing.T) {
	intBus := NewInMemoryBus[int]()
	structBus := NewInMemoryBus[struct{ Name string }]()

	intGot := make(chan int, 1)
	structGot := make(chan struct{ Name string }, 1)
	intBus.Subscribe(context.Background(), func(e int) { intGot <- e })
	structBus.Subscribe(context.Background(), func(e struct{ Name string }) { structGot <- e })

	if err := intBus.Publish(context.Background(), 7); err != nil {
		t.Fatalf("int publish: %v", err)
	}
	if err := structBus.Publish(context.Background(), struct{ Name string }{Name: "dreego"}); err != nil {
		t.Fatalf("struct publish: %v", err)
	}

	if v := <-intGot; v != 7 {
		t.Fatalf("int bus delivered %d, want 7", v)
	}
	if v := <-structGot; v.Name != "dreego" {
		t.Fatalf("struct bus delivered %+v, want Name=dreego", v)
	}
}

func TestEventBusPublishWithCanceledContextReturnsError(t *testing.T) {
	bus := NewInMemoryBus[string]()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := bus.Publish(ctx, "x"); err == nil {
		t.Fatal("publish with canceled context returned nil error, want error")
	} else if !errors.Is(err, context.Canceled) {
		t.Fatalf("publish error = %v, want context.Canceled", err)
	}
}

func TestEventBusHandlerPanicIsRecoveredAndReturned(t *testing.T) {
	bus := NewInMemoryBus[string]()
	bus.Subscribe(context.Background(), func(e string) {
		panic(fmt.Sprintf("boom: %s", e))
	})

	err := bus.Publish(context.Background(), "event")
	if err == nil {
		t.Fatal("publish returned nil error, want recovered panic")
	}
	if got := err.Error(); got != "boom: event" {
		t.Fatalf("publish error = %q, want %q", got, "boom: event")
	}
}

func TestEventBusSubscriptionIDIsOpaqueAndUnique(t *testing.T) {
	bus := NewInMemoryBus[int]()
	subs := make(map[uint64]bool)
	for i := 0; i < 5; i++ {
		sub, err := bus.Subscribe(context.Background(), func(int) {})
		if err != nil {
			t.Fatalf("subscribe %d: %v", i, err)
		}
		if subs[sub.ID()] {
			t.Fatalf("subscription ID %d returned twice", sub.ID())
		}
		subs[sub.ID()] = true
	}
}

// TestEventBusConcurrentPublishWithStableSubscribers: all handlers subscribe
// before any publish starts, so every handler must receive exactly n events.
func TestEventBusConcurrentPublishWithStableSubscribers(t *testing.T) {
	bus := NewInMemoryBus[int]()
	const (
		publishers = 8
		handlers   = 4
		perHandler = publishers
	)

	var received [handlers]atomic.Int64
	for i := 0; i < handlers; i++ {
		idx := i
		if _, err := bus.Subscribe(context.Background(), func(int) {
			received[idx].Add(1)
		}); err != nil {
			t.Fatalf("subscribe %d: %v", idx, err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(publishers)
	for p := 0; p < publishers; p++ {
		go func() {
			defer wg.Done()
			if err := bus.Publish(context.Background(), p); err != nil {
				t.Errorf("concurrent publish: %v", err)
			}
		}()
	}
	wg.Wait()

	for i := 0; i < handlers; i++ {
		if got := received[i].Load(); got != perHandler {
			t.Fatalf("handler %d received %d events, want %d", i, got, perHandler)
		}
	}
}

// TestEventBusConcurrentSubscribeUnsubscribeDoesNotRace: concurrent
// subscribe/unsubscribe must not panic or deadlock; the publish below just
// needs to complete.
func TestEventBusConcurrentSubscribeUnsubscribeDoesNotRace(t *testing.T) {
	bus := NewInMemoryBus[int]()
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				sub, err := bus.Subscribe(context.Background(), func(int) {})
				if err != nil {
					t.Errorf("concurrent subscribe: %v", err)
				}
				bus.Unsubscribe(sub)
			}
		}()
	}
	wg.Wait()

	if err := bus.Publish(context.Background(), 1); err != nil {
		t.Fatalf("publish after concurrent churn: %v", err)
	}
}

// TestEventBusUnsubscribeDuringPublishStillDeliversCurrentEvent: Publish
// snapshots handlers under RLock before delivery, so handler A unsubscribing
// handler B mid-publish must not prevent B from receiving the current event.
func TestEventBusUnsubscribeDuringPublishStillDeliversCurrentEvent(t *testing.T) {
	bus := NewInMemoryBus[string]()
	bGot := make(chan string, 1)

	var bSub Subscription
	bus.Subscribe(context.Background(), func(e string) {
		if e == "first" {
			bus.Unsubscribe(bSub)
		}
	})
	bSub, err := bus.Subscribe(context.Background(), func(e string) {
		bGot <- e
	})
	if err != nil {
		t.Fatalf("subscribe B: %v", err)
	}

	bus.Publish(context.Background(), "first")
	select {
	case v := <-bGot:
		if v != "first" {
			t.Fatalf("B received %q, want %q", v, "first")
		}
	default:
		t.Fatal("B did not receive the in-flight event despite being unsubscribed mid-publish")
	}

	bus.Publish(context.Background(), "second")
	if len(bGot) != 0 {
		t.Fatal("B received an event after being unsubscribed mid-publish")
	}
}

// TestEventBusSelfUnsubscribeDuringDelivery: a handler that unsubscribes
// itself while handling an event must not receive any subsequent events.
func TestEventBusSelfUnsubscribeDuringDelivery(t *testing.T) {
	bus := NewInMemoryBus[string]()
	got := make(chan string, 2)

	var sub Subscription
	var err error
	sub, err = bus.Subscribe(context.Background(), func(e string) {
		bus.Unsubscribe(sub)
		got <- e
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	bus.Publish(context.Background(), "first")
	if v := <-got; v != "first" {
		t.Fatalf("received %q, want %q", v, "first")
	}

	bus.Publish(context.Background(), "second")
	if len(got) != 0 {
		t.Fatal("self-unsubscribed handler received a subsequent event")
	}
}
