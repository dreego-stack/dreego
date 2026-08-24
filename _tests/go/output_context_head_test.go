package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestOutputContextHeadExpressionEscapes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>t := "<script>alert(1)</script>"</server>
<head><title>{{ t }}</title></head>
<body><h1>ok</h1></body>`,
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
		"www/routes/get.dreego": `<server>u := "javascript:alert(1)"</server>
<head><link rel="stylesheet" href="{{ u }}"></head>
<body><h1>ok</h1></body>`,
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
		"www/routes/get.dreego": `<server>u := "0;url=javascript:alert(1)"</server>
<head><meta http-equiv="refresh" content="{{ u }}"></head>
<body><h1>ok</h1></body>`,
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
			"www/routes/get.dreego": "<server>u := \"" + u + "\"</server>\n" +
				`<head><meta http-equiv = "refresh" content="{{ u }}"></head>` + "\n" +
				"<body><h1>ok</h1></body>",
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
		"www/routes/get.dreego": `<server>u := "0;url=javascript:alert(1)"</server>
<head><meta http-equiv = refresh content="{{ u }}"></head>
<body><h1>ok</h1></body>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("meta refresh with unquoted http-equiv must reject javascript: URL, got: %s", body)
	}
}

func TestOutputContextHeadLinkHrefUnquotedRejectsJavascript(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>u := "javascript:alert(1)"</server>
<head><link rel="stylesheet" href={{ u }}></head>
<body><h1>ok</h1></body>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "javascript:") {
		t.Fatalf("head link unquoted href must reject javascript: scheme, got: %s", body)
	}
}

func TestOutputContextHtmxOnJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>s := ` + "`\"><script>alert(1)</script>`" + `</server>
<body><button hx-on:click="{{ s }}">go</button></body>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<client>alert(1)</client>") {
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
			"www/routes/get.dreego": `<server>s := ` + "`\"><script>alert(1)</script>`" + `</server>
<body><div ` + attr + `="{{ s }}">x</div></body>`,
		})
		_, body := c.Get(t, "/")
		if strings.Contains(body, "<client>alert(1)</client>") {
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
		"www/routes/get.dreego": `<server>s := ` + "`\"><script>alert(1)</script>`" + `</server>
<body><button x-on:click="{{ s }}">go</button></body>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<client>alert(1)</client>") {
		t.Fatalf("x-on:click context must not emit raw script markup, got: %s", body)
	}
	if !strings.Contains(body, `\u003cscript\u003e`) {
		t.Fatalf("x-on:click context must JSON-encode the value, got: %s", body)
	}
}

func TestOutputContextShorthandClickJSONEncodes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>s := ` + "`\"><script>alert(1)</script>`" + `</server>
<body><button @click="{{ s }}">go</button></body>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "<client>alert(1)</client>") {
		t.Fatalf("@click context must not emit raw script markup, got: %s", body)
	}
	if !strings.Contains(body, `\u003cscript\u003e`) {
		t.Fatalf("@click context must JSON-encode the value, got: %s", body)
	}
}

func TestOutputContextXBindStyleNeutralizesBreakout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>s := "red; } </style><script>alert(1)</script>"</server>
<body><div x-bind:style="{{ s }}">x</div></body>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "</style>") {
		t.Fatalf("x-bind:style context must neutralize </style>, got: %s", body)
	}
	if strings.Contains(body, "<client>alert(1)</client>") {
		t.Fatalf("x-bind:style context must escape markup, got: %s", body)
	}
}

func TestOutputContextShorthandStyleNeutralizesBreakout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<server>s := "red; } </style><script>alert(1)</script>"</server>
<body><div :style="{{ s }}">x</div></body>`,
	})
	_, body := c.Get(t, "/")
	if strings.Contains(body, "</style>") {
		t.Fatalf(":style context must neutralize </style>, got: %s", body)
	}
	if strings.Contains(body, "<client>alert(1)</client>") {
		t.Fatalf(":style context must escape markup, got: %s", body)
	}
}
