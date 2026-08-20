package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func containsStr(s, sub string) bool {
	return strings.Contains(s, sub)
}

func TestRecoveryUnderCompressPanicProducesValidResponse(t *testing.T) {
	stack := Recovery(nil)(Compress()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	stack.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if enc := w.Header().Get("Content-Encoding"); enc == "gzip" {
		t.Fatalf("panic response must NOT be gzip-encoded, got %q (body would be corrupt)", enc)
	}
	body := w.Body.String()
	if body == "" {
		t.Fatal("expected non-empty 500 body")
	}
	if isGzipStream(w.Body.Bytes()) {
		t.Fatalf("panic response body must be plain, not a gzip stream: %q", body)
	}
	if !containsStr(body, "Internal Server Error") {
		t.Errorf("expected plain 500 body, got %q", body)
	}
}

func TestRecoveryUnderCompressPanicWithHandlerProducesValidResponse(t *testing.T) {
	called := false
	errHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Write([]byte("custom error"))
	})
	stack := Recovery(errHandler)(Compress()(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	})))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	stack.ServeHTTP(w, r)

	if !called {
		t.Fatal("expected error handler to be called")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if enc := w.Header().Get("Content-Encoding"); enc == "gzip" {
		t.Fatalf("panic response must NOT be gzip-encoded, got %q", enc)
	}
	if w.Body.String() != "custom error" {
		t.Errorf("expected 'custom error', got %q", w.Body.String())
	}
}

func TestRecoveryUnderCompressNormalGzipResponse(t *testing.T) {
	stack := Recovery(nil)(Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("compress me"))
	})))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	stack.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip, got %q", got)
	}
	body, err := gzipBody(w.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "compress me" {
		t.Errorf("expected 'compress me', got %q", string(body))
	}
}

func TestRecoveryUnderCompressNoAcceptPlainResponse(t *testing.T) {
	stack := Recovery(nil)(Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	})))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	stack.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "" {
		t.Errorf("expected no encoding, got %q", got)
	}
	if w.Body.String() != "plain" {
		t.Errorf("expected 'plain', got %q", w.Body.String())
	}
}

func isGzipStream(b []byte) bool {
	if len(b) < 2 || b[0] != 0x1f || b[1] != 0x8b {
		return false
	}
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return false
	}
	defer zr.Close()
	_, err = io.ReadAll(zr)
	return err == nil
}
