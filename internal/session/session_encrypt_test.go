package session

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http/httptest"
	"testing"
)

type shortReader struct{}

func (shortReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestCookieStoreEncryptDefaultReadable(t *testing.T) {
	store := NewCookieStore(testSecret)
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

func TestCookieStoreEncryptPropagatesError(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.rand = &shortReader{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	err := store.Set(w, r, "user_id", "42", &Options{Encrypt: true})
	if err == nil {
		t.Error("Set should propagate encryption errors, got nil")
	}

	cookies := w.Result().Cookies()
	for _, c := range cookies {
		if c.Name == "dreego_session" && c.Value != "" {
			t.Error("Set should not write a cookie when encryption fails")
		}
	}
}

func TestCookieStoreEncryptValueNotPlaintext(t *testing.T) {
	store := NewCookieStore(testSecret)
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
	// Value layout: sig (32B) || marker (1B) || base64(nonce||ciphertext).
	// Assert the marker and that the encrypted payload is not the plaintext
	// map. A byte-level substring check for the key/value is flaky: the random
	// ciphertext can coincidentally contain "user_id"/"42" as base64 bytes.
	if len(decoded) < sha256.Size+1 || decoded[sha256.Size] != encMarker {
		t.Fatalf("expected encrypted payload (marker %d), got marker %d", encMarker, decoded[sha256.Size])
	}
	payload, err := base64.RawURLEncoding.DecodeString(string(decoded[sha256.Size+1:]))
	if err != nil {
		t.Fatalf("decode payload failed: %v", err)
	}
	if bytes.Contains(payload, []byte("user_id")) || bytes.Contains(payload, []byte("42")) {
		t.Error("encrypted cookie payload should not contain plaintext key/value")
	}
}

func TestCookieStoreEncryptRoundTrip(t *testing.T) {
	store := NewCookieStore(testSecret)
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
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "role", "admin", &Options{Encrypt: true})

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		c.Value = base64.RawURLEncoding.EncodeToString([]byte("tampered"))
		req.AddCookie(c)
	}

	_, err := store.Get(req, "role")
	if err == nil {
		t.Error("tampered encrypted cookie should return error, got nil")
	}
}

func TestCookieStoreEncryptKeyRotationRejected(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "role", "admin", &Options{Encrypt: true})

	other := NewCookieStore([]byte("rotated-secret-key-32-bytes-!!!!"))
	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	_, err := other.Get(req, "role")
	if err == nil {
		t.Error("cookie encrypted with different key should be rejected with error, got nil")
	}
}
