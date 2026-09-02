package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestMdInlineServes asserts an inline <md> tag inside a lang="html" body is
// converted to a <div> wrapper around the markdown output, with the surrounding
// HTML preserved in order.
func TestMdInlineServes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<body><h1>HTML title</h1><md class="prose">
# Markdown inside

- a
- b
</md><p>After</p></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	for _, want := range []string{
		"<h1>HTML title</h1>",
		`<div class="prose">`,
		"<h1>Markdown inside</h1>",
		"<ul>",
		"<li>a</li>",
		"<li>b</li>",
		"</ul>",
		"<p>After</p>",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q:\n%s", want, body)
		}
	}
	htmlIdx := strings.Index(body, "<h1>HTML title</h1>")
	mdIdx := strings.Index(body, "<h1>Markdown inside</h1>")
	afterIdx := strings.Index(body, "<p>After</p>")
	if !(htmlIdx < mdIdx && mdIdx < afterIdx) {
		t.Fatalf("HTML parts out of order: html=%d md=%d after=%d\n%s", htmlIdx, mdIdx, afterIdx, body)
	}
}
