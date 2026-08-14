package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
