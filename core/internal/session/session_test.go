package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCookieStoreSetAndGet(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	err := store.Set(w, r, "user_id", "42", nil)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	val, err := store.Get(req, "user_id")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "42" {
		t.Errorf("expected '42', got '%s'", val)
	}
}

func TestCookieStoreDelete(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "key", "value", nil)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	w2 := httptest.NewRecorder()
	store.Delete(w2, req, "key")

	req2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w2.Result().Cookies() {
		req2.AddCookie(c)
	}

	val, err := store.Get(req2, "key")
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty after delete, got '%s'", val)
	}
}

func TestCookieStoreDestroy(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "a", "1", nil)
	store.Set(w, r, "b", "2", nil)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	w2 := httptest.NewRecorder()
	store.Destroy(w2, req)

	req2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w2.Result().Cookies() {
		req2.AddCookie(c)
	}

	a, _ := store.Get(req2, "a")
	b, _ := store.Get(req2, "b")
	if a != "" || b != "" {
		t.Errorf("expected all keys empty after destroy, got a=%s b=%s", a, b)
	}
}

func TestCookieStoreOptions(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.SetCookiePolicy(CookiePolicy{Path: "/app"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	opts := &Options{
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		Path:     "/app",
	}
	if err := store.Set(w, r, "key", "value", opts); err != nil {
		t.Fatal(err)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Name != "dreego_session" {
		t.Errorf("expected cookie name dreego_session, got %s", c.Name)
	}
	if c.MaxAge != 3600 {
		t.Errorf("expected MaxAge 3600, got %d", c.MaxAge)
	}
	if !c.Secure {
		t.Error("expected Secure true")
	}
	if !c.HttpOnly {
		t.Error("expected HttpOnly true")
	}
	if c.Path != "/app" {
		t.Errorf("expected Path /app, got %s", c.Path)
	}
}

func TestCookieStoreVerifyInvalidBase64(t *testing.T) {
	store := NewCookieStore(testSecret)
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "dreego_session", Value: "!!!not-base64!!!"})

	_, err := store.Get(r, "key")
	if err == nil {
		t.Error("expected error for invalid base64 cookie, got nil")
	}
}

func TestCookieStoreVerifyTooShort(t *testing.T) {
	store := NewCookieStore(testSecret)
	short := base64.RawURLEncoding.EncodeToString([]byte("short"))
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "dreego_session", Value: short})

	_, err := store.Get(r, "key")
	if err == nil {
		t.Error("expected error for too-short cookie, got nil")
	}
}

func TestCookieStoreCorruptJSON(t *testing.T) {
	store := NewCookieStore(testSecret)
	payload := []byte("{not valid json")
	mac := hmac.New(sha256.New, deriveKeys(testSecret).sig)
	mac.Write(payload)
	sig := mac.Sum(nil)
	value := base64.RawURLEncoding.EncodeToString(append(sig, payload...))

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "dreego_session", Value: value})

	_, err := store.Get(r, "key")
	if err == nil {
		t.Error("expected error for corrupt JSON, got nil")
	}
}

func TestCookieStoreEmptySecretPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty secret")
		}
	}()
	NewCookieStore([]byte{})
}

func TestCookieStoreSetEmptyValueDeletes(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	if err := store.Set(w, r, "key", "value", nil); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		req.AddCookie(c)
	}

	w2 := httptest.NewRecorder()
	if err := store.Set(w2, req, "key", "", nil); err != nil {
		t.Fatalf("Set empty failed: %v", err)
	}

	req2 := httptest.NewRequest("GET", "/", nil)
	for _, c := range w2.Result().Cookies() {
		req2.AddCookie(c)
	}

	val, err := store.Get(req2, "key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty value after Set empty, got %q", val)
	}
}
