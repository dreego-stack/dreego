package session

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestCookieStoreRejectsPerCallPathOverride(t *testing.T) {
	store := NewCookieStore(testSecret)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)

	err := store.Set(rec, req, "user", "1", &Options{Path: "/admin"})
	if !errors.Is(err, ErrCookiePathOverride) {
		t.Fatalf("Set error = %v, want ErrCookiePathOverride", err)
	}
}

func TestCookieStoreDestroyUsesConfiguredPath(t *testing.T) {
	store := NewCookieStore(testSecret)
	store.SetCookiePolicy(CookiePolicy{Path: "/admin"})
	req := httptest.NewRequest("GET", "/admin", nil)

	setRec := httptest.NewRecorder()
	if err := store.Set(setRec, req, "user", "1", nil); err != nil {
		t.Fatal(err)
	}
	destroyRec := httptest.NewRecorder()
	if err := store.Destroy(destroyRec, req); err != nil {
		t.Fatal(err)
	}

	setPath := setRec.Result().Cookies()[0].Path
	destroyPath := destroyRec.Result().Cookies()[0].Path
	if setPath != "/admin" || destroyPath != setPath {
		t.Fatalf("cookie paths differ: set=%q destroy=%q", setPath, destroyPath)
	}
}
