package core

import (
	"strings"
	"testing"
)

// parseExpectError lexes src, parses it, and asserts the returned error
// contains the given substring.
func parseExpectError(t *testing.T, src, want string) {
	t.Helper()
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("error %q does not contain %q", err.Error(), want)
	}
}

func TestParseDuplicateDiv(t *testing.T) {
	parseExpectError(t,
		"<div><p>a</p></div><div><p>b</p></div>",
		"duplicate <div> section")
}

func TestParseDuplicateHead(t *testing.T) {
	parseExpectError(t,
		"<head><title>a</title></head><head><title>b</title></head>",
		"duplicate <head> section")
}

func TestParseUnknownSection(t *testing.T) {
	parseExpectError(t,
		"<div><p>a</p></div><x>foo</x>",
		"unknown section <x>")
}

func TestParseExpectedSectionTag(t *testing.T) {
	parseExpectError(t,
		"<div><p>a</p></div>\nplain text after",
		"expected section tag, got Text")
}

func TestParseGoAttrsVariants(t *testing.T) {
	cases := []struct {
		name  string
		attrs string
		want  string
	}{
		{"double quoted", `type="json"`, "json"},
		{"single quoted", `type='json'`, "json"},
		{"unquoted", `type=json`, "json"},
		{"no type", `method="POST"`, ""},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseGoAttrs(tc.attrs); got != tc.want {
				t.Errorf("parseGoAttrs(%q) = %q, want %q", tc.attrs, got, tc.want)
			}
		})
	}
}
