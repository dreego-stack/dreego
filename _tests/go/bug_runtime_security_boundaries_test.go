package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

type unavailableStore struct{}

func (unavailableStore) Get(*http.Request, string) (string, error) {
	return "", errors.New("unavailable")
}
func (unavailableStore) Set(http.ResponseWriter, *http.Request, string, string, *dreego.Options) error {
	return errors.New("unavailable")
}
func (unavailableStore) Delete(http.ResponseWriter, *http.Request, string) error { return nil }
func (unavailableStore) Destroy(http.ResponseWriter, *http.Request) error        { return nil }

func TestBugCSRFFailsClosedWithoutPersistedToken(t *testing.T) {
	called := false
	h := dreego.CSRF(unavailableStore{})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if called || rec.Code != http.StatusInternalServerError {
		t.Fatalf("called=%v status=%d, want false and 500", called, rec.Code)
	}
}

func TestBugWildcardRedirectCannotBecomeProtocolRelative(t *testing.T) {
	app := dreego.New()
	if err := app.RegisterRedirect("/old/*", "/*", http.StatusFound); err == nil {
		t.Fatal("unsafe wildcard redirect was accepted")
	}
}

func TestBugSessionRejectsInconsistentCookiePath(t *testing.T) {
	store := dreego.NewCookieStore([]byte("01234567890123456789012345678901"))
	err := store.Set(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/admin", nil), "user", "1", &dreego.Options{Path: "/admin"})
	if !errors.Is(err, dreego.ErrCookiePathOverride) {
		t.Fatalf("error = %v, want ErrCookiePathOverride", err)
	}
}
