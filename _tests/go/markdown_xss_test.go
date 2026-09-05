package tests

import (
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

// TestMarkdownXSS asserts the SAFE runtime API (dreego.MarkdownToHTML) never
// emits executable HTML for hostile user content coming from a database.
func TestMarkdownXSS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		src  string
	}{
		{name: "script block", src: "<script>alert(1)</script>"},
		{name: "img onerror", src: `<img src=x onerror=alert(1)>`},
		{name: "iframe javascript", src: `<iframe src=javascript:alert(1)>`},
		{name: "object", src: `<object data=javascript:alert(1)>`},
		{name: "embed", src: `<embed src=javascript:alert(1)>`},
		{name: "svg onload", src: `<svg onload=alert(1)>`},
		{name: "markdown link javascript", src: `[x](javascript:alert(1))`},
		{name: "markdown link mixed case", src: `[x](JaVaScRiPt:alert(1))`},
		{name: "markdown link control char", src: "[x](jav\tascript:alert(1))"},
		{name: "markdown link entity encoded", src: "[x](java&#x09;script:alert(1))"},
		{name: "image data html", src: `![x](data:text/html,<script>alert(1)</script>)`},
		{name: "image data svg", src: `![x](data:image/svg+xml,<svg onload=alert(1)>)`},
		{name: "fence lang breakout", src: "```\"><script>alert(1)</script>\n```"},
		{name: "tag soup", src: `<scr<script>ipt>alert(1)</scr</script>ipt>`},
		{name: "footnote injection", src: "text[^1]\n\n[^1]: <script>alert(1)</script>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := dreego.MarkdownToHTML(tt.src)
			if err != nil {
				t.Fatalf("MarkdownToHTML() error = %v", err)
			}
			assertNoExecutableHTML(t, out)
		})
	}
}

// assertNoExecutableHTML fails if out contains any construct a browser would
// execute. Safe mode escapes raw HTML to &lt;...&gt;, so inert text such as
// "onerror=" inside "&lt;img onerror=...&gt;" is allowed; only unescaped
// executable forms are rejected. Blocked markdown links render as literal text
// (e.g. "javascript:" inside "[x](javascript:...)"), which is also inert.
func assertNoExecutableHTML(t *testing.T, out string) {
	t.Helper()
	lower := strings.ToLower(out)
	for _, bad := range []string{
		"<script",
		"<iframe",
		"<object",
		"<embed",
		"<svg",
		`href="javascript:`,
		`src="javascript:`,
		`src="data:text/html`,
		`src="data:image/svg`,
	} {
		if strings.Contains(lower, bad) {
			t.Errorf("output contains executable HTML %q:\n%s", bad, out)
		}
	}
	for _, attr := range []string{"onerror=", "onload="} {
		if attrInUnescapedTag(lower, attr) {
			t.Errorf("output contains executable attribute %q:\n%s", attr, out)
		}
	}
}

// attrInUnescapedTag reports whether attr appears inside an unescaped <tag>,
// i.e. the nearest preceding '<' is not part of an escaped &lt; sequence.
func attrInUnescapedTag(s, attr string) bool {
	search := s
	for {
		idx := strings.Index(search, attr)
		if idx < 0 {
			return false
		}
		before := search[:idx]
		lt := strings.LastIndex(before, "<")
		esc := strings.LastIndex(before, "&lt;")
		if lt > esc {
			return true
		}
		search = search[idx+len(attr):]
	}
}

// TestMarkdownXSSEscaped asserts the escaped forms are present where the raw
// HTML was neutralized rather than dropped.
func TestMarkdownXSSEscaped(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "script escaped", src: "<script>alert(1)</script>", want: "&lt;script&gt;"},
		{name: "img escaped", src: `<img src=x onerror=alert(1)>`, want: "&lt;img"},
		{name: "iframe escaped", src: `<iframe src=javascript:alert(1)>`, want: "&lt;iframe"},
		{name: "svg escaped", src: `<svg onload=alert(1)>`, want: "&lt;svg"},
		{name: "fence lang escaped", src: "```\"><script>alert(1)</script>\n```", want: "&lt;script&gt;"},
		{name: "footnote escaped", src: "text[^1]\n\n[^1]: <script>alert(1)</script>", want: "&lt;script&gt;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := dreego.MarkdownToHTML(tt.src)
			if err != nil {
				t.Fatalf("MarkdownToHTML() error = %v", err)
			}
			if !strings.Contains(out, tt.want) {
				t.Errorf("output missing escaped form %q:\n%s", tt.want, out)
			}
		})
	}
}

// TestMarkdownXSSBlockedLinksLiteral asserts blocked URL schemes are rendered
// as inert literal text, not as executable anchors or images.
func TestMarkdownXSSBlockedLinksLiteral(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "javascript link", src: `[x](javascript:alert(1))`, want: "[x](javascript:alert(1))"},
		{name: "mixed case link", src: `[x](JaVaScRiPt:alert(1))`, want: "[x](JaVaScRiPt:alert(1))"},
		{name: "control char link", src: "[x](jav\tascript:alert(1))", want: "[x](jav\tascript:alert(1))"},
		{name: "entity encoded link", src: "[x](java&#x09;script:alert(1))", want: "[x](java&amp;#x09;script:alert(1))"},
		{name: "data html image", src: `![x](data:text/html,<script>)`, want: "![x](data:text/html,"},
		{name: "data svg image", src: `![x](data:image/svg+xml,<svg>)`, want: "![x](data:image/svg+xml,"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := dreego.MarkdownToHTML(tt.src)
			if err != nil {
				t.Fatalf("MarkdownToHTML() error = %v", err)
			}
			if !strings.Contains(out, tt.want) {
				t.Errorf("output missing literal blocked link %q:\n%s", tt.want, out)
			}
			if strings.Contains(out, "<a href=") || strings.Contains(out, "<img ") {
				t.Errorf("blocked link rendered as executable element:\n%s", out)
			}
		})
	}
}

// TestMarkdownXSSPositiveControl asserts safe markdown still renders normally.
func TestMarkdownXSSPositiveControl(t *testing.T) {
	t.Parallel()
	out, err := dreego.MarkdownToHTML("# Hi\n\n[ok](https://example.com)")
	if err != nil {
		t.Fatalf("MarkdownToHTML() error = %v", err)
	}
	for _, want := range []string{
		"<h1>Hi</h1>",
		`<a href="https://example.com">ok</a>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}
