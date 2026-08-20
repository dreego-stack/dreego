package server

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSetReadyNoRace(t *testing.T) {
	app := New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			app.SetReady(true)
		}()
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/ready", nil)
			app.readyHandler().ServeHTTP(w, r)
		}()
	}
	wg.Wait()
}

func TestHealthHandlerOK(t *testing.T) {
	app := New()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	app.healthHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", w.Body.String())
	}
}

func TestReadyHandlerReady(t *testing.T) {
	app := New()
	app.SetReady(true)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ready", nil)
	app.readyHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "ready" {
		t.Errorf("expected body 'ready', got %q", w.Body.String())
	}
}

func TestReadyHandlerNotReady(t *testing.T) {
	app := New()
	app.SetReady(false)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ready", nil)
	app.readyHandler().ServeHTTP(w, r)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}
}

func TestReadyHandlerBackToReady(t *testing.T) {
	app := New()
	app.SetReady(false)
	app.SetReady(true)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ready", nil)
	app.readyHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 after SetReady(true), got %d", w.Code)
	}
	if w.Body.String() != "ready" {
		t.Errorf("expected body 'ready', got %q", w.Body.String())
	}
}
