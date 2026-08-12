package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugErrorPageDoctype(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/404.dreego": `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Not Found</title>
</head>
<body>
    <div><p>Not Found</p></div>
</body>
</html>`,
	})
	code, body := c.Get(t, "/missing")
	if code != 404 {
		t.Fatalf("expected HTTP 404, got %d", code)
	}
	if !strings.HasPrefix(body, "<!doctype html>") {
		t.Fatalf("body must start with <!doctype html>, got: %q", body[:min(80, len(body))])
	}
	if strings.Contains(body, "data-scope") {
		t.Fatalf("scope div must not appear in 404 body: %s", body)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
