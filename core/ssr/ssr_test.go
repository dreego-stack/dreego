package ssr

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestListenDelegates(t *testing.T) {
	app := dreego.New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	errCh := make(chan error, 1)
	go func() { errCh <- Listen(app, fmt.Sprintf("127.0.0.1:%d", port)) }()

	time.Sleep(100 * time.Millisecond)

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Shutdown(app, shutCtx); err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Listen returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Listen did not return after Shutdown")
	}
}

func TestDefaultAddr(t *testing.T) {
	if got := DefaultAddr(); got != ":8080" {
		t.Fatalf("DefaultAddr() = %q, want %q", got, ":8080")
	}
}
