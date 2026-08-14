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
		"expected root section, got <x>")
}

func TestParseExpectedSectionTag(t *testing.T) {
	parseExpectError(t,
		"<div><p>a</p></div>\nplain text after",
		"expected root section, got Text")
}

func TestParseRejectsRootHTML(t *testing.T) {
	parseExpectError(t,
		"<h1>Outside a section</h1>",
		"expected root section")
}

func TestParseRejectsRootComponentCall(t *testing.T) {
	parseExpectError(t,
		"<@Card />",
		"expected root section")
}

func TestParseRejectsRootText(t *testing.T) {
	parseExpectError(t,
		"Outside a section",
		"expected root section")
}

func TestParseRejectsDoctypeBeforeSection(t *testing.T) {
	parseExpectError(t,
		"<!doctype html>\n<div><p>x</p></div>",
		"expected root section")
}

// Control: a <go> section at the file start still parses as Go code.
func TestParseGoSectionFirstStillWorks(t *testing.T) {
	tokens, err := Lex("<go>msg := \"hi\"</go>\n<div><p>{{ msg }}</p></div>")
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Go) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Go))
	}
	if file.Go[0].Code != "msg := \"hi\"" {
		t.Fatalf("expected Go code %q, got %q", "msg := \"hi\"", file.Go[0].Code)
	}
}

// Item 9 (v0.0.25-feedback2): `<` inside Go strings. The lexer has no
// Go-string awareness — every <...> inside a <go> section lexes as a tag
// token (TokenTagOpen/TokenTagClose). parseGoSection currently only keeps
// TokenText, so <HASH> / <svg> content is silently dropped from the code.
// Chosen fix (option b): the go-section scanner reconstructs the raw content
// from the tokens (exactly like parseNonDivSection does for head/script/
// style), so the lexer stays as-is and <...> inside strings survives verbatim.

// "TO: <HASH>" in a quoted Go string must survive the parse.
func TestParseGoSectionKeepsLtInQuotedString(t *testing.T) {
	tokens, err := Lex(`<go>msg := "TO: <HASH>"</go>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Go) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Go))
	}
	want := `msg := "TO: <HASH>"`
	if file.Go[0].Code != want {
		t.Fatalf("expected Go code %q, got %q", want, file.Go[0].Code)
	}
}

// A backtick string containing <svg>...</svg> must survive the parse,
// including a self-closing <path .../> tag (regression: the trailing / was
// silently dropped, corrupting the reconstructed tag).
func TestParseGoSectionKeepsLtInBacktickString(t *testing.T) {
	src := "<go>svg := `<svg viewBox=\"0 0 24 24\"><path d=\"M12 2\"/></svg>`</go>"
	tokens, err := Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Go) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Go))
	}
	want := "svg := `<svg viewBox=\"0 0 24 24\"><path d=\"M12 2\"/></svg>`"
	if file.Go[0].Code != want {
		t.Fatalf("expected Go code %q, got %q", want, file.Go[0].Code)
	}
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
