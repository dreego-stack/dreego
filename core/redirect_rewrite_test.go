package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMatchRedirectExactNotPrefix(t *testing.T) {
	t.Parallel()
	rd := redirectRule{from: "/api", to: "/v2/api", status: http.StatusMovedPermanently}
	if _, ok := matchRedirect(rd, "/apiary"); ok {
		t.Fatal("exact rule /api must not match /apiary")
	}
	if _, ok := matchRedirect(rd, "/api/v1"); ok {
		t.Fatal("exact rule /api must not match /api/v1 (use /api/* for that)")
	}
	target, ok := matchRedirect(rd, "/api")
	if !ok || target != "/v2/api" {
		t.Fatalf("exact rule /api should match /api, got target=%q ok=%v", target, ok)
	}
}

func TestMatchRedirectWildcardSegmentBoundary(t *testing.T) {
	t.Parallel()
	rd := redirectRule{from: "/api/*", to: "/v2/*", status: http.StatusFound}
	if _, ok := matchRedirect(rd, "/apiary"); ok {
		t.Fatal("wildcard /api/* must not match /apiary (segment boundary)")
	}
	if _, ok := matchRedirect(rd, "/apiary/x"); ok {
		t.Fatal("wildcard /api/* must not match /apiary/x (segment boundary)")
	}
	target, ok := matchRedirect(rd, "/api")
	if !ok {
		t.Fatalf("wildcard /api/* should match base /api, got ok=%v", ok)
	}
	if target != "/v2" {
		t.Fatalf("base /api -> %q, want /v2", target)
	}
	target, ok = matchRedirect(rd, "/api/users/1")
	if !ok {
		t.Fatalf("wildcard /api/* should match /api/users/1")
	}
	if target != "/v2/users/1" {
		t.Fatalf("/api/users/1 -> %q, want /v2/users/1", target)
	}
}

func TestMatchRewriteExactNotPrefix(t *testing.T) {
	t.Parallel()
	rw := rewriteRule{from: "/api", to: "/v2/api"}
	if matchRewrite(rw, "/apiary") {
		t.Fatal("exact rewrite /api must not match /apiary")
	}
	if matchRewrite(rw, "/api/v1") {
		t.Fatal("exact rewrite /api must not match /api/v1")
	}
	if !matchRewrite(rw, "/api") {
		t.Fatal("exact rewrite /api should match /api")
	}
}

func TestMatchRewriteWildcardSegmentBoundary(t *testing.T) {
	t.Parallel()
	rw := rewriteRule{from: "/api/*", to: "/v2/*"}
	if matchRewrite(rw, "/apiary") {
		t.Fatal("wildcard rewrite /api/* must not match /apiary")
	}
	if matchRewrite(rw, "/apiary/x") {
		t.Fatal("wildcard rewrite /api/* must not match /apiary/x")
	}
	if !matchRewrite(rw, "/api") {
		t.Fatal("wildcard rewrite /api/* should match base /api")
	}
	if !matchRewrite(rw, "/api/users/1") {
		t.Fatal("wildcard rewrite /api/* should match /api/users/1")
	}
}

func TestApplyRewriteWildcard(t *testing.T) {
	t.Parallel()
	rw := rewriteRule{from: "/api/*", to: "/v2/*"}
	got, ok := applyRewrite(rw, "/api/users/1")
	if !ok {
		t.Fatal("applyRewrite should match /api/users/1")
	}
	if got != "/v2/users/1" {
		t.Fatalf("applyRewrite /api/users/1 -> %q, want /v2/users/1", got)
	}
	got, ok = applyRewrite(rw, "/api")
	if !ok {
		t.Fatal("applyRewrite should match base /api")
	}
	if got != "/v2" {
		t.Fatalf("applyRewrite /api -> %q, want /v2", got)
	}
	got, ok = applyRewrite(rw, "/apiary")
	if ok {
		t.Fatalf("applyRewrite must not match /apiary, got %q", got)
	}
}

func TestRedirectHTTPNotMatchNearPrefix(t *testing.T) {
	t.Parallel()
	app := New()
	app.Register("GET", "/apiary", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("apiary"))
	})
	app.Register("GET", "/v2/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("v2api"))
	})
	app.Register("GET", "/v2/users/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("v2user"))
	})
	if err := app.RegisterRedirect("/api", "/v2/api", http.StatusMovedPermanently); err != nil {
		t.Fatalf("exact redirect: %v", err)
	}
	if err := app.RegisterRedirect("/api/*", "/v2/*", http.StatusFound); err != nil {
		t.Fatalf("wildcard redirect: %v", err)
	}

	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/apiary", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || rr.Body.String() != "apiary" {
		t.Fatalf("/apiary should hit its route, got %d %q", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusMovedPermanently {
		t.Fatalf("/api should redirect (exact), got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/v2/api" {
		t.Fatalf("/api Location = %q, want /v2/api", loc)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/users/1", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusFound {
		t.Fatalf("/api/users/1 should redirect (wildcard), got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/v2/users/1" {
		t.Fatalf("/api/users/1 Location = %q, want /v2/users/1", loc)
	}
}

func TestRewriteHTTPNotMatchNearPrefix(t *testing.T) {
	t.Parallel()
	app := New()
	app.Register("GET", "/apiary", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("apiary"))
	})
	app.Register("GET", "/v2/users/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("v2user"))
	})
	if err := app.RegisterRewrite("/api/*", "/v2/*"); err != nil {
		t.Fatalf("rewrite: %v", err)
	}

	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/apiary", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || rr.Body.String() != "apiary" {
		t.Fatalf("/apiary should hit its route (not be rewritten), got %d %q", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/users/1", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || rr.Body.String() != "v2user" {
		t.Fatalf("/api/users/1 should be rewritten to /v2/users/1, got %d %q", rr.Code, rr.Body.String())
	}
}

func TestRedirectTrailingSlashCanonical(t *testing.T) {
	t.Parallel()
	app := New()
	app.Register("GET", "/v2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("v2"))
	})
	if err := app.RegisterRedirect("/api/*", "/v2/*", http.StatusMovedPermanently); err != nil {
		t.Fatalf("redirect: %v", err)
	}
	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusMovedPermanently {
		t.Fatalf("/api/ should redirect, got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/v2/" {
		t.Fatalf("/api/ Location = %q, want /v2/", loc)
	}
}

func TestRewriteExactMatchStillWorks(t *testing.T) {
	t.Parallel()
	app := New()
	app.Register("GET", "/new", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("new"))
	})
	if err := app.RegisterRewrite("/old", "/new"); err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	h := app.Handler()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/old", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK || rr.Body.String() != "new" {
		t.Fatalf("exact rewrite /old -> /new failed, got %d %q", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/old/sub", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("exact rewrite /old should NOT match /old/sub, got %d", rr.Code)
	}
}

func TestRedirectLocationNeverProtocolRelative(t *testing.T) {
	t.Parallel()
	app := New()
	if err := app.RegisterRedirect("/old", "/", http.StatusMovedPermanently); err != nil {
		t.Fatalf("exact redirect to root: %v", err)
	}
	if err := app.RegisterRedirect("/api/*", "/v2/*", http.StatusFound); err != nil {
		t.Fatalf("wildcard redirect: %v", err)
	}
	if err := app.RegisterRedirect("/blog/*", "/v2/api/*", http.StatusFound); err != nil {
		t.Fatalf("wildcard redirect to nested target: %v", err)
	}
	h := app.Handler()

	for _, p := range []string{"/old", "/api", "/api/users/1", "/blog/x", "/blog/y/z"} {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", p, nil)
		h.ServeHTTP(rr, req)
		if rr.Code < 300 || rr.Code > 399 {
			t.Fatalf("%s: expected redirect, got %d", p, rr.Code)
		}
		if loc := rr.Header().Get("Location"); strings.HasPrefix(loc, "//") {
			t.Fatalf("%s: protocol-relative Location %q must never be emitted", p, loc)
		}
	}
}
