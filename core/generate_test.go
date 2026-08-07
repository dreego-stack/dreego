package core

import "testing"

func TestBuildPattern(t *testing.T) {
	cases := map[string]string{
		"./routes":                   "/{$}",
		"./routes/":                  "/{$}",
		"./routes/about":             "/about",
		"./routes/blog/[id]":         "/blog/{id}",
		"./routes/blog/_slug_":       "/blog/{slug}",
		"./routes/blog/(optional)":   "/blog",
		"./routes/blog/[id]/edit":    "/blog/{id}/edit",
		"./routes/blog/(group)/[id]": "/blog/{id}",
	}
	for in, want := range cases {
		if got := buildPattern(in); got != want {
			t.Errorf("buildPattern(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildPageName(t *testing.T) {
	cases := map[string]string{
		"./routes":                 "index",
		"./routes/about":           "about",
		"./routes/blog/[id]":       "blog_id",
		"./routes/blog/_slug_":     "blog_slug",
		"./routes/blog/(optional)": "blog_optional",
		"./routes/blog/[id]/edit":  "blog_id_edit",
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
		"[[id]]":     "id",
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
		"[id]":       "{id}",
		"_slug_":     "{slug}",
		"(optional)": "",
		"about":      "about",
		"[[id]]":     "{id}",
		"__x__":      "{x}",
	}
	for in, want := range cases {
		if got := patternSegment(in); got != want {
			t.Errorf("patternSegment(%q) = %q, want %q", in, got, want)
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
