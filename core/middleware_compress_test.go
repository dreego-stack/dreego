package core

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
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
	body, err := gzipBody(w.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "compress me" {
		t.Errorf("expected decompressed 'compress me', got %q", string(body))
	}
}

func TestCompressRespectsGzipQ0(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip;q=0")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got == "gzip" {
		t.Fatalf("must NOT compress when gzip;q=0, got %q", got)
	}
	if w.Body.String() != "plain" {
		t.Errorf("expected uncompressed body, got %q", w.Body.String())
	}
}

func TestCompressRespectsGzipQ0Mixed(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip;q=0, br")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got == "gzip" {
		t.Fatalf("must NOT compress when gzip;q=0, got %q", got)
	}
}

func TestCompressSetsVary(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Vary"); got != "Accept-Encoding" {
		t.Errorf("expected Vary 'Accept-Encoding', got %q", got)
	}
}

func TestCompressVaryAppendsExisting(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Cookie")
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	vary := w.Header().Values("Vary")
	found := false
	for _, v := range vary {
		if v == "Accept-Encoding" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Vary to contain 'Accept-Encoding', got %v", vary)
	}
	if vary[0] != "Cookie" {
		t.Errorf("expected existing Vary preserved, got %v", vary)
	}
}

func TestCompressSetsVaryEvenWhenNotCompressing(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "br")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Vary"); got != "Accept-Encoding" {
		t.Errorf("expected Vary set even when not compressing, got %q", got)
	}
}

func TestCompressSkipsHEAD(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("body"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("HEAD", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got == "gzip" {
		t.Errorf("must NOT compress HEAD responses, got %q", got)
	}
	if w.Body.String() != "body" {
		t.Errorf("expected uncompressed HEAD body, got %q", w.Body.String())
	}
}

func TestCompressSkips204(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got == "gzip" {
		t.Errorf("must NOT compress 204, got %q", got)
	}
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestCompressSkips304(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got == "gzip" {
		t.Errorf("must NOT compress 304, got %q", got)
	}
	if w.Code != http.StatusNotModified {
		t.Errorf("expected 304, got %d", w.Code)
	}
}

func TestCompressSkipsPreEncoded(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "br")
		w.Write([]byte("already br"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "br" {
		t.Errorf("must NOT double-encode, expected 'br', got %q", got)
	}
	if w.Body.String() != "already br" {
		t.Errorf("expected body preserved, got %q", w.Body.String())
	}
}

func TestCompressSkipsAlreadyCompressedContentType(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte("binary png bytes"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got == "gzip" {
		t.Errorf("must NOT compress already-compressed content type image/png, got %q", got)
	}
	if w.Body.String() != "binary png bytes" {
		t.Errorf("expected body preserved, got %q", w.Body.String())
	}
}

func TestCompressRemovesContentLength(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "11")
		w.Write([]byte("compress me"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip, got %q", got)
	}
	if cl := w.Header().Get("Content-Length"); cl == "11" {
		t.Errorf("original Content-Length must be removed when compressed, got %q", cl)
	}
}

func TestCompressPreservesContentLengthWhenNotCompressed(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "11")
		w.Write([]byte("compress me"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "br")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "" {
		t.Errorf("expected no encoding, got %q", got)
	}
	if cl := w.Header().Get("Content-Length"); cl != "11" {
		t.Errorf("Content-Length must be preserved when not compressed, got %q", cl)
	}
}

func TestGzipResponseWriterWritesGzip(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	body, err := gzipBody(w.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "hello" {
		t.Errorf("expected decompressed 'hello', got %q", string(body))
	}
}

func gzipBody(r io.Reader) ([]byte, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(zr)
}
