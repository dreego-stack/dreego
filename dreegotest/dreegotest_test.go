package dreegotest_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	dreego "codeberg.org/dreego/dreego/core"
	"codeberg.org/dreego/dreego/dreegotest"
)

// Public API contract under test:
//
//	dreegotest.Get(t, path) *Response                 // simulates a GET request through dreego.ServeMux()
//	dreegotest.PostForm(t, path, form) *Response      // simulates an application/x-www-form-urlencoded POST
//	dreegotest.RenderComponent(t, fn, props...) string // renders a single component, props as key/value pairs
//
// *Response exposes:
//   - StatusCode int
//   - Body       string

func TestGetReturnsStatusAndBody(t *testing.T) {
	dreego.Reset()
	dreego.Register("GET", "/dgt-hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World"))
	})
	defer dreego.Reset()

	resp := dreegotest.Get(t, "/dgt-hello")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != "Hello, World" {
		t.Errorf("Body = %q, want %q", resp.Body, "Hello, World")
	}
}

func TestGetReturnsNonOKStatus(t *testing.T) {
	dreego.Reset()
	dreego.Register("GET", "/dgt-missing", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	})
	defer dreego.Reset()

	resp := dreegotest.Get(t, "/dgt-missing")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestPostFormBindsValues(t *testing.T) {
	dreego.Reset()
	dreego.Register("POST", "/dgt-submit", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("got:" + name))
	})
	defer dreego.Reset()

	form := url.Values{"name": {"world"}}
	resp := dreegotest.PostForm(t, "/dgt-submit", form)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if !strings.Contains(resp.Body, "got:world") {
		t.Errorf("Body = %q, want it to contain %q", resp.Body, "got:world")
	}
}

func TestPostFormSetsFormContentType(t *testing.T) {
	dreego.Reset()
	var contentType string
	dreego.Register("POST", "/dgt-ctype", func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	})
	defer dreego.Reset()

	dreegotest.PostForm(t, "/dgt-ctype", url.Values{"a": {"1"}})
	if !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		t.Errorf("Content-Type = %q, want application/x-www-form-urlencoded", contentType)
	}
}

func TestRenderComponentReturnsHTML(t *testing.T) {
	comp := dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {
		return "<section><h1>" + ctx.Get("title") + "</h1></section>", nil
	})
	out := dreegotest.RenderComponent(t, comp, "title", "Welcome")
	if !strings.Contains(out, "<h1>Welcome</h1>") {
		t.Errorf("output = %q, want it to contain <h1>Welcome</h1>", out)
	}
}

func TestRenderComponentEscapesXSS(t *testing.T) {
	comp := dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {
		return "<p>" + ctx.Get("name") + "</p>", nil
	})
	out := dreegotest.RenderComponent(t, comp, "name", "<script>alert(1)</script>")
	if strings.Contains(out, "<script>alert(1)</script>") {
		t.Errorf("output must not contain raw script tag, got %q", out)
	}
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Errorf("output must contain escaped script, got %q", out)
	}
}
