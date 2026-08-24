package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugRouteInlineScriptBodyStaysLiteral(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<body><script>const template = "{{ literal }}";</script></body>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, `const template = "{{ literal }}";`) {
		t.Fatalf("script body was changed: %s", body)
	}
}
