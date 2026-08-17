package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestListenAppliesServerTimeouts(t *testing.T) {
	app := New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	app.Register("GET", "/x", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	port := freePort(t)
	errCh := make(chan error, 1)
	go func() { errCh <- app.Listen(fmt.Sprintf("127.0.0.1:%d", port)) }()
	waitListen(t, port)

	app.mu.RLock()
	srv := app.server
	app.mu.RUnlock()
	if srv == nil {
		t.Fatal("app.server not set during Listen")
	}

	if srv.ReadHeaderTimeout == 0 {
		t.Error("ReadHeaderTimeout not set")
	}
	if srv.ReadTimeout == 0 {
		t.Error("ReadTimeout not set")
	}
	if srv.WriteTimeout == 0 {
		t.Error("WriteTimeout not set")
	}
	if srv.IdleTimeout == 0 {
		t.Error("IdleTimeout not set")
	}
	if srv.MaxHeaderBytes <= 0 {
		t.Error("MaxHeaderBytes not set")
	}

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	_ = proc.Signal(syscall.SIGINT)
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Listen returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Listen did not return after SIGINT")
	}
}

func TestListenWaitsForShutdownCompletion(t *testing.T) {
	app := New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	var inFlight sync.WaitGroup
	inFlight.Add(1)
	completed := atomic.Bool{}
	app.Register("GET", "/slow", func(w http.ResponseWriter, r *http.Request) {
		defer inFlight.Done()
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("done"))
		completed.Store(true)
	})

	port := freePort(t)
	errCh := make(chan error, 1)
	go func() { errCh <- app.Listen(fmt.Sprintf("127.0.0.1:%d", port)) }()
	waitListen(t, port)

	go func() {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/slow", port))
		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}()
	time.Sleep(80 * time.Millisecond)

	start := time.Now()
	proc, _ := os.FindProcess(os.Getpid())
	_ = proc.Signal(syscall.SIGINT)

	select {
	case err := <-errCh:
		elapsed := time.Since(start)
		if err != nil {
			t.Fatalf("Listen returned error: %v", err)
		}
		if elapsed < 100*time.Millisecond {
			t.Fatalf("Listen returned before request drained: %v", elapsed)
		}
		if !completed.Load() {
			t.Fatal("in-flight request did not complete before shutdown exit")
		}
	case <-time.After(15 * time.Second):
		t.Fatal("Listen did not return after shutdown")
	}
}

func TestShutdownExternalDrainsServer(t *testing.T) {
	app := New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	var inFlight sync.WaitGroup
	inFlight.Add(1)
	completed := atomic.Bool{}
	app.Register("GET", "/slow", func(w http.ResponseWriter, r *http.Request) {
		defer inFlight.Done()
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("done"))
		completed.Store(true)
	})

	port := freePort(t)
	errCh := make(chan error, 1)
	go func() { errCh <- app.Listen(fmt.Sprintf("127.0.0.1:%d", port)) }()
	waitListen(t, port)

	go func() {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/slow", port))
		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}()
	time.Sleep(80 * time.Millisecond)

	start := time.Now()
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutErr := make(chan error, 1)
	go func() { shutErr <- app.Shutdown(shutCtx) }()

	select {
	case err := <-errCh:
		elapsed := time.Since(start)
		if err != nil {
			t.Fatalf("Listen returned error: %v", err)
		}
		if elapsed < 100*time.Millisecond {
			t.Fatalf("Listen returned before request drained: %v", elapsed)
		}
		if !completed.Load() {
			t.Fatal("in-flight request did not complete before shutdown exit")
		}
	case <-time.After(15 * time.Second):
		t.Fatal("Listen did not return after external Shutdown")
	}

	select {
	case serr := <-shutErr:
		if serr != nil {
			t.Fatalf("Shutdown returned error: %v", serr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown did not return")
	}

	if err := app.Shutdown(context.Background()); err != nil {
		t.Fatalf("post-Listen Shutdown should be a no-op, got: %v", err)
	}

	app.mu.RLock()
	staleServer := app.server
	staleDone := app.shutdownDone
	app.mu.RUnlock()
	if staleServer != nil {
		t.Fatal("app.server not cleared after Listen returned")
	}
	if staleDone != nil {
		t.Fatal("app.shutdownDone not cleared after Listen returned")
	}
}

func TestListenPropagatesShutdownFailure(t *testing.T) {
	app := New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	cfg := DefaultServerConfig()
	cfg.ShutdownTimeout = 100 * time.Millisecond
	if err := app.SetServerConfig(cfg); err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	released := make(chan struct{})
	app.Register("GET", "/blocking", func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-released
		w.WriteHeader(http.StatusOK)
	})

	port := freePort(t)
	errCh := make(chan error, 1)
	go func() { errCh <- app.Listen(fmt.Sprintf("127.0.0.1:%d", port)) }()
	waitListen(t, port)

	go func() {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/blocking", port))
		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("blocking request did not start")
	}

	proc, _ := os.FindProcess(os.Getpid())
	_ = proc.Signal(syscall.SIGINT)

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected observable error when shutdown deadline exceeded")
		}
		close(released)
	case <-time.After(10 * time.Second):
		close(released)
		t.Fatal("Listen did not return after shutdown deadline exceeded")
	}
}

func TestListenNoGoroutineLeakAcrossLifecycles(t *testing.T) {
	app := New()
	if err := app.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	app.Register("GET", "/x", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	before := runtime.NumGoroutine()
	for i := 0; i < 3; i++ {
		port := freePort(t)
		errCh := make(chan error, 1)
		go func() { errCh <- app.Listen(fmt.Sprintf("127.0.0.1:%d", port)) }()
		waitListen(t, port)

		proc, _ := os.FindProcess(os.Getpid())
		_ = proc.Signal(syscall.SIGINT)
		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("Listen cycle %d error: %v", i, err)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("Listen cycle %d did not return", i)
		}
	}

	time.Sleep(200 * time.Millisecond)
	after := runtime.NumGoroutine()
	if after > before+2 {
		t.Fatalf("goroutine leak: before=%d after=%d", before, after)
	}
}
