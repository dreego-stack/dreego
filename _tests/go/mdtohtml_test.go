package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestMdtohtmlTrustedGenerated asserts the dreego.mdtohtml(x, trusted: true)
// stdlib syntax is rewritten to the trusted core call in the generated code.
func TestMdtohtmlTrustedGenerated(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<server>post := "# Hi"
html, err := dreego.mdtohtml(post, trusted: true)
if err != nil { return "", err }</server>
<body><p>{{ html|raw }}</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "dreego.MarkdownToHTMLTrusted(post)")
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "dreego.mdtohtml(")
}

// TestMdtohtmlSafeGenerated asserts the default dreego.mdtohtml(x) syntax is
// rewritten to the safe core call in the generated code.
func TestMdtohtmlSafeGenerated(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/get.dreego": `<server>post := "# Hi"
html, err := dreego.mdtohtml(post)
if err != nil { return "", err }</server>
<body><p>{{ html|raw }}</p></body>`,
	})
	dreegotest.MustContain(t, gen["www/routes/dree.go"], "dreego.MarkdownToHTML(post)")
	dreegotest.MustNotContain(t, gen["www/routes/dree.go"], "dreego.mdtohtml(")
}

// TestMdtohtmlServes asserts a route using dreego.mdtohtml renders the Markdown
// to HTML at runtime.
func TestMdtohtmlServes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>html, err := dreego.mdtohtml("# Hello")
if err != nil { return "", err }</server>
<body><div>{{ html|raw }}</div></body>`,
	})
	code, body := c.Get(t, "/")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "<h1>Hello</h1>")
}
