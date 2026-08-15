package core

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden fixtures with the current generator output")

// parseFile lexes and parses a .dreego source body into a File.
func parseFile(t *testing.T, src string) *File {
	t.Helper()
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}

// scopeHashFor mirrors the hash computation in generate.go/scanComponents so
// the generated data-scope attribute is deterministic for a given source.
func scopeHashFor(src string) string {
	h := sha256.Sum256([]byte(src))
	return hex.EncodeToString(h[:])[:12]
}

func goldenPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("testdata", "golden", name+".golden")
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := goldenPath(t, name)
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s (run with -update to create): %v", path, err)
	}
	if string(want) != got {
		t.Errorf("golden mismatch for %s:\n--- want ---\n%s\n--- got ---\n%s",
			name, string(want), got)
	}
}

// Golden for a simple route: <head> + <div>, no layout. Covers standalone head
// emission, the scoped wrapper div and static/expression text nodes.
func TestGoldenSimpleRoute(t *testing.T) {
	src := "<head>\n    <title>Home</title>\n</head>\n\n<div>\n    <h1>Welcome</h1>\n    <p>Hello, {{ name }}!</p>\n</div>\n"
	file := parseFile(t, src)
	got, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "home", "index", "/", scopeHashFor(src))
	if err != nil {
		t.Fatalf("GenerateMethodHandler: %v", err)
	}
	assertGolden(t, "simple_route", got)
}

// Golden for a component with props and a scoped <style> block.
func TestGoldenComponentWithStyle(t *testing.T) {
	src := "Component Badge (title string, tone string)\n\n<div class=\"badge\">\n    <span>{{ title }}</span>\n    <em>{{ tone }}</em>\n</div>\n\n<style>\n.badge { font-weight: bold; }\n.badge em { color: #666; }\n</style>\n"
	_, _, body := ParseHeader(src)
	file := parseFile(t, body)
	file.Component = &ComponentDef{
		Name: "Badge",
		Props: []Prop{
			{Name: "title", Type: "string"},
			{Name: "tone", Type: "string"},
		},
	}
	got, err := GenerateComponent(NewGenerator(), file, scopeHashFor(src))
	if err != nil {
		t.Fatalf("GenerateComponent: %v", err)
	}
	assertGolden(t, "component_with_style", got)
}

// Golden for a route rendered inside a layout with {#slot} and {#head}.
func TestGoldenRouteWithLayout(t *testing.T) {
	routeSrc := "<head><title>About</title></head>\n<div><h1>About us</h1></div>\n"
	layoutSrc := "<head></head>\n<div><!doctype html><html><head>{#head}</head><body><main>{#slot}</main></body></html></div>\n"

	file := parseFile(t, routeSrc)
	layout := parseFile(t, layoutSrc)

	got, _, err := GenerateMethodHandler(NewGenerator(), file, layout, "about", "about", "/about", scopeHashFor(routeSrc))
	if err != nil {
		t.Fatalf("GenerateMethodHandler: %v", err)
	}
	assertGolden(t, "route_with_layout", got)
}

// Golden for the router registration code emitted by GenerateRouter.
func TestGoldenRouter(t *testing.T) {
	routes := []RouteInfo{
		{HandlerName: "HandleIndex", RoutePath: "/", Method: "GET"},
		{HandlerName: "HandleAboutPost", RoutePath: "/about", Method: "POST"},
	}
	got := GenerateRouter(routes)
	assertGolden(t, "router", got)
}
