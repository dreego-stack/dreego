package core

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type csrfMockStore struct {
	data map[string]string
}

func newCSRFMockStore() *csrfMockStore {
	return &csrfMockStore{data: map[string]string{}}
}

func (m *csrfMockStore) Get(_ *http.Request, key string) (string, error) {
	return m.data[key], nil
}

func (m *csrfMockStore) Set(_ http.ResponseWriter, _ *http.Request, key, value string, _ *Options) error {
	if value == "" {
		delete(m.data, key)
	} else {
		m.data[key] = value
	}
	return nil
}

func (m *csrfMockStore) Delete(_ http.ResponseWriter, _ *http.Request, key string) error {
	delete(m.data, key)
	return nil
}

func (m *csrfMockStore) Destroy(http.ResponseWriter, *http.Request) error {
	m.data = map[string]string{}
	return nil
}

func csrfReadableCookie(w *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == "csrf_token" && !c.HttpOnly {
			return c
		}
	}
	return nil
}

func seedCSRFToken(t *testing.T, store *csrfMockStore) string {
	t.Helper()
	mw := CSRF(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, r)
	c := csrfReadableCookie(w)
	if c == nil || c.Value == "" {
		t.Fatal("expected seeded csrf_token cookie")
	}
	return c.Value
}

func TestCSRFGetSetsTokenCookieAndSession(t *testing.T) {
	store := newCSRFMockStore()
	mw := CSRF(store)
	hit := false
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if !hit {
		t.Fatal("GET handler not reached")
	}
	c := csrfReadableCookie(w)
	if c == nil {
		t.Fatal("expected readable csrf_token cookie")
	}
	if c.Value == "" {
		t.Error("expected non-empty csrf_token value")
	}
	if store.data["csrf_token"] == "" {
		t.Error("expected csrf_token stored in session")
	}
	if store.data["csrf_token"] != c.Value {
		t.Errorf("cookie %q != session %q", c.Value, store.data["csrf_token"])
	}
}

func TestCSRFGetWithoutTokenPasses(t *testing.T) {
	store := newCSRFMockStore()
	mw := CSRF(store)
	hit := false
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if w.Code == http.StatusForbidden {
		t.Error("GET must not be forbidden without token")
	}
	if !hit {
		t.Error("GET handler should be reached without token")
	}
}

func TestCSRFPostWithoutTokenForbidden(t *testing.T) {
	store := newCSRFMockStore()
	mw := CSRF(store)
	hit := false
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
	if hit {
		t.Error("handler should not be reached on missing token")
	}
}

func TestCSRFPostWithValidHeaderToken(t *testing.T) {
	store := newCSRFMockStore()
	token := seedCSRFToken(t, store)
	mw := CSRF(store)

	hit := false
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("X-CSRF-Token", token)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !hit {
		t.Error("handler should be reached with valid header token")
	}
}

func TestCSRFPostWithInvalidHeaderToken(t *testing.T) {
	store := newCSRFMockStore()
	seedCSRFToken(t, store)
	mw := CSRF(store)

	hit := false
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	r.Header.Set("X-CSRF-Token", "wrong-token")
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
	if hit {
		t.Error("handler should not be reached with invalid header token")
	}
}

func TestCSRFPostWithFormFieldToken(t *testing.T) {
	store := newCSRFMockStore()
	token := seedCSRFToken(t, store)
	mw := CSRF(store)

	form := url.Values{}
	form.Set("csrf_token", token)
	hit := false
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !hit {
		t.Error("handler should be reached via form-field fallback")
	}
}

func TestIsUnsafeMethod(t *testing.T) {
	for _, m := range []string{"POST", "PUT", "DELETE", "PATCH"} {
		if !isUnsafeMethod(m) {
			t.Errorf("expected %s to be unsafe", m)
		}
	}
	for _, m := range []string{"GET", "HEAD", "OPTIONS"} {
		if isUnsafeMethod(m) {
			t.Errorf("expected %s to be safe", m)
		}
	}
}

func TestGenerateCSRFToken(t *testing.T) {
	token := generateCSRFToken()
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if _, err := base64.RawURLEncoding.DecodeString(token); err != nil {
		t.Errorf("token not base64url-decodable: %v", err)
	}
	if token == generateCSRFToken() {
		t.Error("expected distinct tokens")
	}
}
