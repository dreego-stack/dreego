package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisteredMethodsRemainIsolated(t *testing.T) {
	a := New()
	if err := a.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	if err := a.Register("GET", "/resource", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("get"))
	}); err != nil {
		t.Fatal(err)
	}
	if err := a.Register("POST", "/resource", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("post"))
	}); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		method string
		status int
		body   string
	}{
		{"GET", http.StatusOK, "get"},
		{"POST", http.StatusCreated, "post"},
	} {
		t.Run(tc.method, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/resource", nil)
			rec := httptest.NewRecorder()
			a.Handler().ServeHTTP(rec, req)
			body, _ := io.ReadAll(rec.Result().Body)
			if rec.Code != tc.status || string(body) != tc.body {
				t.Fatalf("response: status=%d body=%q", rec.Code, body)
			}
		})
	}
}

func TestRegisterMethodWithoutHandlerReturnsMethodNotAllowed(t *testing.T) {
	a := New()
	if err := a.SetLogging(false); err != nil {
		t.Fatal(err)
	}
	if err := a.Register("POST", "/submit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}); err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	a.Handler().ServeHTTP(rec, httptest.NewRequest("GET", "/submit", nil))
	if rec.Code != http.StatusMethodNotAllowed && rec.Code != http.StatusNotFound {
		t.Fatalf("GET without handler: expected 405 or 404, got %d", rec.Code)
	}
}
