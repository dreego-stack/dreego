package md

import (
	"strings"
	"sync"
	"testing"
)

func TestToHTMLBlocks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "empty", src: "", want: ""},
		{name: "blank lines", src: "\n\n", want: ""},
		{name: "heading", src: "# Title", want: "<h1>Title</h1>"},
		{name: "paragraph", src: "hello world", want: "<p>hello world</p>"},
		{name: "paragraph multi line", src: "first\nsecond", want: "<p>first second</p>"},
		{name: "emphasis", src: "*em* and **strong**", want: "<p><em>em</em> and <strong>strong</strong></p>"},
		{name: "inline code", src: "use `code` here", want: "<p>use <code>code</code> here</p>"},
		{name: "link", src: "[text](https://example.com)", want: `<p><a href="https://example.com">text</a></p>`},
		{name: "image", src: "![alt](https://example.com/i.png)", want: `<p><img src="https://example.com/i.png" alt="alt"></p>`},
		{name: "unordered list", src: "- a\n- b", want: "<ul><li>a</li><li>b</li></ul>"},
		{name: "ordered list", src: "1. a\n2. b", want: "<ol><li>a</li><li>b</li></ol>"},
		{name: "nested list", src: "- a\n  - b\n- c", want: "<ul><li>a<ul><li>b</li></ul></li><li>c</li></ul>"},
		{name: "blockquote", src: "> quoted", want: "<blockquote>quoted</blockquote>"},
		{name: "horizontal rule", src: "---", want: "<hr>"},
		{name: "fenced code", src: "```go\nfunc main() {}\n```", want: `<pre><code class="language-go">func main() {}</code></pre>`},
		{name: "fenced code escaped", src: "```\n<b>raw</b>\n```", want: "<pre><code>&lt;b&gt;raw&lt;/b&gt;</code></pre>"},
		{name: "table", src: "| a | b |\n| --- | --- |\n| x | y |", want: "<table><thead><tr><th>a</th><th>b</th></tr></thead><tbody><tr><td>x</td><td>y</td></tr></tbody></table>"},
		{name: "footnote", src: "text[^1]\n\n[^1]: note", want: `<p>text<sup class="footnote-ref"><a href="#fn-1" id="fnref-1">1</a></sup></p><section class="footnotes"><ol><li id="fn-1">note <a href="#fnref-1" class="footnote-backref">↩</a></li></ol></section>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ToHTML(tt.src, ModeTrusted)
			if err != nil {
				t.Fatalf("ToHTML() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ToHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToHTMLSafeEscapesRawHTML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "script block escaped", src: "<script>alert(1)</script>", want: "&lt;script&gt;alert(1)&lt;/script&gt;"},
		{name: "inline img escaped", src: "before <img onerror=alert(1)> after", want: "<p>before &lt;img onerror=alert(1)&gt; after</p>"},
		{name: "unclosed tag escaped", src: "<iframe src=x", want: "&lt;iframe src=x"},
		{name: "javascript link blocked", src: "[x](javascript:alert(1))", want: "<p>[x](javascript:alert(1))</p>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ToHTML(tt.src, ModeSafe)
			if err != nil {
				t.Fatalf("ToHTML() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("ToHTML() = %q, want %q", got, tt.want)
			}
			lower := strings.ToLower(got)
			for _, forbidden := range []string{"<script", "<iframe", "<svg", "<math", `href="javascript:`, `src="data:image/svg`} {
				if strings.Contains(lower, forbidden) {
					t.Fatalf("safe output contains %q: %q", forbidden, got)
				}
			}
		})
	}
}

func TestToHTMLTrustedPassesRawHTML(t *testing.T) {
	t.Parallel()

	got, err := ToHTML("<div class=\"note\">raw</div>", ModeTrusted)
	if err != nil {
		t.Fatalf("ToHTML() error = %v", err)
	}
	if got != "<div class=\"note\">raw</div>" {
		t.Fatalf("ToHTML() = %q, want raw passthrough", got)
	}
}

func TestToHTMLBareListMarkerErrors(t *testing.T) {
	t.Parallel()

	_, err := ToHTML("-", ModeTrusted)
	if err == nil {
		t.Fatal("expected error for bare list marker")
	}
	if !strings.Contains(err.Error(), "line 1") {
		t.Fatalf("error should mention line number, got: %v", err)
	}
}

func TestParseBlocksReturnsOneStringPerBlock(t *testing.T) {
	t.Parallel()

	blocks, err := ParseBlocks("# Title\n\npara *em*", ModeTrusted)
	if err != nil {
		t.Fatalf("ParseBlocks() error = %v", err)
	}
	want := []string{"<h1>Title</h1>", "<p>para <em>em</em></p>"}
	if len(blocks) != len(want) {
		t.Fatalf("got %d blocks, want %d: %v", len(blocks), len(want), blocks)
	}
	for i := range want {
		if blocks[i] != want[i] {
			t.Errorf("block %d = %q, want %q", i, blocks[i], want[i])
		}
	}
}

func TestSafeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		raw     string
		isImage bool
		want    string
	}{
		{"https://ok.com", false, "https://ok.com"},
		{"mailto:hi@example.com", false, "mailto:hi@example.com"},
		{"/rel/path", false, "/rel/path"},
		{"#fragment", false, "#fragment"},
		{"javascript:alert(1)", false, ""},
		{"vbscript:x", false, ""},
		{"data:image/png;base64,AAA", true, "data:image/png;base64,AAA"},
		{"data:text/html,x", true, ""},
		{"data:image/svg+xml;base64,EE", true, ""},
		{"https://ok.com/?a=1&amp;b=2", false, "https://ok.com/?a=1&b=2"},
	}

	for _, tt := range tests {
		got := SafeURL(tt.raw, tt.isImage)
		if got != tt.want {
			t.Errorf("SafeURL(%q, %v) = %q, want %q", tt.raw, tt.isImage, got, tt.want)
		}
	}
}

func TestSafeFenceLanguageRejectsUnsafe(t *testing.T) {
	t.Parallel()

	if got := SafeFenceLanguage("go"); got != "go" {
		t.Errorf("SafeFenceLanguage(go) = %q, want go", got)
	}
	if got := SafeFenceLanguage("\"><script>"); got != "" {
		t.Errorf("SafeFenceLanguage(unsafe) = %q, want empty", got)
	}
}

func TestConcurrentToHTML(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mode := ModeSafe
			if i%2 == 0 {
				mode = ModeTrusted
			}
			if _, err := ToHTML("# h\n\npara *em* [x](https://ok.com)", mode); err != nil {
				t.Errorf("ToHTML() error = %v", err)
			}
		}(i)
	}
	wg.Wait()
}
