package core

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type failingCSRFStore struct{}

func (failingCSRFStore) Get(*http.Request, string) (string, error) {
	return "", errors.New("read failed")
}

func (failingCSRFStore) Set(http.ResponseWriter, *http.Request, string, string, *Options) error {
	return errors.New("write failed")
}

func (failingCSRFStore) Delete(http.ResponseWriter, *http.Request, string) error { return nil }
func (failingCSRFStore) Destroy(http.ResponseWriter, *http.Request) error        { return nil }

func TestCSRFFailsClosedWhenTokenCannotBePersisted(t *testing.T) {
	called := false
	h := CSRF(failingCSRFStore{})(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))

	if called {
		t.Fatal("handler ran without a persisted CSRF token")
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
