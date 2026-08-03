package core

import (
	"net/http"
	"reflect"
	"testing"
)

func handlerPointer(h http.Handler) uintptr {
	return reflect.ValueOf(h).Pointer()
}

func TestServeMuxReturnsSameHandler(t *testing.T) {
	builtHandler = nil
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
	builtHandler = nil
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
