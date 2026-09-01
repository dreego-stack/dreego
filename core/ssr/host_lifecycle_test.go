package ssr

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestStartIsSynchronousAndWaitObservesShutdown(t *testing.T) {
	host := New(dreego.New())
	if err := host.Start("127.0.0.1:0"); err != nil {
		t.Fatal(err)
	}
	host.mu.Lock()
	running := host.server != nil
	host.mu.Unlock()
	if !running {
		t.Fatal("Start returned before the host entered the running state")
	}
	if err := host.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := host.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestHostHandlesTerminationSignal(t *testing.T) {
	host := New(dreego.New())
	if err := host.Start("127.0.0.1:0"); err != nil {
		t.Fatal(err)
	}
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	if err := host.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestRepeatedHostLifecycleDoesNotLeakGoroutines(t *testing.T) {
	baseline := runtime.NumGoroutine()
	host := New(dreego.New())
	for cycle := 0; cycle < 10; cycle++ {
		if err := host.Start("127.0.0.1:0"); err != nil {
			t.Fatalf("cycle %d Start: %v", cycle, err)
		}
		if err := host.Shutdown(context.Background()); err != nil {
			t.Fatalf("cycle %d Shutdown: %v", cycle, err)
		}
	}
	runtime.GC()
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > baseline+3 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if current := runtime.NumGoroutine(); current > baseline+3 {
		t.Fatalf("goroutines after repeated lifecycle = %d, baseline = %d", current, baseline)
	}
}

func TestCompletedLifecycleResultIsIndependentFromNextLifecycle(t *testing.T) {
	oldErr := errors.New("old lifecycle")
	oldLifecycle := &hostLifecycle{done: make(chan struct{}), err: oldErr}
	close(oldLifecycle.done)
	host := New(dreego.New())
	host.lifecycle = &hostLifecycle{done: make(chan struct{})}

	result := make(chan error, 1)
	go func() { result <- host.wait(oldLifecycle) }()
	select {
	case err := <-result:
		if !errors.Is(err, oldErr) {
			t.Fatalf("wait returned %v, want old lifecycle error", err)
		}
	case <-time.After(time.Second):
		t.Fatal("wait for completed lifecycle blocked on the next lifecycle")
	}
}

func TestShutdownDrainsInFlightRequest(t *testing.T) {
	app := dreego.New()
	started := make(chan struct{})
	release := make(chan struct{})
	completed := make(chan struct{})
	if err := app.Register("GET", "/slow", func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusNoContent)
		close(completed)
	}); err != nil {
		t.Fatal(err)
	}

	host := New(app)
	listener := listenLocal(t)
	serveDone := make(chan error, 1)
	go func() { serveDone <- host.Serve(listener) }()
	waitRunning(t, host)
	go func() {
		response, err := http.Get("http://" + listener.Addr().String() + "/slow")
		if err == nil {
			io.Copy(io.Discard, response.Body)
			response.Body.Close()
		}
	}()
	<-started

	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- host.Shutdown(context.Background()) }()
	select {
	case err := <-shutdownDone:
		t.Fatalf("Shutdown returned before request completed: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	<-completed
	if err := <-shutdownDone; err != nil {
		t.Fatal(err)
	}
	if err := <-serveDone; err != nil {
		t.Fatal(err)
	}
}

func TestShutdownDeadlineIsObservable(t *testing.T) {
	app := dreego.New()
	started := make(chan struct{})
	release := make(chan struct{})
	if err := app.Register("GET", "/blocking", func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusNoContent)
	}); err != nil {
		t.Fatal(err)
	}

	host := New(app)
	listener := listenLocal(t)
	serveDone := make(chan error, 1)
	go func() { serveDone <- host.Serve(listener) }()
	waitRunning(t, host)
	go func() {
		response, requestErr := http.Get("http://" + listener.Addr().String() + "/blocking")
		if requestErr == nil {
			response.Body.Close()
		}
	}()
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err := host.Shutdown(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown = %v, want context deadline exceeded", err)
	}
	close(release)
	if err := <-serveDone; !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Serve = %v, want shutdown deadline", err)
	}
}

func TestHostSupportsRepeatedLifecycles(t *testing.T) {
	host := New(dreego.New())
	for cycle := 0; cycle < 2; cycle++ {
		listener := listenLocal(t)
		serveDone := make(chan error, 1)
		go func() { serveDone <- host.Serve(listener) }()
		waitRunning(t, host)
		if err := host.Shutdown(context.Background()); err != nil {
			t.Fatalf("cycle %d Shutdown: %v", cycle, err)
		}
		if err := <-serveDone; err != nil {
			t.Fatalf("cycle %d Serve: %v", cycle, err)
		}
	}
}

func TestSameAppCanUseIndependentHosts(t *testing.T) {
	app := dreego.New()
	first := New(app)
	second := New(app)
	firstListener := listenLocal(t)
	secondListener := listenLocal(t)
	firstDone := make(chan error, 1)
	secondDone := make(chan error, 1)
	go func() { firstDone <- first.Serve(firstListener) }()
	go func() { secondDone <- second.Serve(secondListener) }()
	waitRunning(t, first)
	waitRunning(t, second)
	if err := first.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := second.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := <-firstDone; err != nil {
		t.Fatal(err)
	}
	if err := <-secondDone; err != nil {
		t.Fatal(err)
	}
}
