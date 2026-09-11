package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompressMultiFlushProducesValidMultimember(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("chunk1"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		w.Write([]byte("chunk2"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip, got %q", got)
	}
	if members := countGzipMembers(w.Body.Bytes()); members < 2 {
		t.Fatalf("expected 2+ gzip members after multiple Flushes, got %d", members)
	}
	body, err := gzipBody(w.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "chunk1chunk2" {
		t.Errorf("expected 'chunk1chunk2', got %q", string(body))
	}
}

func TestCompressFlushInBypassMode(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "br")
		w.Write([]byte("already br"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		w.Write([]byte(" more"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "br" {
		t.Fatalf("expected 'br', got %q", got)
	}
	if w.Body.String() != "already br more" {
		t.Errorf("expected bypass body preserved through Flush, got %q", w.Body.String())
	}
}

func countGzipMembers(b []byte) int {
	n := 0
	for i := 0; i+1 < len(b); i++ {
		if b[i] == 0x1f && b[i+1] == 0x8b {
			n++
		}
	}
	return n
}
