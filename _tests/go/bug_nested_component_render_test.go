package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugNestedNonSelfClosingComponentRenders(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Child.dreego":  "Component Child ()\n<div><section>{#slot}</section></div>",
		"dreego/components/Parent.dreego": "Component Parent ()\n<div><@Child><strong>inside</strong></@Child></div>",
		"dreego/routes/get.dreego":        "<div><@Parent/></div>",
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<@Child>") || !strings.Contains(body, "<strong>inside</strong>") {
		t.Fatalf("nested component did not render its slot: %s", body)
	}
}

func TestBugNestedSelfClosingComponentRenders(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Child.dreego":  "Component Child ()\n<div><strong>child</strong></div>",
		"dreego/components/Parent.dreego": "Component Parent ()\n<div><@Child/></div>",
		"dreego/routes/get.dreego":        "<div><@Parent/></div>",
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, "<strong>child</strong>") {
		t.Fatalf("nested self-closing component did not render: %s", body)
	}
}
