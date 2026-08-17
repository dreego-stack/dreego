package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

func waitListen(t *testing.T, port int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("server on port %d did not start", port)
}

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

func TestMaxBodyReader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(io.Discard, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	wrapped := MaxBodyReader(100)(http.HandlerFunc(handler))

	oversized := make([]byte, 200)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/upload", bytes.NewReader(oversized))
	req.ContentLength = int64(len(oversized))
	wrapped.ServeHTTP(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized: expected 413, got %d", rec.Code)
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/upload", bytes.NewReader([]byte("hello")))
	req2.ContentLength = 5
	wrapped.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("under limit: expected 200, got %d", rec2.Code)
	}
}

func TestDefaultServerTimeoutsAreSecure(t *testing.T) {
	c := DefaultServerConfig()
	cases := map[string]struct {
		got  any
		zero any
	}{
		"ReadHeaderTimeout": {c.ReadHeaderTimeout, time.Duration(0)},
		"ReadTimeout":       {c.ReadTimeout, time.Duration(0)},
		"WriteTimeout":      {c.WriteTimeout, time.Duration(0)},
		"IdleTimeout":       {c.IdleTimeout, time.Duration(0)},
		"MaxHeaderBytes":    {c.MaxHeaderBytes, 0},
		"ShutdownTimeout":   {c.ShutdownTimeout, time.Duration(0)},
	}
	for name, tc := range cases {
		if tc.got == tc.zero {
			t.Errorf("%s default not set", name)
		}
	}
}

func TestSetServerConfigBeforeBuild(t *testing.T) {
	app := New()
	custom := DefaultServerConfig()
	custom.ReadHeaderTimeout = 5 * time.Second
	custom.ShutdownTimeout = 3 * time.Second
	if err := app.SetServerConfig(custom); err != nil {
		t.Fatalf("SetServerConfig: %v", err)
	}
	app.Build()
	if err := app.SetServerConfig(custom); !errors.Is(err, ErrAppBuilt) {
		t.Fatalf("expected ErrAppBuilt after build, got %v", err)
	}
}
