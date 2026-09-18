package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugNestedNonSelfClosingComponentRenders(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Child.dreego":  "DREEFILE component ()\n<body><section>{#slot}</section></body>",
		"www/components/Parent.dreego": "DREEFILE component ()\n<body><@Child><strong>inside</strong></@Child></body>",
		"www/routes/+page.dreego":      "<body><@Parent/></body>",
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<@Child>") || !strings.Contains(body, "<strong>inside</strong>") {
		t.Fatalf("nested component did not render its slot: %s", body)
	}
}

func TestBugNestedSelfClosingComponentRenders(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Child.dreego":  "DREEFILE component ()\n<body><strong>child</strong></body>",
		"www/components/Parent.dreego": "DREEFILE component ()\n<body><@Child/></body>",
		"www/routes/+page.dreego":      "<body><@Parent/></body>",
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, "<strong>child</strong>") {
		t.Fatalf("nested self-closing component did not render: %s", body)
	}
}
