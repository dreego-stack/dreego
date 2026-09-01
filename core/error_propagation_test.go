package core

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/middleware"
	"github.com/dreego-stack/dreego/core/internal/session"
)

func TestSSRContextSessionValErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = session.WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.SessionVal("key"); got != "" {
		t.Errorf("expected empty value on store error, got %q", got)
	}
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store read failure")
	}
	if !strings.Contains(c.SessionError().Error(), "store read failure") {
		t.Errorf("SessionError must contain the cause, got %v", c.SessionError())
	}
}

func TestSSRContextSetSessionValErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = session.WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	c.SetSessionVal("key", "value")
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store write failure")
	}
}

func TestSSRContextDelSessionValErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = session.WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	c.DelSessionVal("key")
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store delete failure")
	}
}

func TestSSRContextDestroySessionErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = session.WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	c.DestroySession()
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store destroy failure")
	}
}

func TestSSRContextCSRFTokenErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = session.WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.CSRFToken(); got != "" {
		t.Errorf("expected empty token on store error, got %q", got)
	}
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface CSRF read failure")
	}
}

func TestSSRContextFormErrorReachable(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("a=%zz"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.FormValue("a"); got != "" {
		t.Errorf("expected empty value on parse error, got %q", got)
	}
	if c.FormError() == nil {
		t.Fatal("expected FormError to surface parse failure")
	}
}

func TestSSRContextFormErrorDistinguishableFromEmpty(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("a="))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.FormValue("a"); got != "" {
		t.Errorf("expected empty value for empty field, got %q", got)
	}
	if c.FormError() != nil {
		t.Errorf("valid empty form must not produce FormError, got %v", c.FormError())
	}
}

func TestRecoveryDoesNotDiscloseCause(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/boom", nil)

	middleware.Recovery(nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(fmt.Errorf("database: connection to db.example.com:5432 failed"))
	})).ServeHTTP(w, r)

	body := w.Body.String()
	if strings.Contains(body, "database") || strings.Contains(body, "db.example.com") {
		t.Errorf("panic cause disclosed in body: %q", body)
	}
}

func TestCSRFStoreFailureReachesErrorPath(t *testing.T) {
	mw := middleware.CSRF(failingStore{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	hit := false

	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if hit {
		t.Fatal("handler ran after CSRF token persistence failed")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if strings.Contains(w.Body.String(), "store read failure") {
		t.Errorf("internal cause disclosed in body: %q", w.Body.String())
	}
}

func TestListenErrorPropagates(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	app := New()
	err = app.Listen(ln.Addr().String())
	if err == nil {
		t.Fatal("expected Listen to return an error when the port is occupied")
	}
}
