package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func isGzipStream(b []byte) bool {
	return len(b) >= 2 && b[0] == 0x1f && b[1] == 0x8b
}

func TestAppStackPanicUnderCompression(t *testing.T) {
	app := New()
	if err := app.Register("GET", "/boom", func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/boom", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	app.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if enc := w.Header().Get("Content-Encoding"); enc == "gzip" {
		t.Fatalf("panic response must not be gzip, got %q", enc)
	}
	body := w.Body.String()
	if body == "" || isGzipStream(w.Body.Bytes()) {
		t.Fatalf("expected plain 500 body, got %q", body)
	}
	if got := w.Header().Get("Vary"); !strings.Contains(got, "Accept-Encoding") {
		t.Errorf("expected Vary Accept-Encoding on panic response, got %q", got)
	}
}
