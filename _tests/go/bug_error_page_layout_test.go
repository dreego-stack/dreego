package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugErrorPageLayout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego": `<head>
    <title>Layout Site</title>
    {#head}
</head>
<body><main>{#slot}</main></body>`,
		"www/routes/404.dreego": `<head>
    <meta charset="utf-8">
    <title>Not Found</title>
    <link rel="stylesheet" href="/err.css">
</head>
<body><!doctype html>
<html lang="en">
<body>
    <div><p>Not Found</p></div>
</body>
</html></body>
<style>
p { color: red; }
</style>
<client>
console.log("err");
</client>`,
	})
	code, body := c.Get(t, "/missing")
	if code != 404 {
		t.Fatalf("expected HTTP 404, got %d", code)
	}
	if !strings.HasPrefix(body, "<!doctype html>") {
		t.Fatalf("body must start with <!doctype html>, got: %q", body[:min(80, len(body))])
	}
	if strings.Contains(body, "Layout Site") {
		t.Fatalf("layout must not wrap 404 page, got: %s", body)
	}
	if !strings.Contains(body, "Not Found") {
		t.Fatalf("404 title missing in body: %s", body)
	}
}
