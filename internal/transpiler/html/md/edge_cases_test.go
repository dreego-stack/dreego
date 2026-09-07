package md

import (
	"strings"
	"testing"
)

func TestMarkdownEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "empty", src: "", want: ""},
		{name: "blank lines", src: "\n\n\n", want: ""},
		{name: "unicode heading", src: "# Grüße 👋", want: "<h1>Grüße 👋</h1>"},
		{name: "windows lines", src: "# One\r\n\r\nTwo", want: "<h1>One</h1><p>Two</p>"},
		{name: "ampersand", src: "A & B", want: "<p>A &amp; B</p>"},
		{name: "quoted code", src: "`<tag>&`", want: "<p><code>&lt;tag&gt;&amp;</code></p>"},
		{name: "empty link text", src: "[](https://example.com)", want: `<p><a href="https://example.com"></a></p>`},
		{name: "relative link", src: "[docs](../guide)", want: `<p><a href="../guide">docs</a></p>`},
		{name: "fragment link", src: "[jump](#section)", want: `<p><a href="#section">jump</a></p>`},
		{name: "unicode image alt", src: "![Überblick](diagram.png)", want: `<p><img src="diagram.png" alt="Überblick"></p>`},
		{name: "empty code fence", src: "```\n```", want: "<pre><code></code></pre>"},
		{name: "code preserves newlines", src: "```go\na := 1\nb := 2\n```", want: `<pre><code class="language-go">a := 1
b := 2</code></pre>`},
		{name: "ordered start value normalizes", src: "7. seven\n8. eight", want: "<ol><li>seven</li><li>eight</li></ol>"},
		{name: "mixed list transition", src: "- one\n1. two", want: "<ul><li>one</li></ul><ol><li>two</li></ol>"},
		{name: "blockquote escapes", src: "> 2 < 3 & 4", want: "<blockquote>2 &lt; 3 &amp; 4</blockquote>"},
		{name: "table escapes cells", src: "| a | b |\n| --- | --- |\n| < | & |", want: "<table><thead><tr><th>a</th><th>b</th></tr></thead><tbody><tr><td>&lt;</td><td>&amp;</td></tr></tbody></table>"},
		{name: "adjacent emphasis", src: "**strong***em*", want: "<p><strong>strong</strong><em>em</em></p>"},
		{name: "plain underscores", src: "snake_case_name", want: "<p>snake_case_name</p>"},
		{name: "trailing spaces", src: "text   ", want: "<p>text</p>"},
		{name: "html comment escaped", src: "<!-- note -->", want: "<p>&lt;!-- note --&gt;</p>"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			nodes, err := ToNodes(test.src, ModeTrusted)
			if err != nil {
				t.Fatalf("ToNodes() error = %v", err)
			}
			var got strings.Builder
			for _, node := range nodes {
				got.WriteString(node.Content)
			}
			if got.String() != test.want {
				t.Fatalf("ToNodes() = %q, want %q", got.String(), test.want)
			}
		})
	}
}

func TestSafeMarkdownNeverEmitsExecutableHTML(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"<script>alert(1)</script>",
		"<svg onload=alert(1)>",
		"<iframe srcdoc='<script>alert(1)</script>'></iframe>",
		"[click](javascript:alert(1))",
		"![image](data:image/svg+xml,<svg onload=alert(1)>)",
		"<math href=javascript:alert(1)>",
	}

	for _, input := range inputs {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			nodes, err := ToNodes(input, ModeSafe)
			if err != nil {
				t.Fatalf("ToNodes() error = %v", err)
			}
			for _, node := range nodes {
				lower := strings.ToLower(node.Content)
				for _, forbidden := range []string{"<script", "<iframe", "<svg", "<math", `href="javascript:`, `src="data:image/svg`} {
					if strings.Contains(lower, forbidden) {
						t.Fatalf("safe output contains %q: %q", forbidden, node.Content)
					}
				}
			}
		})
	}
}
