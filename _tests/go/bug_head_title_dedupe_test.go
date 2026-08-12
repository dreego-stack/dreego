package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadTitleDedupe(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/layouts/default.dreego": `<head>
    <title>Site</title>
    <meta name="description" content="site desc">
    {#head}
</head>
<div><main>{#slot}</main></div>`,
		"dreego/routes/get.dreego": `<head><title>Page</title><meta name="description" content="route desc"></head>
<div><h1>Page</h1></div>`,
	})
	_, body := c.Get(t, "/")
	if n := strings.Count(body, "<title>"); n != 1 {
		t.Fatalf("expected exactly 1 <title>, got %d: %s", n, body)
	}
	if !strings.Contains(body, "<title>Page</title>") {
		t.Fatalf("route title missing in body: %s", body)
	}
	if strings.Contains(body, "Site") {
		t.Fatalf("layout title still present in body: %s", body)
	}
	if n := strings.Count(body, `name="description"`); n != 1 {
		t.Fatalf("expected exactly 1 meta description, got %d: %s", n, body)
	}
	if !strings.Contains(body, `content="route desc"`) {
		t.Fatalf("route meta description missing in body: %s", body)
	}
	if strings.Contains(body, "site desc") {
		t.Fatalf("layout meta description still present in body: %s", body)
	}
}
