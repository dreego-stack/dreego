package transpiler

import (
	"strings"
	"testing"
)

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
