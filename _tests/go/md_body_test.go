package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestMdBodyServes asserts a <body lang="md"> route renders Markdown to HTML at
// generation time: a pipe table, a footnote section, a protected {#if} block,
// and an expression all survive into the served page.
func TestMdBodyServes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>name := "Ada"</server>
<body lang="md">
# Account

| Feature | Status |
| --- | --- |
| Tables | yes |

See the note[^1].

{#if name == "Ada"}
Hello **{{ name }}**.
{/if}

[^1]: The footnote text.
</body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	for _, want := range []string{
		"<h1>Account</h1>",
		"<table>",
		"<th>Feature</th>",
		"<td>Tables</td>",
		`<section class="footnotes">`,
		"The footnote text.",
		"Hello **Ada**.",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q:\n%s", want, body)
		}
	}
	if strings.Contains(body, "{{ name }}") {
		t.Fatalf("unresolved {{ name }} found in HTML: %s", body)
	}
}
