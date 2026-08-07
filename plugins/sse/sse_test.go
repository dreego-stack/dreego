package sse

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	dreego "codeberg.org/dreego/dreego/core"
)

var _ dreego.Plugin = (*SSEPlugin)(nil)

type flushRecorder struct {
	*httptest.ResponseRecorder
}

func (r *flushRecorder) Flush() {}

func newFlushRecorder() *flushRecorder {
	return &flushRecorder{httptest.NewRecorder()}
}

func TestSSERoute(t *testing.T) {
	dreego.Reset()
	p := New()
	dreego.UsePlugin(p)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/sse", nil).WithContext(ctx)
	rec := newFlushRecorder()

	done := make(chan struct{})
	go func() {
		dreego.ServeMux().ServeHTTP(rec, req)
		close(done)
	}()

	deadline := time.After(2 * time.Second)
	for {
		select {
		case <-done:
			t.Fatal("handler returned before context cancel")
		case <-deadline:
			t.Fatal("timed out waiting for stream")
		default:
		}
		if rec.Header().Get("Content-Type") == "text/event-stream" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}

	cancel()
	<-done
}

func TestBroadcastReachesSubscriber(t *testing.T) {
	p := New()
	ch := p.hub.subscribe()
	defer p.hub.unsubscribe(ch)

	p.Broadcast("hello")
	select {
	case msg := <-ch:
		if msg != "hello" {
			t.Fatalf("msg = %q, want hello", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber did not receive broadcast")
	}
}

func TestSSEStreamsBroadcastWithGzip(t *testing.T) {
	dreego.Reset()
	p := New()
	dreego.UsePlugin(p)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/sse", nil).WithContext(ctx)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := newFlushRecorder()

	done := make(chan struct{})
	go func() {
		dreego.ServeMux().ServeHTTP(rec, req)
		close(done)
	}()

	deadline := time.After(2 * time.Second)
	for {
		if gzipBodyContains(rec.Body.Bytes(), ": connected") {
			break
		}
		select {
		case <-done:
			t.Fatal("handler returned before context cancel")
		case <-deadline:
			t.Fatal("timed out waiting for stream")
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}
	if enc := rec.Header().Get("Content-Encoding"); enc != "gzip" {
		t.Fatalf("Content-Encoding = %q, want gzip", enc)
	}

	p.Broadcast("gzip-ping")

	deadline = time.After(2 * time.Second)
	for {
		if gzipBodyContains(rec.Body.Bytes(), "data: gzip-ping") {
			break
		}
		select {
		case <-done:
			t.Fatal("handler returned before context cancel")
		case <-deadline:
			t.Fatal("timed out waiting for broadcast in gzip stream")
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	<-done

	if !gzipBodyContains(rec.Body.Bytes(), "data: gzip-ping") {
		t.Errorf("expected gzip-decompressed stream to contain 'data: gzip-ping'")
	}
}

func gzipBodyContains(b []byte, want string) bool {
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return false
	}
	defer zr.Close()
	body, _ := io.ReadAll(zr)
	return strings.Contains(string(body), want)
}

func TestSSEStreamsBroadcast(t *testing.T) {
	dreego.Reset()
	p := New()
	dreego.UsePlugin(p)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/sse", nil).WithContext(ctx)
	rec := newFlushRecorder()

	done := make(chan struct{})
	go func() {
		dreego.ServeMux().ServeHTTP(rec, req)
		close(done)
	}()

	deadline := time.After(2 * time.Second)
	for {
		if rec.Header().Get("Content-Type") == "text/event-stream" {
			break
		}
		select {
		case <-done:
			t.Fatal("handler returned before context cancel")
		case <-deadline:
			t.Fatal("timed out waiting for stream")
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}

	p.Broadcast("ping")

	deadline = time.After(2 * time.Second)
	for {
		if strings.Contains(rec.Body.String(), "data: ping") {
			break
		}
		select {
		case <-done:
			t.Fatal("handler returned before context cancel")
		case <-deadline:
			t.Fatal("timed out waiting for broadcast in stream")
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	<-done
}

func TestAssets(t *testing.T) {
	p := New()
	f, err := p.Assets().Open("sse.js")
	if err != nil {
		t.Fatalf("open sse.js: %v", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		t.Fatal("sse.js is empty")
	}
}
