package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestOutputContextTextEscapesMarkup(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>v := "<script>alert(1)</script>"</go>
<div><p>{{ v }}</p></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("text context must escape script markup, got: %s", body)
	}
	if !strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatalf("text context must render escaped markup, got: %s", body)
	}
}

func TestOutputContextAttrEscapesQuotes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>v := ` + "`\" onmouseover=\"alert(1)`" + `</go>
<div><a title="{{ v }}">link</a></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, `" onmouseover="`) {
		t.Fatalf("attribute context must escape quotes, got: %s", body)
	}
	if !strings.Contains(body, "&#34; onmouseover=&#34;alert(1)") {
		t.Fatalf("attribute context must render escaped quotes, got: %s", body)
	}
}

func TestOutputContextURLRejectsJavascriptScheme(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>u := "javascript:alert(1)"</go>
<div><a href="{{ u }}">link</a></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("URL context must reject javascript: scheme, got: %s", body)
	}
	if !strings.Contains(body, `href="#`) {
		t.Fatalf("URL context must replace unsafe scheme with #, got: %s", body)
	}
}

func TestOutputContextURLAllowsHTTPS(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>u := "https://example.com/x"</go>
<div><a href="{{ u }}">link</a></div>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, `href="https://example.com/x"`) {
		t.Fatalf("URL context must keep https URLs, got: %s", body)
	}
}

func TestOutputContextURLRejectsDataScheme(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>u := "data:text/html,<script>alert(1)</script>"</go>
<div><img src="{{ u }}"></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "data:text/html") {
		t.Fatalf("URL context must reject data: scheme, got: %s", body)
	}
	if !strings.Contains(body, `src="#`) {
		t.Fatalf("URL context must replace unsafe scheme with #, got: %s", body)
	}
}

func TestOutputContextScriptAttrJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>s := ` + "`\"><script>alert(1)</script>`" + `</go>
<div><button onclick="{{ s }}">go</button></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("script context must not emit raw script markup, got: %s", body)
	}
	if !strings.Contains(body, `\u003cscript\u003e`) {
		t.Fatalf("script context must JSON-encode the value, got: %s", body)
	}
	if !strings.Contains(body, "&#34;") {
		t.Fatalf("script context must escape quotes, got: %s", body)
	}
}

func TestOutputContextStyleNeutralizesBreakout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>s := "red; } </style><script>alert(1)</script>"</go>
<div><div style="{{ s }}">x</div></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "</style>") {
		t.Fatalf("style context must neutralize </style>, got: %s", body)
	}
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("style context must escape markup, got: %s", body)
	}
}

func TestOutputContextRawOptInPassesThrough(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>h := "<b>trusted</b>"</go>
<div><p>{{ h|raw }}</p></div>`,
	})
	_, body := c.Get(t, "/")
	if !strings.Contains(body, "<b>trusted</b>") {
		t.Fatalf("raw opt-in must pass trusted HTML through, got: %s", body)
	}
}

func TestOutputContextComponentURLRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/components/Link.dreego": `Component Link (url string)
<div><a href="{{ url }}">go</a></div>`,
		"dreego/routes/get.dreego": `<go>u := "javascript:alert(1)"</go>
<div><@Link url={u}/></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("component URL context must reject javascript: scheme, got: %s", body)
	}
	if !strings.Contains(body, `href="#`) {
		t.Fatalf("component URL context must replace unsafe scheme with #, got: %s", body)
	}
}
