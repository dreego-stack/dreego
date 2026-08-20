package core

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppConfigurableCookiePolicy(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.SetCookiePolicy(CookiePolicy{
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/app",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "k", "v", nil)

	c := findCookie(t, w.Result().Cookies(), "dreego_session")
	if c.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSite=Strict from policy, got %v", c.SameSite)
	}
	if !c.Secure {
		t.Error("expected Secure=true from policy")
	}
	if !c.HttpOnly {
		t.Error("expected HttpOnly=true from policy")
	}
	if c.Path != "/app" {
		t.Errorf("expected Path=/app from policy, got %q", c.Path)
	}
}

func TestSetWithOptionsPreservesSameSitePolicy(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.SetCookiePolicy(CookiePolicy{
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Path:     "/",
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "k", "v", &Options{HttpOnly: true, Path: "/"})

	c := findCookie(t, w.Result().Cookies(), "dreego_session")
	if c.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSite=Strict from policy with per-call options, got %v", c.SameSite)
	}
}

func TestSetSessionValPreservesCookiePolicy(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.SetCookiePolicy(CookiePolicy{
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/app",
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, store)
	c := NewSSR(w, r)

	c.SetSessionVal("k", "v")

	cookie := findCookie(t, w.Result().Cookies(), "dreego_session")
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSite=Strict from policy via SetSessionVal, got %v", cookie.SameSite)
	}
	if cookie.Path != "/app" {
		t.Errorf("expected Path=/app from policy via SetSessionVal, got %q", cookie.Path)
	}
	if !cookie.Secure {
		t.Error("expected Secure=true from policy via SetSessionVal")
	}
	if !cookie.HttpOnly {
		t.Error("expected HttpOnly=true from policy via SetSessionVal")
	}
}

func TestCSRFSessionWritePreservesCookiePolicy(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.SetCookiePolicy(CookiePolicy{
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Path:     "/",
	})
	mw := CSRF(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, r)

	sess := findCookie(t, w.Result().Cookies(), "dreego_session")
	if sess.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSite=Strict from policy on CSRF session write, got %v", sess.SameSite)
	}
}

func TestAppValidatesSessionStoreSecretAtBuild(t *testing.T) {
	app := New()
	if err := app.SetSessionStore(storeWithFailedValidation{}); err != nil {
		t.Fatalf("SetSessionStore failed: %v", err)
	}
	if err := app.Build(); err == nil {
		t.Fatal("expected error when building app with invalid store")
	}
}

func TestNewCookieStoreShortSecretPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for short secret")
		}
	}()
	NewCookieStore([]byte("x"))
}

func TestSetCookiePolicyPartialKeepsSecureDefaults(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.SetCookiePolicy(CookiePolicy{SameSite: http.SameSiteStrictMode})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "k", "v", nil)

	c := findCookie(t, w.Result().Cookies(), "dreego_session")
	if c.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSite=Strict from partial policy, got %v", c.SameSite)
	}
	if !c.HttpOnly {
		t.Error("partial policy must keep default HttpOnly=true")
	}
	if c.Path != "/" {
		t.Errorf("partial policy must keep default Path=/, got %q", c.Path)
	}
}

type storeWithFailedValidation struct{ customSecureStore }

func (storeWithFailedValidation) Validate() error { return errors.New("invalid store config") }

func TestReplaceableStoreDoesNotDisableSecurity(t *testing.T) {
	app := New()
	custom := &customSecureStore{}
	if err := app.SetSessionStore(custom); err != nil {
		t.Fatalf("SetSessionStore failed: %v", err)
	}
	if app.SessionStore() == nil {
		t.Error("expected non-nil store after SetSessionStore")
	}
}

type customSecureStore struct{}

func (customSecureStore) Get(*http.Request, string) (string, error) { return "", nil }
func (customSecureStore) Set(http.ResponseWriter, *http.Request, string, string, *Options) error {
	return nil
}
func (customSecureStore) Delete(http.ResponseWriter, *http.Request, string) error { return nil }
func (customSecureStore) Destroy(http.ResponseWriter, *http.Request) error        { return nil }

func findCookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	t.Fatalf("expected cookie %s, got %d cookies", name, len(cookies))
	return nil
}

func csrfReadablePolicyCookie(t *testing.T, cookies []*http.Cookie) *http.Cookie {
	t.Helper()
	for _, c := range cookies {
		if c.Name == "csrf_token" && !c.HttpOnly {
			return c
		}
	}
	return nil
}
