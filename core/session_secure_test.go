package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionCookieSameSiteDefault(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "k", "v", &Options{
		HttpOnly: true,
		Path:     "/",
	})

	cookies := w.Result().Cookies()
	var sess *http.Cookie
	for _, c := range cookies {
		if c.Name == "dreego_session" {
			sess = c
		}
	}
	if sess == nil {
		t.Fatal("expected dreego_session cookie")
	}
	if sess.SameSite == http.SameSiteDefaultMode {
		t.Error("expected SameSite not DefaultMode (unspecified)")
	}
}

func TestSessionCookieSecurePassedThrough(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "k", "v", &Options{
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	cookies := w.Result().Cookies()
	var sess *http.Cookie
	for _, c := range cookies {
		if c.Name == "dreego_session" {
			sess = c
		}
	}
	if sess == nil {
		t.Fatal("expected dreego_session cookie")
	}
	if !sess.Secure {
		t.Error("expected Secure=true when Options.Secure true")
	}
}
