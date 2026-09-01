package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/middleware"
)

func TestRecoveryCatchesPanicNoHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/boom", nil)

	middleware.Recovery(nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
	if body := w.Body.String(); !strings.Contains(body, http.StatusText(http.StatusInternalServerError)) {
		t.Errorf("expected body to contain %q, got %q", http.StatusText(http.StatusInternalServerError), body)
	}
}

func TestRecoveryCatchesPanicWithHandler(t *testing.T) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/boom", nil)

	middleware.Recovery(handler)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})).ServeHTTP(w, r)

	if !called {
		t.Fatal("expected error handler to be called")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected Content-Type text/html, got %q", ct)
	}
}

func TestRecoveryNoPanic(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ok", nil)

	middleware.Recovery(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})).ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRecoveryWritesStatusBeforeHandler(t *testing.T) {
	var observed int
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		observed = w.(*httptest.ResponseRecorder).Code
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/boom", nil)

	middleware.Recovery(handler)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})).ServeHTTP(w, r)

	if observed != http.StatusInternalServerError {
		t.Errorf("expected WriteHeader(%d) before error handler, observed %d", http.StatusInternalServerError, observed)
	}
}
