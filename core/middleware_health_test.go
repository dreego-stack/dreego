package core

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSetReadyNoRace(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			SetReady(true)
		}()
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/ready", nil)
			readyHandler().ServeHTTP(w, r)
		}()
	}
	wg.Wait()
}

func TestHealthHandlerOK(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	healthHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", w.Body.String())
	}
}

func TestReadyHandlerReady(t *testing.T) {
	SetReady(true)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ready", nil)
	readyHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "ready" {
		t.Errorf("expected body 'ready', got %q", w.Body.String())
	}
}

func TestReadyHandlerNotReady(t *testing.T) {
	SetReady(false)
	defer SetReady(true)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ready", nil)
	readyHandler().ServeHTTP(w, r)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}
	if w.Body.String() != "not ready" {
		t.Errorf("expected body 'not ready', got %q", w.Body.String())
	}
}

func TestReadyHandlerBackToReady(t *testing.T) {
	SetReady(false)
	SetReady(true)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ready", nil)
	readyHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 after SetReady(true), got %d", w.Code)
	}
	if w.Body.String() != "ready" {
		t.Errorf("expected body 'ready', got %q", w.Body.String())
	}
}
