package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/middleware"
)

func TestCompressPanicDoesNotCorruptResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/boom", nil)
	r.Header.Set("Accept-Encoding", "gzip")

	stack := middleware.Recovery(nil)(middleware.Compress()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})))
	stack.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
	if body := w.Body.String(); body != "" && body != http.StatusText(http.StatusInternalServerError)+"\n" {
		t.Errorf("expected empty or generic error body, got %q", body)
	}
	if ce := w.Header().Get("Content-Encoding"); ce == "gzip" {
		t.Errorf("Content-Encoding: gzip must be removed when downstream panics, got %q", ce)
	}
}
