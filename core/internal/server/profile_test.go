package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dreego-stack/dreego/internal/session"
)

func TestPatternPrefixStripsMuxSuffixes(t *testing.T) {
	cases := map[string]string{
		"/hooks":          "/hooks",
		"/hooks/{$}":      "/hooks",
		"/hooks/*":        "/hooks",
		"/users/{id}":     "/users",
		"/{$}":            "/",
		"":                "/",
		"/":               "/",
		"/blog/2024/{id}": "/blog/2024",
	}
	for in, want := range cases {
		if got := patternPrefix(in); got != want {
			t.Errorf("patternPrefix(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPathUnderPrefixBoundaries(t *testing.T) {
	cases := []struct {
		prefix, path string
		want         bool
	}{
		{"/", "/anything", true},
		{"/hooks", "/hooks", true},
		{"/hooks", "/hooks/github", true},
		{"/hooks", "/hooksx", false},
		{"/admin/secret", "/admin", false},
	}
	for _, tc := range cases {
		if got := pathUnderPrefix(tc.prefix, tc.path); got != tc.want {
			t.Errorf("pathUnderPrefix(%q, %q) = %v, want %v", tc.prefix, tc.path, got, tc.want)
		}
	}
}

func TestNoProfilesKeepsGlobalSessionAndCSRF(t *testing.T) {
	app := New()
	store := session.NewCookieStore([]byte("01234567890123456789012345678901"))
	if err := app.SetSessionStore(store); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		if session.StoreFromCtx(r.Context()) == nil {
			http.Error(w, "no store", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("global session store not applied: %d", rec.Code)
	}
}

func TestProfiledAppLeavesUnboundRoutesWithoutStore(t *testing.T) {
	app := New()
	store := session.NewCookieStore([]byte("01234567890123456789012345678901"))
	disabled := false
	if err := app.Profile("hooks", Profile{Session: store, CSRF: &disabled}); err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("/hooks", "hooks"); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("GET", "/marketing", func(w http.ResponseWriter, r *http.Request) {
		if session.StoreFromCtx(r.Context()) != nil {
			http.Error(w, "unexpected store", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest("GET", "/marketing", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("unbound route leaked session store: %d %s", rec.Code, rec.Body.String())
	}
}

func TestApplyProfileRequiresLeadingSlash(t *testing.T) {
	app := New()
	store := session.NewCookieStore([]byte("01234567890123456789012345678901"))
	if err := app.Profile("hooks", Profile{Session: store}); err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("hooks", "hooks"); err == nil {
		t.Fatal("ApplyProfile must reject a pattern without a leading slash")
	}
	if err := app.ApplyProfile("/hooks", "hooks"); err != nil {
		t.Fatalf("ApplyProfile must accept a rooted pattern: %v", err)
	}
}
