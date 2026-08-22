package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMaxBodyReaderLimitsRequestBody(t *testing.T) {
	var got error
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, got = io.ReadAll(r.Body)
	})
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("0123456789"))
	MaxBodyReader(4)(next).ServeHTTP(httptest.NewRecorder(), r)
	if got == nil {
		t.Fatal("ReadAll error = nil, want request body limit error")
	}
}

func TestMaxBodyReaderLeavesBodyUnchangedWhenDisabled(t *testing.T) {
	var body string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ReadAll() error = %v", err)
		}
		body = string(data)
	})
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("payload"))
	MaxBodyReader(0)(next).ServeHTTP(httptest.NewRecorder(), r)
	if body != "payload" {
		t.Fatalf("body = %q, want payload", body)
	}
}

func TestRequestIDPreservesIncomingIDAndContext(t *testing.T) {
	const want = "client-request-id"
	var got string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = RequestIDFromCtx(r.Context())
	})
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-Request-ID", want)
	w := httptest.NewRecorder()
	RequestID()(next).ServeHTTP(w, r)
	if got != want || w.Header().Get("X-Request-ID") != want {
		t.Fatalf("request ID = %q, header = %q, want %q", got, w.Header().Get("X-Request-ID"), want)
	}
}

func TestRequestIDGeneratesHexIdentifier(t *testing.T) {
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	w := httptest.NewRecorder()
	RequestID()(next).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	id := w.Header().Get("X-Request-ID")
	if len(id) != 16 {
		t.Fatalf("generated ID length = %d, want 16", len(id))
	}
	for _, c := range id {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Fatalf("generated ID %q contains non-hex character %q", id, c)
		}
	}
}
