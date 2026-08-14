package core

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func okHandler(http.ResponseWriter, *http.Request) {}

func TestNewAppReturnsNonNil(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestAppServesRegisteredRoute(t *testing.T) {
	app := New()
	app.Register("GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
	app.Build()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hello", nil)
	app.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello" {
		t.Fatalf("expected body 'hello', got %q", rec.Body.String())
	}
}

func TestTwoAppsAreIsolated(t *testing.T) {
	app1 := New()
	app1.Register("GET", "/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("app1"))
	})
	app1.Build()

	app2 := New()
	app2.Register("GET", "/b", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("app2"))
	})
	app2.Build()

	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest("GET", "/a", nil)
	app1.Handler().ServeHTTP(rec1, req1)
	if rec1.Body.String() != "app1" {
		t.Fatalf("app1 /a: expected 'app1', got %q", rec1.Body.String())
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/b", nil)
	app2.Handler().ServeHTTP(rec2, req2)
	if rec2.Body.String() != "app2" {
		t.Fatalf("app2 /b: expected 'app2', got %q", rec2.Body.String())
	}

	recMiss := httptest.NewRecorder()
	reqMiss := httptest.NewRequest("GET", "/a", nil)
	app2.Handler().ServeHTTP(recMiss, reqMiss)
	if recMiss.Code != http.StatusNotFound {
		t.Fatalf("app2 /a: expected 404 (route only in app1), got %d", recMiss.Code)
	}
}

func TestAppHasHealthAndReady(t *testing.T) {
	app := New()
	app.Build()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	app.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected /health 200, got %d", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	if string(body) != "ok" {
		t.Fatalf("expected /health body 'ok', got %q", string(body))
	}
}

func TestAppSetCSRF(t *testing.T) {
	app := New()
	app.SetCSRF(false)
	app.Build()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/anything", nil)
	app.Handler().ServeHTTP(rec, req)

	if rec.Code == http.StatusForbidden {
		t.Fatal("CSRF should be disabled but got 403")
	}
}

func TestAppSetErrorHandler(t *testing.T) {
	app := New()
	app.SetErrorHandler(500, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("custom error"))
	})
	app.Build()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	app.Handler().ServeHTTP(rec, req)
}

func TestAppRejectsConfigurationAfterBuild(t *testing.T) {
	app := New()
	app.Build()

	checks := []struct {
		name string
		err  error
	}{
		{"route", app.Register("GET", "/late", okHandler)},
		{"redirect", app.RegisterRedirect("/old", "/new", http.StatusFound)},
		{"rewrite", app.RegisterRewrite("/old/*", "/new/*")},
		{"middleware", app.Use(func(next http.Handler) http.Handler { return next })},
		{"logging", app.SetLogging(false)},
		{"csrf", app.SetCSRF(false)},
		{"error handler", app.SetErrorHandler(500, okHandler)},
		{"session", app.SetSessionStore(mockStore{})},
		{"csp", app.SetCSP("default-src 'none'")},
		{"rule", app.RegisterRule("late", func(string) string { return "" })},
	}

	for _, check := range checks {
		if !errors.Is(check.err, ErrAppBuilt) {
			t.Errorf("%s mutation error = %v, want ErrAppBuilt", check.name, check.err)
		}
	}
}

func TestBuildDoesNotHoldConfigLockWhileBuildingMiddleware(t *testing.T) {
	app := New()
	mutationResult := make(chan error, 1)
	if err := app.Use(func(next http.Handler) http.Handler {
		mutationResult <- app.SetLogging(false)
		return next
	}); err != nil {
		t.Fatal(err)
	}

	built := make(chan struct{})
	go func() {
		app.Handler()
		close(built)
	}()

	select {
	case <-built:
	case <-time.After(time.Second):
		t.Fatal("building middleware deadlocked while configuring the app")
	}

	if err := <-mutationResult; !errors.Is(err, ErrAppBuilt) {
		t.Fatalf("middleware mutation error = %v, want ErrAppBuilt", err)
	}
}

func TestMiddlewareCanAccessHandlerDuringBuild(t *testing.T) {
	app := New()
	nestedHandler := make(chan http.Handler, 1)
	if err := app.Use(func(next http.Handler) http.Handler {
		nestedHandler <- app.Handler()
		return next
	}); err != nil {
		t.Fatal(err)
	}

	built := make(chan struct{})
	go func() {
		app.Handler()
		close(built)
	}()

	select {
	case <-built:
	case <-time.After(time.Second):
		t.Fatal("middleware access to Handler deadlocked during build")
	}

	var handler http.Handler
	select {
	case handler = <-nestedHandler:
	case <-time.After(time.Second):
		t.Fatal("middleware did not receive a handler")
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("nested handler status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestWaitingHandlerRecoversAfterBuildPanic(t *testing.T) {
	app := New()
	entered := make(chan struct{})
	release := make(chan struct{})
	var builds atomic.Int32
	if err := app.Use(func(next http.Handler) http.Handler {
		if builds.Add(1) == 1 {
			close(entered)
			<-release
			panic("build failed")
		}
		return next
	}); err != nil {
		t.Fatal(err)
	}

	panicked := make(chan struct{})
	go func() {
		defer func() {
			recover()
			close(panicked)
		}()
		app.Build()
	}()
	<-entered

	waiting := app.Handler()
	close(release)
	<-panicked

	done := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		rec := httptest.NewRecorder()
		waiting.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
		done <- rec
	}()

	select {
	case rec := <-done:
		if rec.Code != http.StatusOK {
			t.Fatalf("recovered handler status = %d, want %d", rec.Code, http.StatusOK)
		}
	case <-time.After(time.Second):
		t.Fatal("waiting handler blocked after build panic")
	}
}

func TestAppRejectsDuplicateAndReservedRoutes(t *testing.T) {
	app := New()
	if err := app.Register("GET", "/same", okHandler); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("GET", "/same", okHandler); !errors.Is(err, ErrRouteConflict) {
		t.Fatalf("duplicate route error = %v, want ErrRouteConflict", err)
	}
	if err := app.Register("GET", "/health", okHandler); !errors.Is(err, ErrRouteConflict) {
		t.Fatalf("reserved route error = %v, want ErrRouteConflict", err)
	}
}
