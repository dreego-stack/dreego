package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestOutputContextTextEscapesMarkup(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>v := "<script>alert(1)</script>"</go>
<div><p>{{ v }}</p></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "<script>alert(1)</script>")
	dreegotest.MustContainBody(t, body, "&lt;script&gt;alert(1)&lt;/script&gt;")
}

func TestOutputContextAttrEscapesQuotes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>v := ` + "`\" onmouseover=\"alert(1)`" + `</go>
<div><a title="{{ v }}">link</a></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, `" onmouseover="`)
	dreegotest.MustContainBody(t, body, "&#34; onmouseover=&#34;alert(1)")
}

func TestOutputContextURLRejectsJavascriptScheme(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>u := "javascript:alert(1)"</go>
<div><a href="{{ u }}">link</a></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "javascript:")
	dreegotest.MustContainBody(t, body, `href="#`)
}

func TestOutputContextURLAllowsHTTPS(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>u := "https://example.com/x"</go>
<div><a href="{{ u }}">link</a></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `href="https://example.com/x"`)
}

func TestOutputContextURLRejectsDataScheme(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>u := "data:text/html,<script>alert(1)</script>"</go>
<div><img src="{{ u }}"></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "data:text/html")
	dreegotest.MustContainBody(t, body, `src="#`)
}

func TestOutputContextScriptAttrJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>s := ` + "`\"><script>alert(1)</script>`" + `</go>
<div><button onclick="{{ s }}">go</button></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "<script>alert(1)</script>")
	dreegotest.MustContainBody(t, body, `\u003cscript\u003e`)
	dreegotest.MustContainBody(t, body, "&#34;")
}

func TestOutputContextStyleNeutralizesBreakout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>s := "red; } </style><script>alert(1)</script>"</go>
<div><div style="{{ s }}">x</div></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "</style>")
	dreegotest.MustNotContainBody(t, body, "<script>alert(1)</script>")
}

func TestOutputContextRawOptInPassesThrough(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<go>h := "<b>trusted</b>"</go>
<div><p>{{ h|raw }}</p></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, "<b>trusted</b>")
}

func TestOutputContextComponentURLRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Link.dreego": `Component Link (url string)
<div><a href="{{ url }}">go</a></div>`,
		"www/routes/get.dreego": `<go>u := "javascript:alert(1)"</go>
<div><@Link url={u}/></div>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustNotContainBody(t, body, "javascript:")
	dreegotest.MustContainBody(t, body, `href="#`)
}
