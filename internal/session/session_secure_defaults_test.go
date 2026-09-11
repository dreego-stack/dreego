package session

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
