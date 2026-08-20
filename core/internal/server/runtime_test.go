package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func handlerPointer(h http.Handler) uintptr {
	return reflect.ValueOf(h).Pointer()
}

func TestHandlerReturnsSameHandler(t *testing.T) {
	app := New()
	h1 := app.Handler()
	h2 := app.Handler()
	if h1 == nil {
		t.Fatal("Handler returned nil")
	}
	if handlerPointer(h1) != handlerPointer(h2) {
		t.Error("Handler should return the same cached handler")
	}
}

func TestBuildIsIdempotent(t *testing.T) {
	app := New()
	app.Build()
	h1 := app.builtHandler
	app.Build()
	h2 := app.builtHandler
	if h1 == nil {
		t.Fatal("Build did not set handler")
	}
	if handlerPointer(h1) != handlerPointer(h2) {
		t.Error("Build should not replace an existing handler")
	}
}

func TestNewAppHasFreshState(t *testing.T) {
	app := New()
	app.Build()
	if app.builtHandler == nil {
		t.Fatal("Build did not set handler")
	}

	app.Register("GET", "/runtime-reset-test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	app2 := New()
	if app2.builtHandler != nil {
		t.Error("New app should not have a built handler")
	}

	app2.Register("GET", "/runtime-reset-test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h2 := app2.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/runtime-reset-test", nil)
	h2.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("fresh app should serve routes registered before build, got %d", rr.Code)
	}
}
