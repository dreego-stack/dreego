package core

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testSecret = []byte("test-secret-key-32-bytes-long!!!")
var testSecret2 = []byte("rotated-secret-key-32-bytes-!!!!")

func TestNewCookieStoreRejectsEmptySecret(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewCookieStore with empty secret should panic")
		}
	}()
	NewCookieStore([]byte{})
}

func TestNewCookieStoreRejectsShortSecret(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewCookieStore with short secret should panic")
		}
	}()
	NewCookieStore([]byte("short"))
}

func TestNewCookieStoreAcceptsValidSecret(t *testing.T) {
	store := NewCookieStore(testSecret)
	if store == nil {
		t.Fatal("expected non-nil store with valid secret")
	}
}

func TestSessionCookieSameSiteDefaultLax(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "k", "v", &Options{HttpOnly: true, Path: "/"})

	c := findCookie(t, w.Result().Cookies(), "dreego_session")
	if c.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax by default, got %v", c.SameSite)
	}
}

func TestSessionCookieSecureFromTLS(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.TLS = &tls.ConnectionState{}

	store.Set(w, r, "k", "v", &Options{HttpOnly: true, Path: "/"})

	c := findCookie(t, w.Result().Cookies(), "dreego_session")
	if !c.Secure {
		t.Error("expected Secure=true when request is over TLS")
	}
}

func TestSessionCookieSecureFromTrustedProxy(t *testing.T) {
	cases := []struct {
		name     string
		proxies  []string
		remote   string
		fwdProto string
		wantSec  bool
	}{
		{"trusted", []string{"10.0.0.1"}, "10.0.0.1:12345", "https", true},
		{"untrusted", []string{"10.0.0.1"}, "192.168.1.1:12345", "https", false},
		{"none-configured", nil, "10.0.0.1:12345", "https", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := NewCookieStore(testSecret)
			store.SetTrustedProxies(tc.proxies)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tc.remote
			r.Header.Set("X-Forwarded-Proto", tc.fwdProto)
			store.Set(w, r, "k", "v", &Options{HttpOnly: true, Path: "/"})
			c := findCookie(t, w.Result().Cookies(), "dreego_session")
			if c.Secure != tc.wantSec {
				t.Errorf("Secure = %v, want %v", c.Secure, tc.wantSec)
			}
		})
	}
}

func TestStoreGetReturnsErrorOnCorruptCookie(t *testing.T) {
	store := NewCookieStore(testSecret)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "dreego_session", Value: "corrupt-but-not-empty"})

	_, err := store.Get(r, "k")
	if err == nil {
		t.Error("expected error on corrupt cookie, got nil")
	}
}

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
	r = WithStore(r, failingStore{})
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
	r := WithStore(httptest.NewRequest("GET", "/", nil), failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	val := c.SessionVal("k")
	if val != "" {
		t.Errorf("expected empty value on store error, got %q", val)
	}
}

func TestCSRFCookieUsesSecurePolicyWithTLS(t *testing.T) {
	store := NewCookieStore(testSecret)
	mw := CSRF(store)
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
			store := NewCookieStore(testSecret)
			store.SetTrustedProxies(tc.proxies)
			mw := CSRF(store)
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

func TestOversizedSessionStateReturnsError(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	big := strings.Repeat("x", 5000)
	err := store.Set(w, r, "big", big, nil)
	if err == nil {
		t.Error("expected error for oversized session state, got nil")
	}
	if !errors.Is(err, ErrSessionTooLarge) {
		t.Errorf("expected ErrSessionTooLarge, got %v", err)
	}
}

func TestDestroyPreservesCookiePolicy(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.TLS = &tls.ConnectionState{}

	store.Set(w, r, "k", "v", &Options{HttpOnly: true, Secure: true, Path: "/", Encrypt: true})

	destroyW := httptest.NewRecorder()
	destroyReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		destroyReq.AddCookie(c)
	}
	destroyReq.TLS = &tls.ConnectionState{}

	store.Destroy(destroyW, destroyReq)

	c := findCookie(t, destroyW.Result().Cookies(), "dreego_session")
	if c.MaxAge != -1 {
		t.Errorf("destroy MaxAge = %d, want -1", c.MaxAge)
	}
	if !c.Secure {
		t.Error("destroy should preserve Secure")
	}
	if !c.HttpOnly {
		t.Error("destroy should preserve HttpOnly")
	}
	if c.SameSite == http.SameSiteDefaultMode {
		t.Error("destroy should preserve SameSite")
	}
}

func TestDeletePreservesCookiePolicy(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.TLS = &tls.ConnectionState{}

	store.Set(w, r, "k", "v", &Options{HttpOnly: true, Secure: true, Path: "/"})

	delW := httptest.NewRecorder()
	delReq := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		delReq.AddCookie(c)
	}
	delReq.TLS = &tls.ConnectionState{}

	store.Delete(delW, delReq, "k")

	c := findCookie(t, delW.Result().Cookies(), "dreego_session")
	if !c.Secure {
		t.Error("delete should preserve Secure")
	}
	if !c.HttpOnly {
		t.Error("delete should preserve HttpOnly")
	}
	if c.SameSite == http.SameSiteDefaultMode {
		t.Error("delete should preserve SameSite")
	}
}
