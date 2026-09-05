package core

import (
	"strings"
	"testing"
)

// FuzzMarkdownToHTMLSafe asserts the SAFE runtime API never panics and never
// emits executable HTML for arbitrary hostile input. Run manually with
// `go test -fuzz=FuzzMarkdownToHTMLSafe ./core/`; CI does not run fuzz targets.
func FuzzMarkdownToHTMLSafe(f *testing.F) {
	seeds := []string{
		"<script>alert(1)</script>",
		`<img src=x onerror=alert(1)>`,
		`<iframe src=javascript:alert(1)>`,
		`<object data=javascript:alert(1)>`,
		`<embed src=javascript:alert(1)>`,
		`<svg onload=alert(1)>`,
		`[x](javascript:alert(1))`,
		`[x](JaVaScRiPt:alert(1))`,
		"[x](jav\tascript:alert(1))",
		"[x](java&#x09;script:alert(1))",
		`![x](data:text/html,<script>alert(1)</script>)`,
		`![x](data:image/svg+xml,<svg onload=alert(1)>)`,
		"```\"><script>alert(1)</script>\n```",
		`<scr<script>ipt>alert(1)</scr</script>ipt>`,
		"text[^1]\n\n[^1]: <script>alert(1)</script>",
		"# Hi\n\n[ok](https://example.com)",
		"<ifrAme",
		"hello <iframe src=javascript:alert(1)//",
		"hello <img src=x onerror=alert(1)",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 64<<10 {
			t.Skip("input too large")
		}
		out, err := MarkdownToHTML(input)
		if err != nil {
			return
		}
		lower := strings.ToLower(out)
		for _, bad := range []string{
			"<script",
			"<iframe",
			`href="javascript:`,
			`src="javascript:`,
			`src="data:text/html`,
			`src="data:image/svg`,
		} {
			if strings.Contains(lower, bad) {
				t.Fatalf("output contains executable HTML %q for input %q:\n%s", bad, input, out)
			}
		}
		for _, attr := range []string{"onerror=", "onload="} {
			if attrInUnescapedTag(lower, attr) {
				t.Fatalf("output contains executable attribute %q for input %q:\n%s", attr, input, out)
			}
		}
	})
}

// attrInUnescapedTag reports whether attr appears inside an unescaped <tag>,
// i.e. the nearest preceding '<' is not part of an escaped &lt; sequence and the
// tag is still open (no '>' between that '<' and the attr). An attr in plain
// text after a closed tag (e.g. "<p>onload=0</p>") is inert and not flagged.
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
			gt := strings.LastIndex(before, ">")
			if gt < lt {
				return true
			}
		}
		search = search[idx+len(attr):]
	}
}
