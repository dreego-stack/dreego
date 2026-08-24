package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestOutputContextTextEscapesMarkup(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>v := "<script>alert(1)</script>"</server>
<body><p>{{ v }}</p></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "<script>alert(1)</script>")
	dreegotest.MustContainBody(t, body, "&lt;script&gt;alert(1)&lt;/script&gt;")
}

func TestOutputContextAttrEscapesQuotes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>v := ` + "`\" onmouseover=\"alert(1)`" + `</server>
<body><a title="{{ v }}">link</a></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, `" onmouseover="`)
	dreegotest.MustContainBody(t, body, "&#34; onmouseover=&#34;alert(1)")
}

func TestOutputContextURLRejectsJavascriptScheme(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>u := "javascript:alert(1)"</server>
<body><a href="{{ u }}">link</a></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "javascript:")
	dreegotest.MustContainBody(t, body, `href="#`)
}

func TestOutputContextURLAllowsHTTPS(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>u := "https://example.com/x"</server>
<body><a href="{{ u }}">link</a></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `href="https://example.com/x"`)
}

func TestOutputContextURLRejectsDataScheme(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>u := "data:text/html,<script>alert(1)</script>"</server>
<body><img src="{{ u }}"></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "data:text/html")
	dreegotest.MustContainBody(t, body, `src="#`)
}

func TestOutputContextScriptAttrJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>s := ` + "`\"><script>alert(1)</script>`" + `</server>
<body><button onclick="{{ s }}">go</button></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "<client>alert(1)</client>")
	dreegotest.MustContainBody(t, body, `\u003cscript\u003e`)
	dreegotest.MustContainBody(t, body, "&#34;")
}

func TestOutputContextStyleNeutralizesBreakout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>s := "red; } </style><script>alert(1)</script>"</server>
<body><div style="{{ s }}">x</div></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "</style>")
	dreegotest.MustNotContainBody(t, body, "<client>alert(1)</client>")
}

func TestOutputContextRawOptInPassesThrough(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>h := "<b>trusted</b>"</server>
<body><p>{{ h|raw }}</p></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, "<b>trusted</b>")
}

func TestOutputContextComponentURLRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Link.dreego": `Component Link (url string)
<body><a href="{{ url }}">go</a></body>`,
		"www/routes/get.dreego": `<server>u := "javascript:alert(1)"</server>
<body><@Link url={u}/></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "javascript:")
	dreegotest.MustContainBody(t, body, `href="#`)
}
