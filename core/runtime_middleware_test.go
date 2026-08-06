package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStore struct{}

func (mockStore) Get(*http.Request, string) (string, error) { return "", nil }
func (mockStore) Set(http.ResponseWriter, *http.Request, string, string, *Options) error {
	return nil
}
func (mockStore) Delete(http.ResponseWriter, *http.Request, string) error { return nil }
func (mockStore) Destroy(http.ResponseWriter, *http.Request) error        { return nil }

func TestRegisterRedirect(t *testing.T) {
	Reset()
	RegisterRedirect("/old", "/new", http.StatusMovedPermanently)
	if len(redirects) != 1 {
		t.Fatalf("expected 1 redirect rule, got %d", len(redirects))
	}
	if redirects[0].from != "/old" || redirects[0].to != "/new" || redirects[0].status != http.StatusMovedPermanently {
		t.Errorf("redirect rule not registered correctly: %+v", redirects[0])
	}
}

func TestRedirectServesRedirect(t *testing.T) {
	Reset()
	Register("GET", "/redirect-target", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	RegisterRedirect("/old", "/redirect-target", http.StatusMovedPermanently)

	h := ServeMux()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/old", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("expected status %d, got %d", http.StatusMovedPermanently, rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/redirect-target" {
		t.Errorf("expected Location %q, got %q", "/redirect-target", loc)
	}
}

func TestRedirectWildcardMatch(t *testing.T) {
	Reset()
	Register("GET", "/wild-target", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	RegisterRedirect("/blog/*", "/wild-target", http.StatusFound)

	h := ServeMux()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/blog", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("expected status %d for wildcard base match, got %d", http.StatusFound, rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/wild-target" {
		t.Errorf("expected Location %q, got %q", "/wild-target", loc)
	}
}

func TestRedirectWildcardPrefixMatch(t *testing.T) {
	// A "/blog/*" rule prefix-matches deeper paths: "/blog/2024/01" is
	// redirected, and the matched prefix is replaced by the target.
	Reset()
	RegisterRedirect("/blog/*", "/new", http.StatusFound)
	Register("GET", "/blog/2024/01", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := ServeMux()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/blog/2024/01", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("expected status %d for wildcard prefix match, got %d", http.StatusFound, rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/new/2024/01" {
		t.Errorf("expected Location %q, got %q", "/new/2024/01", loc)
	}
}

func TestRedirectNoMatch(t *testing.T) {
	Reset()
	RegisterRedirect("/old", "/new", http.StatusMovedPermanently)
	Register("GET", "/keep", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := ServeMux()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/keep", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected next handler to run for non-matching path, got %d", rr.Code)
	}
}

func TestRegisterRewrite(t *testing.T) {
	Reset()
	RegisterRewrite("/old/*", "/new/*")
	if len(rewrites) != 1 {
		t.Fatalf("expected 1 rewrite rule, got %d", len(rewrites))
	}
	if rewrites[0].from != "/old/*" || rewrites[0].to != "/new/*" {
		t.Errorf("rewrite rule not registered correctly: %+v", rewrites[0])
	}
}

func TestRewriteRewritesPath(t *testing.T) {
	Reset()
	Register("GET", "/new/route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	RegisterRewrite("/old/*", "/new/*")

	h := ServeMux()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/old/route", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected rewritten path to hit the route, got %d", rr.Code)
	}
}

func TestRewriteNoMatch(t *testing.T) {
	Reset()
	Register("GET", "/unchanged", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	RegisterRewrite("/old/*", "/new/*")

	h := ServeMux()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/unchanged", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected unmodified path to hit its route, got %d", rr.Code)
	}
}

func TestMatchRewrite(t *testing.T) {
	exact := matchRewrite(rewriteRule{from: "/foo", to: "/bar"}, "/foo")
	if !exact {
		t.Error("matchRewrite should match an exact path")
	}

	prefix := matchRewrite(rewriteRule{from: "/foo/*", to: "/bar/*"}, "/foo/deep/nested")
	if !prefix {
		t.Error("matchRewrite should prefix-match a wildcard rule")
	}

	trimmed := matchRewrite(rewriteRule{from: "/foo/*", to: "/bar/*"}, "/foo")
	if !trimmed {
		t.Error("matchRewrite should match the trimmed prefix itself")
	}
}

func TestSessionMiddlewareSetsStore(t *testing.T) {
	Reset()
	Register("GET", "/session-test", func(w http.ResponseWriter, r *http.Request) {
		got := StoreFromCtx(r.Context())
		if got == nil {
			http.Error(w, "no store", http.StatusInternalServerError)
			return
		}
		if _, ok := got.(mockStore); !ok {
			http.Error(w, "wrong store type", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	SetSessionStore(mockStore{})

	h := ServeMux()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/session-test", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected session store to be available in handler, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestSetSessionStore(t *testing.T) {
	Reset()
	SetSessionStore(mockStore{})
	if got := SessionStore(); got == nil {
		t.Fatal("SetSessionStore did not set the session store")
	}
	if _, ok := SessionStore().(mockStore); !ok {
		t.Errorf("expected mockStore, got %T", SessionStore())
	}
}

func TestSessionStore(t *testing.T) {
	Reset()
	if SessionStore() != nil {
		t.Errorf("expected nil store after Reset, got %T", SessionStore())
	}
	SetSessionStore(mockStore{})
	if SessionStore() == nil {
		t.Error("expected non-nil store after SetSessionStore")
	}
}

func TestSetLoggingDisablesLogging(t *testing.T) {
	Reset()
	if !loggingEnabled {
		t.Fatal("logging should default to enabled")
	}
	SetLogging(false)
	if loggingEnabled {
		t.Error("SetLogging(false) should disable logging")
	}
	SetLogging(true)
	if !loggingEnabled {
		t.Error("SetLogging(true) should re-enable logging")
	}
}

func TestSetCSRF(t *testing.T) {
	Reset()
	if !csrfEnabled {
		t.Fatal("CSRF should default to enabled")
	}
	SetCSRF(false)
	if csrfEnabled {
		t.Error("SetCSRF(false) should disable CSRF")
	}
	SetCSRF(true)
	if !csrfEnabled {
		t.Error("SetCSRF(true) should re-enable CSRF")
	}
}

func TestSetErrorHandler(t *testing.T) {
	Reset()
	if _, ok := errorHandlers[http.StatusInternalServerError]; ok {
		t.Fatal("error handler should be unset after Reset")
	}
	SetErrorHandler(http.StatusInternalServerError, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("custom 500"))
	})
	if h, ok := errorHandlers[http.StatusInternalServerError]; !ok || h == nil {
		t.Error("SetErrorHandler did not register the 500 handler")
	}
}
