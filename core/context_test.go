package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/core/internal/middleware"
)

func TestNewSSRNilRequest(t *testing.T) {
	c := NewSSR(nil, nil)
	if c.Context == nil {
		t.Fatal("expected non-nil background context")
	}
	c.Set("k", "v")
	if got := c.Data("k"); got != "v" {
		t.Errorf("Data: expected 'v', got %v", got)
	}
	if got := c.Get("k"); got != "v" {
		t.Errorf("Get: expected 'v', got %q", got)
	}
}

func TestNewSSRWithRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(r.Context(), "probe", "yes")
	r = r.WithContext(ctx)
	c := NewSSR(nil, r)
	if c.Context.Value("probe") != "yes" {
		t.Error("expected context to be derived from r.Context()")
	}
}

func TestSSRContextDataNilMap(t *testing.T) {
	c := &SSRContext{Context: context.Background()}
	if c.Data("missing") != nil {
		t.Error("expected Data on nil map to return nil")
	}
}

func TestSSRContextGetNilMap(t *testing.T) {
	c := &SSRContext{Context: context.Background()}
	if c.Get("missing") != "" {
		t.Error("expected Get on nil map to return empty string")
	}
}

func TestSSRContextSetNilMap(t *testing.T) {
	c := &SSRContext{Context: context.Background()}
	c.Set("k", "v")
	if got := c.Get("k"); got != "v" {
		t.Errorf("expected Get after Set on nil map to return 'v', got %q", got)
	}
}

func TestSSRContextParam(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.SetPathValue("id", "42")
	c := NewSSR(nil, r)
	if got := c.Param("id"); got != "42" {
		t.Errorf("expected param '42', got %q", got)
	}
	if got := c.Param("missing"); got != "" {
		t.Errorf("expected missing param '', got %q", got)
	}
}

func TestSSRContextQuery(t *testing.T) {
	r := httptest.NewRequest("GET", "/?q=hello&empty=", nil)
	c := NewSSR(nil, r)
	if got := c.Query("q"); got != "hello" {
		t.Errorf("expected query 'hello', got %q", got)
	}
	if got := c.Query("missing"); got != "" {
		t.Errorf("expected missing query '', got %q", got)
	}
}

func TestSSRContextFormValue(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("name=anna&age=30"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := NewSSR(nil, r)
	if got := c.FormValue("name"); got != "anna" {
		t.Errorf("expected form value 'anna', got %q", got)
	}
	if got := c.FormValue("missing"); got != "" {
		t.Errorf("expected missing form value '', got %q", got)
	}
}

func TestSSRContextFormValueParseError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("a=%zz"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := NewSSR(nil, r)
	if got := c.FormValue("a"); got != "" {
		t.Errorf("expected empty on parse error, got %q", got)
	}
}

func TestSSRContextSessionValNilStore(t *testing.T) {
	c := NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if got := c.SessionVal("key"); got != "" {
		t.Errorf("expected empty with nil store, got %q", got)
	}
}

func TestSSRContextSetSessionValNilStore(t *testing.T) {
	w := httptest.NewRecorder()
	c := NewSSR(w, httptest.NewRequest("GET", "/", nil))
	c.SetSessionVal("key", "value")
	if len(w.Result().Cookies()) != 0 {
		t.Error("expected no cookies set with nil store")
	}
}

func TestSSRContextDelSessionValNilStore(t *testing.T) {
	c := NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	c.DelSessionVal("key")
}

func TestSSRContextDestroySessionNilStore(t *testing.T) {
	c := NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	c.DestroySession()
}

func TestSSRContextSessionVal(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, store)
	c := NewSSR(w, r)

	c.SetSessionVal("user_id", "42")

	req := httptest.NewRequest("GET", "/", nil)
	for _, ck := range w.Result().Cookies() {
		req.AddCookie(ck)
	}
	req = WithStore(req, store)
	c2 := NewSSR(httptest.NewRecorder(), req)
	if got := c2.SessionVal("user_id"); got != "42" {
		t.Errorf("expected session value '42', got %q", got)
	}
}

func TestSSRContextCSRFToken(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, store)
	c := NewSSR(w, r)
	c.SetSessionVal("csrf_token", "tok123")

	req := httptest.NewRequest("GET", "/", nil)
	for _, ck := range w.Result().Cookies() {
		req.AddCookie(ck)
	}
	req = WithStore(req, store)
	c2 := NewSSR(httptest.NewRecorder(), req)
	if got := c2.CSRFToken(); got != "tok123" {
		t.Errorf("expected csrf token 'tok123', got %q", got)
	}
}

func TestSSRContextRequestIDMissing(t *testing.T) {
	c := NewSSR(nil, httptest.NewRequest("GET", "/", nil))
	if got := c.RequestID(); got != "" {
		t.Errorf("expected empty request id, got %q", got)
	}
}

func TestSSRContextRequestIDPresent(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), middleware.RequestIDKey, "abc123"))
	c := NewSSR(nil, r)
	if got := c.RequestID(); got != "abc123" {
		t.Errorf("expected request id 'abc123', got %q", got)
	}
}

func TestSSRContextErrors(t *testing.T) {
	c := NewSSR(nil, httptest.NewRequest("GET", "/", nil))
	c.Set("error_email", "invalid")
	if got := c.Errors("email"); got != "invalid" {
		t.Errorf("expected error for email, got %q", got)
	}
	if got := c.Errors("missing"); got != "" {
		t.Errorf("expected empty for missing field, got %q", got)
	}
}

func TestSSRContextOld(t *testing.T) {
	c := NewSSR(nil, httptest.NewRequest("GET", "/", nil))
	c.Set("old_name", "Lukas")
	if got := c.Old("name"); got != "Lukas" {
		t.Errorf("expected old value 'Lukas', got %q", got)
	}
	if got := c.Old("missing"); got != "" {
		t.Errorf("expected empty for missing field, got %q", got)
	}
}

func TestSSRContextRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	c := NewSSR(w, r)
	err := c.Redirect("/login", http.StatusFound)
	if err != ErrRedirect {
		t.Errorf("expected ErrRedirect, got %v", err)
	}
	if w.Code != http.StatusFound {
		t.Errorf("expected status %d, got %d", http.StatusFound, w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/login" {
		t.Errorf("expected Location /login, got %q", loc)
	}
}
