package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestOutputContextHeadExpressionEscapes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>t := "<script>alert(1)</script>"</go>
<head><title>{{ t }}</title></head>
<div><h1>ok</h1></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("head expression must escape markup, got: %s", body)
	}
	if !strings.Contains(body, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Fatalf("head expression must render escaped markup, got: %s", body)
	}
}

func TestOutputContextHeadLinkHrefRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>u := "javascript:alert(1)"</go>
<head><link rel="stylesheet" href="{{ u }}"></head>
<div><h1>ok</h1></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("head link href must reject javascript: scheme, got: %s", body)
	}
	if !strings.Contains(body, `href="#`) {
		t.Fatalf("head link href must replace unsafe scheme with #, got: %s", body)
	}
}

func TestOutputContextHeadMetaRefreshRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>u := "0;url=javascript:alert(1)"</go>
<head><meta http-equiv="refresh" content="{{ u }}"></head>
<div><h1>ok</h1></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("head meta refresh must reject javascript: URL, got: %s", body)
	}
	if !strings.Contains(body, `content="0;url=#`) {
		t.Fatalf("head meta refresh must replace unsafe URL with #, got: %s", body)
	}
}

func TestOutputContextHeadMetaRefreshWhitespaceEqualsRejectsJavascript(t *testing.T) {
	t.Parallel()
	for _, u := range []string{
		"0; url = javascript:alert(1)",
		"0;url =javascript:alert(1)",
		"0;URL = javascript:alert(1)",
	} {
		c := dreegotest.Serve(t, map[string]string{
			"dreego/routes/get.dreego": "<go>u := \"" + u + "\"</go>\n" +
				`<head><meta http-equiv = "refresh" content="{{ u }}"></head>` + "\n" +
				"<div><h1>ok</h1></div>",
		})
		_, body := c.Get(t, "/")
		if strings.Contains(body, "javascript:") {
			t.Fatalf("meta refresh with whitespace must reject javascript: URL %q, got: %s", u, body)
		}
	}
}

func TestOutputContextHeadMetaRefreshUnquotedEquivRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>u := "0;url=javascript:alert(1)"</go>
<head><meta http-equiv = refresh content="{{ u }}"></head>
<div><h1>ok</h1></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("meta refresh with unquoted http-equiv must reject javascript: URL, got: %s", body)
	}
}

func TestOutputContextHeadLinkHrefUnquotedRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>u := "javascript:alert(1)"</go>
<head><link rel="stylesheet" href={{ u }}></head>
<div><h1>ok</h1></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("head link unquoted href must reject javascript: scheme, got: %s", body)
	}
}

func TestOutputContextHtmxOnJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>s := ` + "`\"><script>alert(1)</script>`" + `</go>
<div><button hx-on:click="{{ s }}">go</button></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("hx-on:click context must not emit raw script markup, got: %s", body)
	}
	if !strings.Contains(body, `\u003cscript\u003e`) {
		t.Fatalf("hx-on:click context must JSON-encode the value, got: %s", body)
	}
}

func TestOutputContextAlpineEvaluatorJSONEncodes(t *testing.T) {
	t.Parallel()
	for _, attr := range []string{"x-data", "x-init", "x-effect", "x-html"} {
		c := dreegotest.Serve(t, map[string]string{
			"dreego/routes/get.dreego": `<go>s := ` + "`\"><script>alert(1)</script>`" + `</go>
<div><div ` + attr + `="{{ s }}">x</div></div>`,
		})
		_, body := c.Get(t, "/")
		if strings.Contains(body, "<script>alert(1)</script>") {
			t.Fatalf("%s context must not emit raw script markup, got: %s", attr, body)
		}
		if !strings.Contains(body, `\u003cscript\u003e`) {
			t.Fatalf("%s context must JSON-encode the value, got: %s", attr, body)
		}
	}
}

func TestOutputContextXOnClickJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>s := ` + "`\"><script>alert(1)</script>`" + `</go>
<div><button x-on:click="{{ s }}">go</button></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("x-on:click context must not emit raw script markup, got: %s", body)
	}
	if !strings.Contains(body, `\u003cscript\u003e`) {
		t.Fatalf("x-on:click context must JSON-encode the value, got: %s", body)
	}
}

func TestOutputContextShorthandClickJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>s := ` + "`\"><script>alert(1)</script>`" + `</go>
<div><button @click="{{ s }}">go</button></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("@click context must not emit raw script markup, got: %s", body)
	}
	if !strings.Contains(body, `\u003cscript\u003e`) {
		t.Fatalf("@click context must JSON-encode the value, got: %s", body)
	}
}

func TestOutputContextXBindStyleNeutralizesBreakout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>s := "red; } </style><script>alert(1)</script>"</go>
<div><div x-bind:style="{{ s }}">x</div></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "</style>") {
		t.Fatalf("x-bind:style context must neutralize </style>, got: %s", body)
	}
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf("x-bind:style context must escape markup, got: %s", body)
	}
}

func TestOutputContextShorthandStyleNeutralizesBreakout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go>s := "red; } </style><script>alert(1)</script>"</go>
<div><div :style="{{ s }}">x</div></div>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "</style>") {
		t.Fatalf(":style context must neutralize </style>, got: %s", body)
	}
	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Fatalf(":style context must escape markup, got: %s", body)
	}
}
