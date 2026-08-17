package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestReferenceHello(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeFixture(t, "hello")
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Hello from Dreego") {
		t.Fatalf("home page missing greeting: %s", body)
	}
	if !strings.Contains(body, "<title>Hello</title>") {
		t.Fatalf("home page missing title: %s", body)
	}
	code, body = c.Get(t, "/about")
	if code != 200 {
		t.Fatalf("GET /about = %d, want 200", code)
	}
	if !strings.Contains(body, "About this app") {
		t.Fatalf("about page missing content: %s", body)
	}
	code, body = c.Get(t, "/users/42")
	if code != 200 {
		t.Fatalf("GET /users/42 = %d, want 200", code)
	}
	if !strings.Contains(body, "User 42") {
		t.Fatalf("dynamic route missing param: %s", body)
	}
	code, _ = c.Get(t, "/missing")
	if code != 404 {
		t.Fatalf("GET /missing = %d, want 404", code)
	}
}

func TestReferenceForms(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeFixture(t, "forms")
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Guestbook") {
		t.Fatalf("guestbook page missing heading: %s", body)
	}
	code, body = c.Get(t, "/entries")
	if code != 200 {
		t.Fatalf("GET /entries = %d, want 200", code)
	}
	if !strings.Contains(body, "No entries yet") {
		t.Fatalf("empty entries page missing message: %s", body)
	}
	token := c.Cookie("csrf_token")
	if token == "" {
		t.Fatalf("no csrf_token cookie issued")
	}
	code, _, headers := c.Request(t, "POST", "/", "name=Ada&message=Hello+world&csrf_token="+token, map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("POST / = %d, want 303 redirect", code)
	}
	if !strings.Contains(headers.Get("Location"), "/entries") {
		t.Fatalf("POST / redirect location = %q, want /entries", headers.Get("Location"))
	}
	code, body = c.Get(t, "/entries")
	if code != 200 {
		t.Fatalf("GET /entries after POST = %d, want 200", code)
	}
	if !strings.Contains(body, "Ada") || !strings.Contains(body, "Hello world") {
		t.Fatalf("entry missing after POST: %s", body)
	}
	code, body = c.Get(t, "/counter")
	if code != 200 {
		t.Fatalf("GET /counter = %d, want 200", code)
	}
	if !strings.Contains(body, "Count: 0") {
		t.Fatalf("counter page missing initial count: %s", body)
	}
	code, _, _ = c.Request(t, "POST", "/counter", "csrf_token="+c.Cookie("csrf_token"), map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 200 {
		t.Fatalf("POST /counter = %d, want 200", code)
	}
	code, body = c.Get(t, "/counter")
	if code != 200 {
		t.Fatalf("GET /counter after POST = %d, want 200", code)
	}
	if !strings.Contains(body, "Count: 1") {
		t.Fatalf("counter did not increment: %s", body)
	}
}

func TestReferenceComponents(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeFixture(t, "components")
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Welcome to the shop") {
		t.Fatalf("shop page missing heading: %s", body)
	}
	if !strings.Contains(body, "Dreego Mug") {
		t.Fatalf("product card missing name: %s", body)
	}
	if !strings.Contains(body, "data-scope=") {
		t.Fatalf("component scoped style missing data-scope: %s", body)
	}
	if !strings.Contains(body, "In stock") {
		t.Fatalf("in-stock badge missing: %s", body)
	}
	if !strings.Contains(body, "Sold out") {
		t.Fatalf("sold-out badge missing: %s", body)
	}
	code, body = c.Get(t, "/products/1")
	if code != 200 {
		t.Fatalf("GET /products/1 = %d, want 200", code)
	}
	if !strings.Contains(body, "Dreego Mug") {
		t.Fatalf("product detail missing name: %s", body)
	}
	code, body = c.Get(t, "/products/2")
	if code != 200 {
		t.Fatalf("GET /products/2 = %d, want 200", code)
	}
	if !strings.Contains(body, "Dreego Tee") {
		t.Fatalf("product detail missing name: %s", body)
	}
}

func TestReferencePlugin(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeFixture(t, "plugin")
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Plugin demo") {
		t.Fatalf("plugin home missing heading: %s", body)
	}
	code, body = c.Get(t, "/plugin/hello")
	if code != 200 {
		t.Fatalf("GET /plugin/hello = %d, want 200", code)
	}
	if !strings.Contains(body, "Hello from the plugin") {
		t.Fatalf("plugin route missing body: %s", body)
	}
	code, body = c.Get(t, "/plugin/hello/42")
	if code != 200 {
		t.Fatalf("GET /plugin/hello/42 = %d, want 200", code)
	}
	if !strings.Contains(body, "Hello 42") {
		t.Fatalf("plugin dynamic route missing param: %s", body)
	}
	code, body = c.Get(t, "/plugin/health")
	if code != 200 {
		t.Fatalf("GET /plugin/health = %d, want 200", code)
	}
	if !strings.Contains(body, "plugin ok") {
		t.Fatalf("plugin health route missing body: %s", body)
	}
}
