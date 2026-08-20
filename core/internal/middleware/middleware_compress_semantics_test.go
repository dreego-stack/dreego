package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCompressStreamingFlush(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("chunk1"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		w.Write([]byte("chunk2"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip, got %q", got)
	}
	body, err := gzipBody(w.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "chunk1chunk2" {
		t.Errorf("expected 'chunk1chunk2', got %q", string(body))
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

func TestCompressPreservesHijacker(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(http.Hijacker); !ok {
			t.Error("expected http.Hijacker through Compress when supported by upstream")
		}
	}))

	w := newHijackableRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestCompressPreservesPusher(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(http.Pusher); !ok {
			t.Error("expected http.Pusher through Compress when supported by upstream")
		}
	}))

	w := newPusherRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestCompressPreservesReaderFrom(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(io.ReaderFrom); !ok {
			t.Error("expected io.ReaderFrom through Compress when supported by upstream")
		}
	}))

	w := newReaderFromRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestCompressRespectsGzipQ0WithWildcard(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip;q=0, *")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got == "gzip" {
		t.Fatalf("explicit gzip;q=0 must beat wildcard, got %q", got)
	}
	if w.Body.String() != "plain" {
		t.Errorf("expected uncompressed body, got %q", w.Body.String())
	}
}

func TestCompressAcceptsGzipCaseInsensitive(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("compress me"))
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "GZip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip for 'GZip', got %q", got)
	}
	body, err := gzipBody(w.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "compress me" {
		t.Errorf("expected 'compress me', got %q", string(body))
	}
}

func TestCompressCopiesViaIO(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		src := io.LimitReader(strings.NewReader("copy me"), 7)
		if _, err := io.Copy(w, src); err != nil {
			t.Errorf("io.Copy failed: %v", err)
		}
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)

	if got := w.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip through io.Copy, got %q", got)
	}
	body, err := gzipBody(w.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "copy me" {
		t.Errorf("expected 'copy me', got %q", string(body))
	}
}

func TestCompressHijackForwards(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("expected http.Hijacker")
		}
		if _, _, err := h.Hijack(); err != nil {
			t.Errorf("Hijack must forward to upstream, got error %v", err)
		}
	}))

	w := newHijackableRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestCompressPushForwards(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := w.(http.Pusher)
		if !ok {
			t.Fatal("expected http.Pusher")
		}
		if err := p.Push("/asset.css", nil); err != nil {
			t.Errorf("Push must forward to upstream, got error %v", err)
		}
	}))

	w := newPusherRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestCompressHijackUnsupported(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("expected http.Hijacker")
		}
		if _, _, err := h.Hijack(); err != http.ErrNotSupported {
			t.Errorf("expected http.ErrNotSupported without upstream support, got %v", err)
		}
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestCompressPushUnsupported(t *testing.T) {
	mw := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := w.(http.Pusher)
		if !ok {
			t.Fatal("expected http.Pusher")
		}
		if err := p.Push("/asset.css", nil); err != http.ErrNotSupported {
			t.Errorf("expected http.ErrNotSupported without upstream support, got %v", err)
		}
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	mw.ServeHTTP(w, r)
}

func TestCompressForwardsInformationalStatus(t *testing.T) {
	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", "</style.css>; rel=preload; as=style")
		w.WriteHeader(http.StatusEarlyHints)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, err := http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected final status 404, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("expected gzip, got %q", got)
	}
	body, err := gzipBody(resp.Body)
	if err != nil {
		t.Fatalf("read gzip body failed: %v", err)
	}
	if string(body) != "not found" {
		t.Errorf("expected 'not found', got %q", string(body))
	}
}
