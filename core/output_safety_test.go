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
