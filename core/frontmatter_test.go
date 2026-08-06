package core

import "testing"

// Architecture decision (Option C):
// Core must stay dependency-free (stdlib only). We do NOT pull in a full YAML
// library. ParseFrontmatter implements a minimal, YAML-like key:value frontmatter
// parser covering the metadata needs of docs/blogs (title, date, description,
// tags, list values). It intentionally does NOT implement the full YAML spec
// (no nested maps, no anchors, no multi-document YAML).
//
// API (defined by these tests, RED against current core):
//
//	ParseFrontmatter(src string) (frontmatter map[string]string, body string)
//
// Behavior:
//   - Leading "---\n" opens a frontmatter block; the next line "---" closes it.
//   - Only the block at the very top of src is treated as frontmatter.
//   - Each non-empty, non-comment line inside the block is "key: value".
//     The first ':' separates key from value; anything after is the value
//     (so "a: b" yields value "b", and quoted "a: b" yields "a: b").
//   - A list value "tags: [go, web]" yields the string "go, web".
//   - Lines before any "---" or a src with no "---" produce an empty map and
//     the whole src as body.

func TestParseFrontmatterExtractsKeys(t *testing.T) {
	src := "---\ntitle: X\n---\n<body>"
	fm, body := ParseFrontmatter(src)
	if len(fm) != 1 {
		t.Fatalf("expected 1 key, got %d: %#v", len(fm), fm)
	}
	if fm["title"] != "X" {
		t.Errorf("title = %q, want %q", fm["title"], "X")
	}
	if body != "<body>" {
		t.Errorf("body = %q, want %q", body, "<body>")
	}
}

func TestParseFrontmatterNoFrontmatter(t *testing.T) {
	src := "<body>"
	fm, body := ParseFrontmatter(src)
	if len(fm) != 0 {
		t.Errorf("expected empty map, got %#v", fm)
	}
	if body != src {
		t.Errorf("body = %q, want whole src %q", body, src)
	}
}

func TestParseFrontmatterMultiLineValues(t *testing.T) {
	src := "---\ntitle: My Page\ndate: 2026-08-06\ndescription: A page\n---\n<div>...</div>"
	fm, body := ParseFrontmatter(src)
	if len(fm) != 3 {
		t.Fatalf("expected 3 keys, got %d: %#v", len(fm), fm)
	}
	if fm["title"] != "My Page" {
		t.Errorf("title = %q, want %q", fm["title"], "My Page")
	}
	if fm["date"] != "2026-08-06" {
		t.Errorf("date = %q, want %q", fm["date"], "2026-08-06")
	}
	if fm["description"] != "A page" {
		t.Errorf("description = %q, want %q", fm["description"], "A page")
	}
	if body != "<div>...</div>" {
		t.Errorf("body = %q, want %q", body, "<div>...</div>")
	}
}

func TestParseFrontmatterListValue(t *testing.T) {
	src := "---\ntags: [go, web]\n---\n<body>"
	fm, _ := ParseFrontmatter(src)
	if len(fm) != 1 {
		t.Fatalf("expected 1 key, got %d: %#v", len(fm), fm)
	}
	// List values are normalized to a plain comma-joined string (no brackets).
	if fm["tags"] != "go, web" {
		t.Errorf("tags = %q, want %q", fm["tags"], "go, web")
	}
}

func TestParseFrontmatterColonInValue(t *testing.T) {
	src := "---\ndescription: \"a: b\"\n---\n<body>"
	fm, _ := ParseFrontmatter(src)
	if len(fm) != 1 {
		t.Fatalf("expected 1 key, got %d: %#v", len(fm), fm)
	}
	// A ':' inside the value must be preserved; only the first ':' splits.
	if fm["description"] != "a: b" {
		t.Errorf("description = %q, want %q", fm["description"], "a: b")
	}
}

func TestParseFrontmatterEmptyBlock(t *testing.T) {
	src := "---\n---\n<body>"
	fm, body := ParseFrontmatter(src)
	if len(fm) != 0 {
		t.Errorf("expected empty map, got %#v", fm)
	}
	if body != "<body>" {
		t.Errorf("body = %q, want %q", body, "<body>")
	}
}

func TestParseFrontmatterIntegratesWithData(t *testing.T) {
	src := "---\ntitle: My Page\n---\n<body>"
	fm, _ := ParseFrontmatter(src)
	c := NewSSR(nil, nil)
	for k, v := range fm {
		c.Set(k, v)
	}
	if got := c.Data("title"); got != "My Page" {
		t.Errorf("c.Data(\"title\") = %#v, want %q", got, "My Page")
	}
}
