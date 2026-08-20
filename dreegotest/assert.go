package dreegotest

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// MustStatus fails the test unless the response status code equals want.
func MustStatus(t testing.TB, code, want int) {
	t.Helper()
	if code != want {
		t.Fatalf("status = %d, want %d", code, want)
	}
}

// MustContainBody fails the test unless body contains substr.
func MustContainBody(t testing.TB, body, substr string) {
	t.Helper()
	if !strings.Contains(body, substr) {
		t.Fatalf("body missing %q\n---\n%s", substr, body)
	}
}

// MustNotContainBody fails the test if body contains substr.
func MustNotContainBody(t testing.TB, body, substr string) {
	t.Helper()
	if strings.Contains(body, substr) {
		t.Fatalf("body must not contain %q\n---\n%s", substr, body)
	}
}

// MustHeader fails the test unless the header key has value want.
func MustHeader(t testing.TB, h http.Header, key, want string) {
	t.Helper()
	got := h.Get(key)
	if got != want {
		t.Fatalf("header %s = %q, want %q", key, got, want)
	}
}

// MustEqual fails the test unless got equals want, printing a %#v diff.
func MustEqual[T any](t testing.TB, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

// MustNotEqual fails the test if got equals want.
func MustNotEqual[T any](t testing.TB, got, want T) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want something else", got)
	}
}

// MustRedirect fails the test unless status is 3xx and the Location header
// equals wantLocation.
func MustRedirect(t testing.TB, status int, h http.Header, wantLocation string) {
	t.Helper()
	if status < 300 || status >= 400 {
		t.Fatalf("status = %d, want redirect (3xx)", status)
	}
	if got := h.Get("Location"); got != wantLocation {
		t.Fatalf("Location = %q, want %q", got, wantLocation)
	}
}
