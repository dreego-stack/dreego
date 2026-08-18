package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugRouteInlineScriptBodyStaysLiteral(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<div><script>const template = "{{ literal }}";</script></div>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, `const template = "{{ literal }}";`) {
		t.Fatalf("script body was changed: %s", body)
	}
}
