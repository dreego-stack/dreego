package core

import (
	"strings"
	"testing"
)

func TestSafeTextEscapesMarkup(t *testing.T) {
	if got := SafeText(`<script>alert(1)</script>`); got != `&lt;script&gt;alert(1)&lt;/script&gt;` {
		t.Errorf("SafeText = %q, want escaped markup", got)
	}
}

func TestSafeAttrEscapesQuotes(t *testing.T) {
	if got := SafeAttr(`" onmouseover="alert(1)`); !strings.Contains(got, "&#34;") {
		t.Errorf("SafeAttr must escape quotes, got %q", got)
	}
}

func TestSafeURLAllowsSafeSchemes(t *testing.T) {
	for _, in := range []string{
		"https://example.com/x",
		"http://example.com",
		"mailto:a@b.de",
		"tel:+4912345",
		"/relative/path",
		"#anchor",
		"//cdn.example.com/lib.js",
		"",
	} {
		if got := SafeURL(in); got == "#" {
			t.Errorf("SafeURL(%q) = %q, want value kept", in, got)
		}
	}
}

func TestSafeURLRejectsUnsafeSchemes(t *testing.T) {
	for _, in := range []string{
		"javascript:alert(1)",
		"JaVaScRiPt:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"vbscript:msgbox(1)",
		"file:///etc/passwd",
	} {
		if got := SafeURL(in); got != "#" {
			t.Errorf("SafeURL(%q) = %q, want #", in, got)
		}
	}
}

func TestSafeURLRejectsObfuscatedScheme(t *testing.T) {
	for _, in := range []string{
		"java\nscript:alert(1)",
		"java\tscript:alert(1)",
		"java\rscript:alert(1)",
		"java script:alert(1)",
		"java\x0bscript:alert(1)",
		"java\x0cscript:alert(1)",
		"java\x00script:alert(1)",
	} {
		if got := SafeURL(in); got != "#" {
			t.Errorf("SafeURL(%q) = %q, want #", in, got)
		}
	}
}

func TestSafeURLRejectsObfuscatedSchemeInSrcset(t *testing.T) {
	if got := SafeURL("a.jpg 1x, javascript:alert(1)"); got != "#" {
		t.Errorf("SafeURL with mixed srcset value = %q, want #", got)
	}
}

func TestSafeURLAllowsValidSrcset(t *testing.T) {
	for _, in := range []string{
		"a.jpg 1x, https://cdn.example.com/b.jpg 2x",
		"a.jpg 480w, b.jpg 800w",
	} {
		if got := SafeURL(in); got == "#" {
			t.Errorf("SafeURL(%q) = %q, want value kept", in, got)
		}
	}
}

func TestSafeScriptJSONEncodes(t *testing.T) {
	got := SafeScript(`"><script>alert(1)</script>`)
	if !strings.Contains(got, `\u003cscript\u003e`) {
		t.Errorf("SafeScript must JSON-encode markup, got %q", got)
	}
	if !strings.Contains(got, "&#34;") {
		t.Errorf("SafeScript must escape quotes, got %q", got)
	}
	if strings.Contains(got, "<script>") {
		t.Errorf("SafeScript must not contain raw script markup, got %q", got)
	}
}

func TestSafeStyleNeutralizesBreakout(t *testing.T) {
	got := SafeStyle(`red; } </style><script>alert(1)</script>`)
	if strings.Contains(got, "</style>") {
		t.Errorf("SafeStyle must neutralize </style>, got %q", got)
	}
	if strings.Contains(got, "<script>") {
		t.Errorf("SafeStyle must escape markup, got %q", got)
	}
}

func TestSafeRawPassesThrough(t *testing.T) {
	if got := SafeRaw(`<b>trusted</b>`); got != `<b>trusted</b>` {
		t.Errorf("SafeRaw = %q, want passthrough", got)
	}
}

func TestSafeRefreshRejectsUnsafeURL(t *testing.T) {
	if got := SafeRefresh(`0;url=javascript:alert(1)`); strings.Contains(got, "javascript:") {
		t.Errorf("SafeRefresh must reject javascript: URL, got %q", got)
	}
	if got := SafeRefresh(`5`); got != "5" {
		t.Errorf("SafeRefresh without url= must pass through, got %q", got)
	}
	if got := SafeRefresh(`0;url=https://example.com`); !strings.Contains(got, "https://example.com") {
		t.Errorf("SafeRefresh must keep https URLs, got %q", got)
	}
}

func TestSafeRefreshRejectsWhitespaceAndCaseBypass(t *testing.T) {
	cases := []struct{ in, want string }{
		{"0; url = javascript:alert(1)", "0; url = #"},
		{"0;url =javascript:alert(1)", "0;url =#"},
		{"0;URL = javascript:alert(1)", "0;URL = #"},
		{"0;Url\t=\tjavascript:alert(1)", "0;Url\t=\t#"},
		{"0; url=javascript:alert(1)", "0; url=#"},
		{"0;url =\n javascript:alert(1)", "0;url =\n #"},
		{"0;url=java\nscript:alert(1)", "0;url=#"},
		{"0;url=java\rscript:alert(1)", "0;url=#"},
		{"0;url=java\tscript:alert(1)", "0;url=#"},
	}
	for _, c := range cases {
		got := SafeRefresh(c.in)
		if got != c.want {
			t.Errorf("SafeRefresh(%q) = %q, want sanitized %q", c.in, got, c.want)
		}
		if strings.Contains(got, "javascript:") {
			t.Errorf("SafeRefresh(%q) = %q, must reject javascript: value", c.in, got)
		}
	}
}

func TestSafeRefreshKeepsWhitespaceVariantWithSafeURL(t *testing.T) {
	for _, in := range []string{
		"0; url = https://example.com",
		"0;URL=https://example.com",
	} {
		got := SafeRefresh(in)
		if !strings.Contains(got, "https://example.com") {
			t.Errorf("SafeRefresh(%q) = %q, want https URL kept", in, got)
		}
		if strings.Contains(got, "javascript:") {
			t.Errorf("SafeRefresh(%q) = %q, must not contain javascript:", in, got)
		}
	}
}

func TestSafeStyleNeutralizesBreakoutCaseInsensitive(t *testing.T) {
	for _, in := range []string{
		"red; } </style><script>alert(1)</script>",
		"red; } </STYLE><script>alert(1)</script>",
		"red; } </StYlE><SCRIPT>alert(1)</SCRIPT>",
	} {
		got := SafeStyle(in)
		if strings.Contains(got, "</style>") || strings.Contains(got, "</STYLE>") || strings.Contains(got, "</StYlE>") {
			t.Errorf("SafeStyle(%q) = %q, must escape closing style tag", in, got)
		}
		if strings.Contains(got, "<script>") || strings.Contains(got, "<SCRIPT>") {
			t.Errorf("SafeStyle(%q) = %q, must escape markup", in, got)
		}
	}
	if got := SafeStyle(`<!-- --> red`); strings.Contains(got, "<!--") {
		t.Errorf("SafeStyle must escape comment open, got %q", got)
	}
}

func TestAttrNameAt(t *testing.T) {
	tag := `<a href="/x" onclick="go()" title="hi">`
	cases := []struct {
		i    int
		want string
	}{
		{strings.Index(tag, "/x"), "href"},
		{strings.Index(tag, "go()"), "onclick"},
		{strings.Index(tag, "hi"), "title"},
	}
	for _, c := range cases {
		if got := attrNameAt(tag, c.i); got != c.want {
			t.Errorf("attrNameAt(%q, %d) = %q, want %q", tag, c.i, got, c.want)
		}
	}
	if got := attrNameAt(tag, 0); got != "" {
		t.Errorf("attrNameAt at tag start = %q, want empty", got)
	}
	crTag := "<a\r\nhref=\"/x\">"
	if got := attrNameAt(crTag, strings.Index(crTag, "/x")); got != "href" {
		t.Errorf("attrNameAt with CR whitespace = %q, want href", got)
	}
	wsTag := `<a href = "/x">`
	if got := attrNameAt(wsTag, strings.Index(wsTag, "/x")); got != "href" {
		t.Errorf("attrNameAt with whitespace around = = %q, want href", got)
	}
	unqTag := `<a href=/x>`
	if got := attrNameAt(unqTag, strings.Index(unqTag, "/x")); got != "href" {
		t.Errorf("attrNameAt with unquoted value = %q, want href", got)
	}
}

func TestAttrSafeFuncClassifies(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<a href="{{ u }}">`, "SafeURL"},
		{`<img src="{{ u }}">`, "SafeURL"},
		{`<form action="{{ u }}">`, "SafeURL"},
		{`<button onclick="{{ s }}">`, "SafeScript"},
		{`<div style="{{ s }}">`, "SafeStyle"},
		{`<div title="{{ s }}">`, "SafeAttr"},
		{`<p>{{ s }}</p>`, "SafeAttr"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		tagStart := strings.LastIndex(c.content[:i], "<")
		if got := attrSafeFunc(c.content, tagStart, i); got != c.want {
			t.Errorf("attrSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestAttrSafeFuncClassifiesDirectives(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<button x-on:click="{{ s }}">`, "SafeScript"},
		{`<button @click="{{ s }}">`, "SafeScript"},
		{`<div x-on:mouseover="{{ s }}">`, "SafeScript"},
		{`<div x-bind:style="{{ s }}">`, "SafeStyle"},
		{`<div :style="{{ s }}">`, "SafeStyle"},
		{`<svg><use xlink:href="{{ u }}"></use></svg>`, "SafeURL"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		tagStart := strings.LastIndex(c.content[:i], "<")
		if got := attrSafeFunc(c.content, tagStart, i); got != c.want {
			t.Errorf("attrSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestAttrSafeFuncClassifiesHtmxAlpineScriptContexts(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<button hx-on:click="{{ s }}">`, "SafeScript"},
		{`<button hx-on::before-request="{{ s }}">`, "SafeScript"},
		{`<div x-data="{{ s }}">`, "SafeScript"},
		{`<div x-init="{{ s }}">`, "SafeScript"},
		{`<div x-effect="{{ s }}">`, "SafeScript"},
		{`<div x-html="{{ s }}">`, "SafeScript"},
		{`<div x-show="{{ s }}">`, "SafeScript"},
		{`<div x-model="{{ s }}">`, "SafeScript"},
		{`<div x-text="{{ s }}">`, "SafeScript"},
		{`<div x-transition="{{ s }}">`, "SafeScript"},
		{`<div once="{{ s }}">`, "SafeScript"},
		{`<div only="{{ s }}">`, "SafeScript"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		tagStart := strings.LastIndex(c.content[:i], "<")
		if got := attrSafeFunc(c.content, tagStart, i); got != c.want {
			t.Errorf("attrSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestIsScriptAttrClassification(t *testing.T) {
	for _, name := range []string{
		"onclick", "onload", "onerror", "onmouseover",
		"x-on:click", "x-on.mouseover", "@click",
		"hx-on:click", "hx-on::before-request",
		"x-data", "x-init", "x-effect", "x-html",
		"x-show", "x-model", "x-text", "x-transition",
	} {
		if !isScriptAttr(name) {
			t.Errorf("isScriptAttr(%q) = false, want true", name)
		}
	}
	for _, name := range []string{
		"on", "on1", "on-foo", "on:foo", "on.foo",
		"title", "href", "style", "x-bind:style", ":style", "x-bind:href",
	} {
		if isScriptAttr(name) {
			t.Errorf("isScriptAttr(%q) = true, want false", name)
		}
	}
}

func TestHeadSafeFuncClassifies(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<title>{{ t }}</title>`, "SafeText"},
		{`<meta name="description" content="{{ t }}">`, "SafeAttr"},
		{`<link href="{{ u }}">`, "SafeURL"},
		{`<meta http-equiv="refresh" content="{{ u }}">`, "SafeRefresh"},
		{`<meta http-equiv=refresh content="{{ u }}">`, "SafeRefresh"},
		{`<script src="{{ u }}"></script>`, "SafeURL"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		if got := headSafeFunc(c.content, i); got != c.want {
			t.Errorf("headSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestHeadSafeFuncClassifiesWhitespaceAroundEquals(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<meta http-equiv = "refresh" content="{{ u }}">`, "SafeRefresh"},
		{`<meta http-equiv = refresh content="{{ u }}">`, "SafeRefresh"},
		{`<meta HTTP-EQUIV="refresh" content="{{ u }}">`, "SafeRefresh"},
		{`<meta http-equiv = "refresh" CONTENT = "{{ u }}">`, "SafeRefresh"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		if got := headSafeFunc(c.content, i); got != c.want {
			t.Errorf("headSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestHeadSafeFuncClassifiesUnquotedValue(t *testing.T) {
	cases := []struct {
		content string
		want    string
	}{
		{`<link href={{ u }}>`, "SafeURL"},
		{`<meta name=description content={{ t }}>`, "SafeAttr"},
		{`<meta http-equiv=refresh content={{ u }}>`, "SafeRefresh"},
	}
	for _, c := range cases {
		i := strings.Index(c.content, "{{")
		if i < 0 {
			t.Fatalf("no placeholder in %q", c.content)
		}
		if got := headSafeFunc(c.content, i); got != c.want {
			t.Errorf("headSafeFunc(%q) = %q, want %q", c.content, got, c.want)
		}
	}
}

func TestAttrValueWhitespaceTolerance(t *testing.T) {
	cases := []struct {
		tag, attr, want string
	}{
		{`<meta http-equiv="refresh">`, "http-equiv", "refresh"},
		{`<meta http-equiv = "refresh">`, "http-equiv", "refresh"},
		{`<meta HTTP-EQUIV = refresh>`, "http-equiv", "refresh"},
		{`<meta http-equiv=refresh>`, "http-equiv", "refresh"},
		{`<meta charset="utf-8" http-equiv = 'x'>`, "http-equiv", "x"},
		{`<meta name="description">`, "http-equiv", ""},
	}
	for _, c := range cases {
		if got := attrValue(c.tag, c.attr); got != c.want {
			t.Errorf("attrValue(%q, %q) = %q, want %q", c.tag, c.attr, got, c.want)
		}
	}
}

func TestGenHeadURLAttrUsesSafeURL(t *testing.T) {
	out, err := genHead(`<link href="{{ u }}">`, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dreego.SafeURL") {
		t.Errorf("genHead link href must use SafeURL, got:\n%s", out)
	}
}

func TestGenHeadMetaRefreshUsesSafeRefresh(t *testing.T) {
	out, err := genHead(`<meta http-equiv="refresh" content="{{ u }}">`, "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dreego.SafeRefresh") {
		t.Errorf("genHead meta refresh content must use SafeRefresh, got:\n%s", out)
	}
}
