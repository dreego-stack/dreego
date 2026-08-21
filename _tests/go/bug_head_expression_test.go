package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadExpression(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>doc := struct{ Title string }{Title: "My Docs Title"}</go>
<head><title>{{ doc.Title }}</title></head>
<div><h1>{{ doc.Title }}</h1></div>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, "<title>My Docs Title</title>") {
		t.Fatalf("head expression not resolved, got: %s", body)
	}
	if strings.Contains(body, "{{ doc.Title }}") {
		t.Fatalf("unresolved {{ doc.Title }} found in HTML: %s", body)
	}
}
