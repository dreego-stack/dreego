package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectEncodedSlashInPath(t *testing.T) {
	t.Parallel()
	app := New()
	app.Register("GET", "/v2/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hit:" + r.URL.Path))
	})
	if err := app.RegisterRedirect("/api/*", "/v2/*", http.StatusFound); err != nil {
		t.Fatalf("redirect: %v", err)
	}
	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/users%2F1", nil)
	if req.URL.Path != "/api/users/1" {
		t.Fatalf("net/http decodes %%2F: r.URL.Path = %q, want /api/users/1", req.URL.Path)
	}
	if req.URL.RawPath != "/api/users%2F1" {
		t.Fatalf("r.URL.RawPath = %q, want /api/users%%2F1", req.URL.RawPath)
	}
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", rr.Code)
	}
	loc := rr.Header().Get("Location")
	if loc != "/v2/users/1" {
		t.Fatalf("Location = %q, want /v2/users/1 (net/http decodes %%2F before the middleware sees Path)", loc)
	}
}

func TestRedirectDotSegmentsInPath(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/api/*", "/v2/*", http.StatusFound); err != nil {
		t.Fatalf("redirect: %v", err)
	}
	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/../etc", nil)
	if req.URL.Path != "/api/../etc" {
		t.Fatalf("net/http does not clean dot segments in Path: got %q, want /api/../etc", req.URL.Path)
	}
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", rr.Code)
	}
	loc := rr.Header().Get("Location")
	if loc != "/etc" {
		t.Fatalf("Location = %q, want /etc (middleware computes /v2/../etc, http.Redirect cleans it)", loc)
	}
}

func TestRedirectDoubleSlashInPath(t *testing.T) {
	t.Parallel()
	app := New()
	app.Register("GET", "/v2/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hit:" + r.URL.Path))
	})
	if err := app.RegisterRedirect("/api/*", "/v2/*", http.StatusFound); err != nil {
		t.Fatalf("redirect: %v", err)
	}
	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api//users", nil)
	if req.URL.Path != "/api//users" {
		t.Fatalf("net/http keeps double slash: got %q, want /api//users", req.URL.Path)
	}
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected redirect, got %d", rr.Code)
	}
	loc := rr.Header().Get("Location")
	if loc != "/v2/users" {
		t.Fatalf("Location = %q, want /v2/users (middleware emits /v2//users, http.Redirect collapses it)", loc)
	}
}

func TestRedirectDoubleSlashSuffixInternally(t *testing.T) {
	t.Parallel()
	rd := redirectRule{from: "/api/*", to: "/v2/*", status: http.StatusFound}
	target, ok := matchRedirect(rd, "/api//users")
	if !ok {
		t.Fatal("wildcard /api/* must match /api//users (HasPrefix /api/)")
	}
	if target != "/v2//users" {
		t.Fatalf("matchRedirect /api//users -> %q, want /v2//users (raw suffix preserved before Redirect cleans)", target)
	}
}

func TestRedirectRootPathDoesNotMatchNonRootRule(t *testing.T) {
	t.Parallel()
	app := New()
	app.Register("GET", "/v2/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	if err := app.RegisterRedirect("/api/*", "/v2/*", http.StatusFound); err != nil {
		t.Fatalf("redirect: %v", err)
	}
	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)

	if rr.Code == http.StatusFound || rr.Code == http.StatusMovedPermanently {
		t.Fatalf("/ must not match /api/*, got redirect %d Location=%q", rr.Code, rr.Header().Get("Location"))
	}
}
