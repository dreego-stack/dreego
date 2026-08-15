package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestSelfClosingComponentAllowsFollowingSibling(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Card.dreego": "Component Card ()\n<div><article>{#slot}</article></div>",
		"dreego/routes/get.dreego":      `<div><@Card/><p>sibling</p></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<p>sibling</p>") {
		t.Fatalf("following sibling missing, got: %s", body)
	}
}
