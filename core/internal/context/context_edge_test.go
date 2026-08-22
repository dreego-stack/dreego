package context

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSSRContextDataLifecycleAndTypedGet(t *testing.T) {
	c := NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if c.Data("missing") != nil || c.Get("missing") != "" {
		t.Fatal("missing data should return zero values")
	}
	c.Set("name", "dreego")
	c.Set("number", 42)
	if c.Data("name") != "dreego" || c.Get("name") != "dreego" {
		t.Fatalf("stored string not returned: data=%v get=%q", c.Data("name"), c.Get("name"))
	}
	if got := c.Get("number"); got != "" {
		t.Fatalf("non-string data returned as string: %q", got)
	}
	c.Delete("name")
	if c.Data("name") != nil {
		t.Fatal("deleted data still present")
	}
}

func TestSSRContextReadsQueryAndFormValues(t *testing.T) {
	r := httptest.NewRequest("POST", "/search?q=dreego", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Body = http.NoBody
	c := NewSSR(httptest.NewRecorder(), r)
	if got := c.Query("q"); got != "dreego" {
		t.Fatalf("query: got %q", got)
	}
	if got := c.FormValue("name"); got != "" || c.FormError() != nil {
		t.Fatalf("empty form value: got %q error=%v", got, c.FormError())
	}
}

func TestSSRContextRedirectReturnsSentinelAndLocation(t *testing.T) {
	rec := httptest.NewRecorder()
	c := NewSSR(rec, httptest.NewRequest("GET", "/from", nil))
	if err := c.Redirect("/to", http.StatusSeeOther); !errors.Is(err, ErrRedirect) {
		t.Fatalf("Redirect error: %v", err)
	}
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/to" {
		t.Fatalf("redirect response: status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
}

func TestSSRContextRequestMetadataDefaults(t *testing.T) {
	c := NewSSR(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if c.RequestID() != "" {
		t.Fatalf("unexpected request id: %q", c.RequestID())
	}
	if c.SessionError() != nil || c.FormError() != nil {
		t.Fatal("new context should not contain errors")
	}
}
