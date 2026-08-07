package core

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCompressSkippedWithoutAccept(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "" {
		t.Errorf("expected no Content-Encoding, got %q", got)
	}
	if w.Body.String() != "plain" {
		t.Errorf("expected uncompressed body 'plain', got %q", w.Body.String())
	}
}

func TestCompressGzipApplied(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("compress me"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected Content-Encoding gzip, got %q", got)
	}
	zr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatalf("gzip reader failed: %v", err)
	}
	defer zr.Close()
	body, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "compress me" {
		t.Errorf("expected decompressed body 'compress me', got %q", string(body))
	}
}

func TestCompressPreservesFlusher(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(http.Flusher); !ok {
			t.Error("expected http.Flusher through Compress with gzip")
		}
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestGzipResponseWriterWritesGzip(t *testing.T) {
	var buf strings.Builder
	gw := gzip.NewWriter(&buf)
	gwr := &gzipResponseWriter{Writer: gw}

	n, err := gwr.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes written, got %d", n)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	zr, err := gzip.NewReader(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("gzip reader failed: %v", err)
	}
	defer zr.Close()
	body, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "hello" {
		t.Errorf("expected decompressed 'hello', got %q", string(body))
	}
}
