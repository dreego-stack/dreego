package dreegotest_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/dreegotest"
)

func TestGetReturnsStatusAndBody(t *testing.T) {
	app := dreego.New()
	app.Register("GET", "/dgt-hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World"))
	})
	client := dreegotest.NewApp(app)
	resp := client.Get(t, "/dgt-hello")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Body != "Hello, World" {
		t.Errorf("Body = %q, want %q", resp.Body, "Hello, World")
	}
}

func TestGetReturnsNonOKStatus(t *testing.T) {
	app := dreego.New()
	app.Register("GET", "/dgt-missing", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	})
	resp := dreegotest.NewApp(app).Get(t, "/dgt-missing")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestPostFormBindsValues(t *testing.T) {
	app := dreego.New()
	app.Register("POST", "/dgt-submit", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("got:" + name))
	})
	form := url.Values{"name": {"world"}}
	resp := dreegotest.NewApp(app).PostForm(t, "/dgt-submit", form)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if !strings.Contains(resp.Body, "got:world") {
		t.Errorf("Body = %q, want it to contain %q", resp.Body, "got:world")
	}
}

func TestPostFormSetsFormContentType(t *testing.T) {
	app := dreego.New()
	var contentType string
	app.Register("POST", "/dgt-ctype", func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	})
	dreegotest.NewApp(app).PostForm(t, "/dgt-ctype", url.Values{"a": {"1"}})
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
