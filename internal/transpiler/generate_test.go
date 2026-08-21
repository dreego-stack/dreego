package transpiler

import (
	"testing"
)

func TestBuildPattern(t *testing.T) {
	cases := map[string]string{
		"":                   "/{$}",
		"about":              "/about",
		"blog/[id]":          "/blog/{id}",
		"blog/_slug_":        "/blog/{slug}",
		"blog/(optional)":    "/blog",
		"blog/[id]/edit":     "/blog/{id}/edit",
		"blog/(group)/[id]":  "/blog/{id}",
		"blog/[...catchall]": "/blog/{catchall...}",
	}
	for in, want := range cases {
		if got := buildPattern(in); got != want {
			t.Errorf("buildPattern(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildPageName(t *testing.T) {
	cases := map[string]string{
		"":                "index",
		"about":           "about",
		"blog/[id]":       "blog_id",
		"blog/_slug_":     "blog_slug",
		"blog/(optional)": "blog_optional",
		"blog/[id]/edit":  "blog_id_edit",
	}
	for in, want := range cases {
		if got := buildPageName(in); got != want {
			t.Errorf("buildPageName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCleanSegment(t *testing.T) {
	cases := map[string]string{
		"[id]":       "id",
		"_slug_":     "slug",
		"(optional)": "optional",
		"about":      "about",
		"__x__":      "x",
	}
	for in, want := range cases {
		if got := cleanSegment(in); got != want {
			t.Errorf("cleanSegment(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPatternSegment(t *testing.T) {
	cases := map[string]string{
		"[id]":          "{id}",
		"_slug_":        "{slug}",
		"(optional)":    "",
		"about":         "about",
		"__x__":         "{x}",
		"[...catchall]": "{catchall...}",
	}
	for in, want := range cases {
		if got := patternSegment(in); got != want {
			t.Errorf("patternSegment(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDoubleBracketSegment(t *testing.T) {
	cases := map[string]string{
		"":                  "",
		"blog/[id]":         "",
		"blog/[[opt]]":      "[[opt]]",
		"blog/[[opt]]/get":  "[[opt]]",
		"blog/[id]/[[opt]]": "[[opt]]",
		"blog/[[a]]/[[b]]":  "[[a]]",
	}
	for in, want := range cases {
		if got := doubleBracketSegment(in); got != want {
			t.Errorf("doubleBracketSegment(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestErrorCatchPattern(t *testing.T) {
	cases := map[string]string{
		"/blog":      "/blog/{p...}",
		"/blog/{$}":  "/blog/{p...}",
		"/blog/{id}": "/blog/{id}/{p...}",
		"/{$}":       "/{p...}",
	}
	for in, want := range cases {
		if got := errorCatchPattern(in); got != want {
			t.Errorf("errorCatchPattern(%q) = %q, want %q", in, got, want)
		}
	}
}
