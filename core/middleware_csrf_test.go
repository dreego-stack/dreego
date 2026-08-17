package core

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFCookieSecureWithTLS(t *testing.T) {
	store := NewCookieStore(testSecret)
	mw := CSRF(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.TLS = &tls.ConnectionState{}
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	cookies := w.Result().Cookies()
	var csrfReadable *http.Cookie
	for _, c := range cookies {
		if c.Name == "csrf_token" && !c.HttpOnly {
			csrfReadable = c
		}
	}
	if csrfReadable == nil {
		t.Fatal("expected readable csrf_token cookie")
	}
	if !csrfReadable.Secure {
		t.Error("expected readable csrf cookie Secure=true when TLS")
	}
}

func TestCSRFCookieSecureWithoutTLS(t *testing.T) {
	store := NewCookieStore(testSecret)
	mw := CSRF(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	cookies := w.Result().Cookies()
	var csrfReadable *http.Cookie
	for _, c := range cookies {
		if c.Name == "csrf_token" && !c.HttpOnly {
			csrfReadable = c
		}
	}
	if csrfReadable == nil {
		t.Fatal("expected readable csrf_token cookie")
	}
	if csrfReadable.Secure {
		t.Error("expected readable csrf cookie Secure=false without TLS")
	}
}

func TestCSRFCookieSameSiteSet(t *testing.T) {
	store := NewCookieStore(testSecret)
	mw := CSRF(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	cookies := w.Result().Cookies()
	var csrfReadable *http.Cookie
	for _, c := range cookies {
		if c.Name == "csrf_token" && !c.HttpOnly {
			csrfReadable = c
		}
	}
	if csrfReadable == nil {
		t.Fatal("expected readable csrf_token cookie")
	}
	if csrfReadable.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected readable csrf cookie SameSite=Strict, got %v", csrfReadable.SameSite)
	}
}
