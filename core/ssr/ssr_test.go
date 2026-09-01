package ssr

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestHostServesAndShutsDown(t *testing.T) {
	app := dreego.New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("GET", "/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}); err != nil {
		t.Fatal(err)
	}
	listener := listenLocal(t)
	host := New(app)
	done := make(chan error, 1)
	go func() { done <- host.Serve(listener) }()
	waitRunning(t, host)

	response, err := http.Get("http://" + listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusNoContent)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := host.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if err := host.Shutdown(context.Background()); err != nil {
		t.Fatalf("second Shutdown = %v, want nil", err)
	}
}

func TestHostRejectsConcurrentServe(t *testing.T) {
	host := New(dreego.New())
	firstListener := listenLocal(t)
	done := make(chan error, 1)
	go func() { done <- host.Serve(firstListener) }()
	waitRunning(t, host)

	secondListener := listenLocal(t)
	if err := host.Serve(secondListener); !errors.Is(err, ErrServerRunning) {
		t.Fatalf("second Serve = %v, want ErrServerRunning", err)
	}
	if err := host.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestHostUsesSecureConfigDefaults(t *testing.T) {
	host := New(dreego.New(), ServerConfig{})
	listener := listenLocal(t)
	done := make(chan error, 1)
	go func() { done <- host.Serve(listener) }()
	waitRunning(t, host)

	host.mu.Lock()
	server := host.server
	host.mu.Unlock()
	if server.ReadHeaderTimeout <= 0 || server.ReadTimeout <= 0 || server.WriteTimeout <= 0 || server.IdleTimeout <= 0 {
		t.Fatal("server timeouts must all have secure defaults")
	}
	if server.MaxHeaderBytes <= 0 {
		t.Fatal("MaxHeaderBytes must have a secure default")
	}
	if err := host.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestListenReturnsBindError(t *testing.T) {
	listener := listenLocal(t)
	defer listener.Close()
	if err := New(dreego.New()).Listen(listener.Addr().String()); err == nil {
		t.Fatal("expected occupied address to fail")
	}
}

func TestDefaultAddr(t *testing.T) {
	if got := DefaultAddr(); got != ":8080" {
		t.Fatalf("DefaultAddr() = %q, want :8080", got)
	}
}

func listenLocal(t *testing.T) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return listener
}

func waitRunning(t *testing.T, host *Host) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		host.mu.Lock()
		running := host.server != nil
		host.mu.Unlock()
		if running {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("host did not start")
}
