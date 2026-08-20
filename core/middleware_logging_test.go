package core

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestResponseWriterStatusDefault(t *testing.T) {
	rw := &responseWriter{ResponseWriter: httptest.NewRecorder()}
	if rw.status != 0 {
		t.Errorf("expected default status 0, got %d", rw.status)
	}
}

func TestResponseWriterWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}
	rw.WriteHeader(http.StatusNotFound)
	if rw.status != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rw.status)
	}
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected recorder status 404, got %d", rec.Code)
	}
}

func TestResponseWriterUnwrap(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}
	if got := rw.Unwrap(); got != rec {
		t.Errorf("expected Unwrap to return the underlying writer")
	}
}

func TestResponseControllerReachesUnderlyingWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}
	rc := http.NewResponseController(rw)
	if err := rc.Flush(); err != nil {
		t.Fatalf("ResponseController.Flush through Unwrap failed: %v", err)
	}
}

func TestJSONLHandlerOutput(t *testing.T) {
	var buf bytes.Buffer
	h := &jsonlHandler{w: &buf}
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "request", 0)
	rec.AddAttrs(slog.String("method", "GET"), slog.Int("status", 200))
	if err := h.Handle(context.Background(), rec); err != nil {
		t.Fatalf("handle failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["method"] != "GET" {
		t.Errorf("expected method GET, got %v", m["method"])
	}
	if m["status"] != float64(200) {
		t.Errorf("expected status 200, got %v", m["status"])
	}
}

func TestRequestLoggingSetsStatus(t *testing.T) {
	mw := RequestLogging()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/missing", nil)
	mw.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected downstream status 404, got %d", w.Code)
	}
}

func TestRequestLoggingIncludesRID(t *testing.T) {
	var gotRID string
	inner := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRID = RequestIDFromCtx(r.Context())
	}))
	mw := RequestLogging()(inner)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw.ServeHTTP(w, r)

	if gotRID == "" {
		t.Error("expected request ID in context through RequestLogging")
	}
}

func TestRequestLoggingPreservesFlusher(t *testing.T) {
	mw := RequestLogging()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := w.(http.Flusher); !ok {
			t.Error("expected http.Flusher through RequestLogging")
		}
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw.ServeHTTP(w, r)
}
