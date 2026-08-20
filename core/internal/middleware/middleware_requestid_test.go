package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestRequestIDGenerates(t *testing.T) {
	var gotHeader, gotCtx string
	mw := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = w.Header().Get("X-Request-ID")
		gotCtx = RequestIDFromCtx(r.Context())
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw.ServeHTTP(w, r)

	if gotHeader == "" {
		t.Fatal("expected X-Request-ID header to be set")
	}
	if gotCtx == "" {
		t.Fatal("expected request ID in context")
	}
	if gotHeader != gotCtx {
		t.Errorf("header ID %q != context ID %q", gotHeader, gotCtx)
	}
	if !regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(gotHeader) {
		t.Errorf("expected 16 hex chars, got %q", gotHeader)
	}
}

func TestRequestIDPreservesIncoming(t *testing.T) {
	var gotCtx string
	mw := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCtx = RequestIDFromCtx(r.Context())
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Request-ID", "incoming-id-123")
	mw.ServeHTTP(w, r)

	if gotCtx != "incoming-id-123" {
		t.Errorf("expected incoming ID preserved, got %q", gotCtx)
	}
	if got := w.Header().Get("X-Request-ID"); got != "incoming-id-123" {
		t.Errorf("expected response header to echo incoming ID, got %q", got)
	}
}

func TestRequestIDFromCtxMissing(t *testing.T) {
	if got := RequestIDFromCtx(context.Background()); got != "" {
		t.Errorf("expected empty string for missing ID, got %q", got)
	}
}

func TestRequestIDFromCtxWrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), RequestIDKey, 42)
	if got := RequestIDFromCtx(ctx); got != "" {
		t.Errorf("expected empty string for wrong type, got %q", got)
	}
}
