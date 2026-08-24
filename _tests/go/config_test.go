package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestConfigInvalidJSON(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json": `{ broken json !!!`,
		"www/routes/get.dreego":  `<body><p>hello</p></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if !strings.Contains(body, "hello") {
		t.Fatalf("expected 'hello' in response, got: %s", body)
	}
}

func TestConfigRedirect(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json": `{
    "redirects": [
        { "from": "/old", "to": "/new", "status": 301 }
    ]
}`,
		"www/routes/get.dreego": `<body><p>new page</p></body>`,
	})
	code, _, headers := c.Request(t, "GET", "/old", "", nil)
	if code != 301 {
		t.Fatalf("expected 301 redirect, got %d", code)
	}
	if loc := headers.Get("Location"); !strings.HasSuffix(loc, "/new") {
		t.Fatalf("expected redirect to /new, got %q", loc)
	}
}

func TestConfigRewrite(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json": `{
    "rewrites": [
        { "from": "/old", "to": "/new" }
    ]
}`,
		"www/routes/new/index.dreego": `<body><p>rewritten content</p></body>`,
	})
	code, body := c.Get(t, "/old")
	if code != 200 {
		t.Fatalf("expected 200 for rewritten path, got %d", code)
	}
	if !strings.Contains(body, "rewritten content") {
		t.Fatalf("expected rewritten content, got: %s", body)
	}
}
