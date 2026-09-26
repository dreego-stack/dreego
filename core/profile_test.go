package core

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/session"
)

func boolPtr(v bool) *bool { return &v }

func TestApplyProfileRejectsUnknownName(t *testing.T) {
	app := New()
	if err := app.ApplyProfile("/hooks", "missing"); err != nil {
		t.Fatalf("ApplyProfile before registration should defer validation: %v", err)
	}
	err := app.Build()
	if err == nil {
		t.Fatal("expected Build error for unknown profile name")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Fatalf("error must name the profile, got %v", err)
	}
}

func TestProfileMutationsRejectedAfterBuild(t *testing.T) {
	app := New()
	if err := app.Profile("hooks", Profile{Session: NewCookieStore(testSecret)}); err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("/hooks", "hooks"); err != nil {
		t.Fatal(err)
	}
	app.Build()
	if err := app.ApplyProfile("/other", "hooks"); !errors.Is(err, ErrAppBuilt) {
		t.Fatalf("post-build ApplyProfile = %v, want ErrAppBuilt", err)
	}
	if err := app.Profile("late", Profile{Session: NewCookieStore(testSecret)}); !errors.Is(err, ErrAppBuilt) {
		t.Fatalf("post-build Profile = %v, want ErrAppBuilt", err)
	}
}

func TestProfileRejectsDuplicateAndEmptyName(t *testing.T) {
	app := New()
	store := NewCookieStore(testSecret)
	if err := app.Profile("hooks", Profile{Session: store}); err != nil {
		t.Fatalf("first Profile: %v", err)
	}
	if err := app.Profile("hooks", Profile{Session: store}); err == nil {
		t.Fatal("expected duplicate profile error")
	}
	if err := app.Profile("  ", Profile{Session: store}); err == nil {
		t.Fatal("expected empty profile name error")
	}
}

func TestProfileDisablesCSRFOnlyForBoundRoutes(t *testing.T) {
	global := NewCookieStore(testSecret)
	hookStore := NewCookieStore(testSecret)
	app := New()
	if err := app.SetSessionStore(global); err != nil {
		t.Fatal(err)
	}
	if err := app.Profile("hooks", Profile{Session: hookStore, CSRF: boolPtr(false)}); err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("/hooks", "hooks"); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("POST", "/hooks/github", okHandler); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("POST", "/form", okHandler); err != nil {
		t.Fatal(err)
	}

	h := app.Handler()

	hook := httptest.NewRecorder()
	h.ServeHTTP(hook, httptest.NewRequest("POST", "/hooks/github", nil))
	if hook.Code != http.StatusOK {
		t.Fatalf("profile route without CSRF token = %d, want 200", hook.Code)
	}

	form := httptest.NewRecorder()
	h.ServeHTTP(form, httptest.NewRequest("POST", "/form", nil))
	if form.Code != http.StatusOK {
		t.Fatalf("unprofiled route must not receive profile CSRF = %d, want 200", form.Code)
	}
}

func TestProfileAppliesSessionOnlyToItsRoutes(t *testing.T) {
	store := NewCookieStore(testSecret)
	app := New()
	if err := app.Profile("marketing", Profile{Session: store}); err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("/shop", "marketing"); err != nil {
		t.Fatal(err)
	}
	app.Register("GET", "/shop", func(w http.ResponseWriter, r *http.Request) {
		if session.StoreFromCtx(r.Context()) == nil {
			http.Error(w, "no store", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	app.Register("GET", "/plain", func(w http.ResponseWriter, r *http.Request) {
		if session.StoreFromCtx(r.Context()) != nil {
			http.Error(w, "unexpected store", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	h := app.Handler()
	if rec := httptest.NewRecorder(); func() int {
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/shop", nil))
		return rec.Code
	}() != http.StatusOK {
		t.Fatalf("profiled route should carry the profile store")
	}
	if rec := httptest.NewRecorder(); func() int {
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/plain", nil))
		return rec.Code
	}() != http.StatusOK {
		t.Fatalf("non-profiled route must not carry the profile store")
	}
}

func TestProfileCookiePolicyOverridesNameAndPath(t *testing.T) {
	store := NewCookieStore(testSecret)
	app := New()
	err := app.Profile("app", Profile{
		Session: store,
		Cookie:  &ProfileCookie{Name: "app_session", Path: "/app", SameSite: http.SameSiteStrictMode},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("/app", "app"); err != nil {
		t.Fatal(err)
	}
	app.Register("GET", "/app", func(w http.ResponseWriter, r *http.Request) {
		c := NewSSR(w, r)
		c.SetSessionVal("k", "v")
	})

	h := app.Handler()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/app", nil))

	cookie := findCookie(t, rec.Result().Cookies(), "app_session")
	if cookie.Path != "/app" {
		t.Errorf("cookie path = %q, want /app", cookie.Path)
	}
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Errorf("cookie SameSite = %v, want Strict", cookie.SameSite)
	}
}

func TestProfileNearestAncestorWins(t *testing.T) {
	outer := NewCookieStore(testSecret)
	inner := NewCookieStore(testSecret)
	app := New()
	if err := app.Profile("outer", Profile{Session: outer}); err != nil {
		t.Fatal(err)
	}
	if err := app.Profile("inner", Profile{Session: inner}); err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("/admin", "outer"); err != nil {
		t.Fatal(err)
	}
	if err := app.ApplyProfile("/admin/secret", "inner"); err != nil {
		t.Fatal(err)
	}

	app.Register("GET", "/admin/secret/x", func(w http.ResponseWriter, r *http.Request) {
		if session.StoreFromCtx(r.Context()) != Store(inner) {
			http.Error(w, "wrong store", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest("GET", "/admin/secret/x", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("nearest ancestor profile not applied: %d %s", rec.Code, rec.Body.String())
	}
}

func TestProfileCSRF403UsesErrorHandler(t *testing.T) {
	store := NewCookieStore(testSecret)
	app := New()
	if err := app.SetSessionStore(store); err != nil {
		t.Fatal(err)
	}
	if err := app.SetErrorHandler(http.StatusForbidden, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("nope-page"))
	}); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("POST", "/submit", okHandler); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest("POST", "/submit", nil))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
	if rec.Body.String() != "nope-page" {
		t.Fatalf("body = %q, want custom error handler output", rec.Body.String())
	}
}

func TestProfileCSRF403DefaultPlaintext(t *testing.T) {
	store := NewCookieStore(testSecret)
	app := New()
	if err := app.SetSessionStore(store); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("POST", "/submit", okHandler); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest("POST", "/submit", nil))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid csrf token") {
		t.Fatalf("body = %q, want plaintext default", rec.Body.String())
	}
}
