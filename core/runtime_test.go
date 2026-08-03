package core

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func handlerPointer(h http.Handler) uintptr {
	return reflect.ValueOf(h).Pointer()
}

func TestServeMuxReturnsSameHandler(t *testing.T) {
	Reset()
	h1 := ServeMux()
	h2 := ServeMux()
	if h1 == nil {
		t.Fatal("ServeMux returned nil")
	}
	if handlerPointer(h1) != handlerPointer(h2) {
		t.Error("ServeMux should return the same cached handler")
	}
}

func TestBuildIsIdempotent(t *testing.T) {
	Reset()
	Build()
	h1 := builtHandler
	Build()
	h2 := builtHandler
	if h1 == nil {
		t.Fatal("Build did not set handler")
	}
	if handlerPointer(h1) != handlerPointer(h2) {
		t.Error("Build should not replace an existing handler")
	}
}

func TestResetClearsCache(t *testing.T) {
	Reset()
	Build()
	if builtHandler == nil {
		t.Fatal("Build did not set handler")
	}

	Register("GET", "/runtime-reset-test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	Reset()
	if builtHandler != nil {
		t.Error("Reset should set builtHandler to nil")
	}

	h2 := ServeMux()
	if h2 == nil {
		t.Fatal("ServeMux returned nil")
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/runtime-reset-test", nil)
	h2.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reset/ServeMux should build a fresh handler that includes routes registered after the previous build, got %d", rr.Code)
	}
}
