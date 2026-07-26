package session

import (
	"net/http/httptest"
	"testing"
)

func TestCookieStoreSetAndGet(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
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
	store := NewCookieStore([]byte("secret-key"))
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
	store := NewCookieStore([]byte("secret-key"))
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

func TestCookieStoreTamperDetection(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "role", "admin", nil)

	req := httptest.NewRequest("GET", "/", nil)
	for _, c := range w.Result().Cookies() {
		c.Value = "tampered-data"
		req.AddCookie(c)
	}

	val, err := store.Get(req, "role")
	if err != nil {
		t.Fatalf("tampered data should return empty, not error: %v", err)
	}
	if val != "" {
		t.Errorf("tampered cookie should return empty, got '%s'", val)
	}
}

func TestCookieStoreMultipleKeys(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	store.Set(w, r, "a", "1", nil)
	store.Set(w, r, "b", "2", nil)

	cookies := w.Result().Cookies()
	if len(cookies) < 1 {
		t.Fatal("no cookies set")
	}
	last := cookies[len(cookies)-1]

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(last)

	a, _ := store.Get(req, "a")
	b, _ := store.Get(req, "b")
	if a != "1" || b != "2" {
		t.Errorf("expected a=1 b=2, got a=%s b=%s", a, b)
	}
}

func TestCookieStoreOptions(t *testing.T) {
	store := NewCookieStore([]byte("secret-key"))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	opts := &Options{
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		Path:     "/app",
	}
	store.Set(w, r, "key", "value", opts)

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
