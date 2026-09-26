package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/session"
)

var _ Context = (*SSRContext)(nil)

func TestContextInterfaceExposesStateAndRequestMethods(t *testing.T) {
	var c Context = NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/?q=hi", nil))
	c.Set("name", "ada")
	if got := c.Get("name"); got != "ada" {
		t.Errorf("Get after Set = %q, want ada", got)
	}
	if got := c.Query("q"); got != "hi" {
		t.Errorf("Query(q) = %q, want hi", got)
	}
	c.Delete("name")
	if got := c.Data("name"); got != nil {
		t.Errorf("Data after Delete = %v, want nil", got)
	}
}

func TestContextInterfaceFormValueAndDestroySession(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("name=ada"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var c Context = NewSSR(httptest.NewRecorder(), r)
	if got := c.FormValue("name"); got != "ada" {
		t.Errorf("FormValue(name) = %q, want ada", got)
	}
	c.DestroySession()
	if err := c.SessionError(); err != nil {
		t.Errorf("DestroySession without store error = %v", err)
	}
}

func TestFlashSetGetConsumesValue(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r = session.WithStore(r, store)
	c := NewSSR(w, r)

	c.Flash("notice", "saved")
	if got := c.FlashPeek("notice"); got != "saved" {
		t.Fatalf("FlashPeek = %q, want saved", got)
	}

	next := flashRequest(t, w, store)
	c2 := NewSSR(httptest.NewRecorder(), next)
	if got := c2.FlashGet("notice"); got != "saved" {
		t.Fatalf("FlashGet = %q, want saved", got)
	}
	if got := c2.FlashPeek("notice"); got != "" {
		t.Fatalf("FlashPeek after consume = %q, want empty", got)
	}
}

func TestFlashPeekDoesNotConsume(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := session.WithStore(httptest.NewRequest("GET", "/", nil), store)
	c := NewSSR(w, r)
	c.Flash("error", "boom")

	if got := c.FlashPeek("error"); got != "boom" {
		t.Fatalf("FlashPeek = %q, want boom", got)
	}
	if got := c.FlashPeek("error"); got != "boom" {
		t.Fatalf("second FlashPeek = %q, want boom (no consume)", got)
	}
}

func TestFlashWithoutStoreIsNoOp(t *testing.T) {
	c := NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	c.Flash("notice", "saved")
	if got := c.FlashGet("notice"); got != "" {
		t.Fatalf("FlashGet without store = %q, want empty", got)
	}
	if c.SessionError() != nil {
		t.Fatalf("unexpected session error: %v", c.SessionError())
	}
}

func TestCSRFInputRendersHiddenField(t *testing.T) {
	store := NewCookieStore(testSecret)
	w := httptest.NewRecorder()
	r := session.WithStore(httptest.NewRequest("GET", "/", nil), store)
	c := NewSSR(w, r)
	c.SetSessionVal("csrf_token", "tok<1>&2")

	got := c.CSRFInput()
	want := `<input type="hidden" name="csrf_token" value="tok&lt;1&gt;&amp;2">`
	if got != want {
		t.Fatalf("CSRFInput() = %q, want %q", got, want)
	}
}

func TestCSRFInputEmptyToken(t *testing.T) {
	c := NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if got := c.CSRFInput(); got != `<input type="hidden" name="csrf_token" value="">` {
		t.Fatalf("CSRFInput() = %q", got)
	}
}

func flashRequest(t *testing.T, w *httptest.ResponseRecorder, store Store) *http.Request {
	t.Helper()
	req := httptest.NewRequest("GET", "/", nil)
	for _, ck := range w.Result().Cookies() {
		req.AddCookie(ck)
	}
	return session.WithStore(req, store)
}
