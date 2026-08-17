package core

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMaxBodyReader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(io.Discard, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	wrapped := MaxBodyReader(100)(http.HandlerFunc(handler))

	oversized := make([]byte, 200)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/upload", bytes.NewReader(oversized))
	req.ContentLength = int64(len(oversized))
	wrapped.ServeHTTP(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized: expected 413, got %d", rec.Code)
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/upload", bytes.NewReader([]byte("hello")))
	req2.ContentLength = 5
	wrapped.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("under limit: expected 200, got %d", rec2.Code)
	}
}
