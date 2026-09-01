package core

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/middleware"
	"github.com/dreego-stack/dreego/core/internal/session"
)

var testSecret = []byte("test-secret-key-32-bytes-long!!!")

var testSecret2 = []byte("rotated-secret-key-32-bytes-!!!!")

type failingStore struct{}

func (failingStore) Get(*http.Request, string) (string, error) {
	return "", errors.New("store read failure")
}
func (failingStore) Set(http.ResponseWriter, *http.Request, string, string, *Options) error {
	return errors.New("store write failure")
}
func (failingStore) Delete(http.ResponseWriter, *http.Request, string) error {
	return errors.New("store delete failure")
}
func (failingStore) Destroy(http.ResponseWriter, *http.Request) error {
	return errors.New("store destroy failure")
}

func TestSetSessionValSurfacesStoreError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = session.WithStore(r, failingStore{})
	c := NewSSR(w, r)

	c.SetSessionVal("k", "v")

	body := w.Body.String()
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 on store failure, got %d", w.Code)
	}
	if !strings.Contains(body, "internal server error") {
		t.Errorf("expected generic error body, got %q", body)
	}
	if strings.Contains(body, "store write failure") {
		t.Errorf("internal error must not leak to client, got %q", body)
	}
}

func TestSessionValReturnsEmptyOnStoreGetError(t *testing.T) {
	r := session.WithStore(httptest.NewRequest("GET", "/", nil), failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	val := c.SessionVal("k")
	if val != "" {
		t.Errorf("expected empty value on store error, got %q", val)
	}
}

func TestCSRFCookieUsesSecurePolicyWithTLS(t *testing.T) {
	store := NewCookieStore(testSecret)
	mw := middleware.CSRF(store)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.TLS = &tls.ConnectionState{}
	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, r)

	readable := csrfReadablePolicyCookie(t, w.Result().Cookies())
	if readable == nil {
		t.Fatal("expected readable csrf_token cookie")
	}
	if !readable.Secure {
		t.Error("CSRF cookie should be Secure=true under TLS")
	}
	if readable.SameSite != http.SameSiteStrictMode {
		t.Errorf("CSRF cookie SameSite = %v, want Strict", readable.SameSite)
	}
}

func TestCSRFCookieSecureBehindProxies(t *testing.T) {
	cases := []struct {
		name    string
		proxies []string
		remote  string
		wantSec bool
	}{
		{"trusted", []string{"10.0.0.1"}, "10.0.0.1:12345", true},
		{"untrusted", []string{"10.0.0.1"}, "192.168.99.99:12345", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := NewCookieStore(testSecret2)
			store.SetTrustedProxies(tc.proxies)
			mw := middleware.CSRF(store)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tc.remote
			r.Header.Set("X-Forwarded-Proto", "https")
			mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, r)
			c := csrfReadablePolicyCookie(t, w.Result().Cookies())
			if c == nil {
				t.Fatal("expected readable csrf_token cookie")
			}
			if c.Secure != tc.wantSec {
				t.Errorf("CSRF Secure = %v, want %v", c.Secure, tc.wantSec)
			}
		})
	}
}
