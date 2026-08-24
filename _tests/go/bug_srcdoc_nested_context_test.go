package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugSrcdocUsesNestedHTMLContext(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>payload := "<script>alert(1)</script>"</server>
<body><iframe srcdoc="{{ payload }}"></iframe></body>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, `srcdoc="&amp;lt;script&amp;gt;alert(1)&amp;lt;/script&amp;gt;"`) {
		t.Fatalf("srcdoc was not escaped for nested parsing: %s", body)
	}
}
