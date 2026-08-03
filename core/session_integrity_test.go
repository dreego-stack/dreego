package core

import (
	"net/http/httptest"
	"testing"
)

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
