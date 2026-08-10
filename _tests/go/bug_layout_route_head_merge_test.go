package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLayoutRouteHeadMerge(t *testing.T) {
	c := dreegotest.Serve(t, map[string]string{
		"dreego/layouts/default.dreego": `<html>
<head>
<meta charset="utf-8">
{#head}
</head>
<body><main>{#slot}</main></body>
</html>`,
		"dreego/routes/get.dreego": `<head><title>Merged Title</title><script src="route-script.js"></script></head>
<div><p>page content</p></div>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, `<meta charset="utf-8">`) {
		t.Fatalf("layout head meta missing in body: %s", body)
	}
	if !strings.Contains(body, "<title>Merged Title</title>") {
		t.Fatalf("route head title missing in body: %s", body)
	}
	if !strings.Contains(body, "route-script.js") {
		t.Fatalf("route head script missing in body: %s", body)
	}
}
