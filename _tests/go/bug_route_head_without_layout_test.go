package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugRouteHeadWithoutLayout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<head><title>No Layout Title</title><script src="route.js"></script></head>
<div><p>hello</p></div>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, "<title>No Layout Title</title>") {
		t.Fatalf("route head title missing in body: %s", body)
	}
	if !strings.Contains(body, "route.js") {
		t.Fatalf("route head script missing in body: %s", body)
	}
}
