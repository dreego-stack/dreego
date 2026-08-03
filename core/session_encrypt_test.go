package core

import (
	"encoding/base64"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCookieStoreEncryptDefaultReadable(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "user_id", "42", nil)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	val, _ := store.Get(req, "user_id")
	if val != "42" {
		t.Errorf("expected '42', got '%s'", val)
	}
}

func TestCookieStoreEncryptValueNotPlaintext(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "user_id", "42", &Options{Encrypt: true})

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	decoded, err := base64.RawURLEncoding.DecodeString(cookies[0].Value)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if strings.Contains(string(decoded), "user_id") || strings.Contains(string(decoded), "42") {
		t.Error("encrypted cookie should not contain plaintext key/value")
	}
}

func TestCookieStoreEncryptRoundTrip(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "user_id", "42", &Options{Encrypt: true})

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	val, _ := store.Get(req, "user_id")
	if val != "42" {
		t.Errorf("expected '42', got '%s'", val)
	}
}

func TestCookieStoreEncryptTamperRejected(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "role", "admin", &Options{Encrypt: true})

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		c.Value = base64.RawURLEncoding.EncodeToString([]byte("tampered"))
		req.AddCookie(c)
	}

	val, _ := store.Get(req, "role")
	if val != "" {
		t.Errorf("tampered encrypted cookie should return empty, got '%s'", val)
	}
}

func TestCookieStoreEncryptKeyRotationRejected(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "role", "admin", &Options{Encrypt: true})

	other := NewCookieStore([]byte("other-secret"))
	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	val, _ := other.Get(req, "role")
	if val != "" {
		t.Errorf("cookie encrypted with different key should be rejected, got '%s'", val)
	}
}
