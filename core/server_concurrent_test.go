package core

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestListenRejectsConcurrentServer(t *testing.T) {
	app := New()
	first := make(chan error, 1)
	go func() { first <- app.Listen("127.0.0.1:0") }()

	deadline := time.Now().Add(time.Second)
	for {
		app.mu.RLock()
		running := app.server != nil
		app.mu.RUnlock()
		if running {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("first server did not start")
		}
		time.Sleep(time.Millisecond)
	}

	if err := app.Listen("127.0.0.1:0"); !errors.Is(err, ErrServerRunning) {
		t.Fatalf("second Listen error = %v, want ErrServerRunning", err)
	}
	if err := app.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := <-first; err != nil {
		t.Fatal(err)
	}
}
