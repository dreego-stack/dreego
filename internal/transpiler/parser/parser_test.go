package parser

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/lexer"
)

// parseExpectError lexes src, parses it, and asserts the returned error
// contains the given substring.
func parseExpectError(t *testing.T, src, want string) {
	t.Helper()
	tokens, err := lexer.Lex(src)
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

func TestParseDuplicateBody(t *testing.T) {
	parseExpectError(t,
		"<body><p>a</p></body><body><p>b</p></body>",
		"duplicate <body> section")
}

func TestParseDuplicateHead(t *testing.T) {
	parseExpectError(t,
		"<head><title>a</title></head><head><title>b</title></head>",
		"duplicate <head> section")
}

func TestParseDuplicateMethodTemplateReportsMethod(t *testing.T) {
	parseExpectError(t,
		`<body method="post">one</body><body method="POST">two</body>`,
		"duplicate <body> section for method POST")
}

func TestParseUnknownSection(t *testing.T) {
	parseExpectError(t,
		"<body><p>a</p></body><x>foo</x>",
		"expected root section, got <x>")
}

func TestParseExpectedSectionTag(t *testing.T) {
	parseExpectError(t,
		"<body><p>a</p></body>\nplain text after",
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
		"<!doctype html>\n<body><p>x</p></body>",
		"expected root section")
}

// Control: a <server> section at the file start still parses as Go code.
func TestParseServerSectionFirstStillWorks(t *testing.T) {
	tokens, err := lexer.Lex("<server>msg := \"hi\"</server>\n<body><p>{{ msg }}</p></body>")
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Server) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Server))
	}
	if file.Server[0].Code != "msg := \"hi\"" {
		t.Fatalf("expected Go code %q, got %q", "msg := \"hi\"", file.Server[0].Code)
	}
}

// Item 9 (v0.0.25-feedback2): `<` inside Go strings. The lexer has no
// Go-string awareness — every <...> inside a <server> section lexes as a tag
// token (TokenTagOpen/TokenTagClose). parseServerSection currently only keeps
// TokenText, so <HASH> / <svg> content is silently dropped from the code.
// Chosen fix (option b): the go-section scanner reconstructs the raw content
// from the tokens (exactly like parseNonDivSection does for head/script/
// style), so the lexer stays as-is and <...> inside strings survives verbatim.

// "TO: <HASH>" in a quoted Go string must survive the parse.
func TestParseServerSectionKeepsLtInQuotedString(t *testing.T) {
	tokens, err := lexer.Lex(`<server>msg := "TO: <HASH>"</server>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Server) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Server))
	}
	want := `msg := "TO: <HASH>"`
	if file.Server[0].Code != want {
		t.Fatalf("expected Go code %q, got %q", want, file.Server[0].Code)
	}
}

// A backtick string containing <svg>...</svg> must survive the parse,
// including a self-closing <path .../> tag (regression: the trailing / was
// silently dropped, corrupting the reconstructed tag).
func TestParseServerSectionKeepsLtInBacktickString(t *testing.T) {
	src := "<server>svg := `<svg viewBox=\"0 0 24 24\"><path d=\"M12 2\"/></svg>`</server>"
	tokens, err := lexer.Lex(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Server) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Server))
	}
	want := "svg := `<svg viewBox=\"0 0 24 24\"><path d=\"M12 2\"/></svg>`"
	if file.Server[0].Code != want {
		t.Fatalf("expected Go code %q, got %q", want, file.Server[0].Code)
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
			if got := parseServerAttrs(tc.attrs); got != tc.want {
				t.Errorf("parseServerAttrs(%q) = %q, want %q", tc.attrs, got, tc.want)
			}
		})
	}
}

func TestParseGoMethodVariants(t *testing.T) {
	cases := []struct {
		name  string
		attrs string
		want  string
	}{
		{"double quoted", `method="post"`, "POST"},
		{"single quoted", `method='PUT'`, "PUT"},
		{"unquoted", `method=delete`, "DELETE"},
		{"no method", `type="json"`, ""},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseServerMethod(tc.attrs); got != tc.want {
				t.Errorf("parseServerMethod(%q) = %q, want %q", tc.attrs, got, tc.want)
			}
		})
	}
}

func TestParseServerSectionMethodAttr(t *testing.T) {
	tokens, err := lexer.Lex(`<server method="post">x := 1</server><body><p>{{ x }}</p></body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(file.Server) != 1 {
		t.Fatalf("expected 1 Go section, got %d", len(file.Server))
	}
	if file.Server[0].Method != "POST" {
		t.Fatalf("expected method POST, got %q", file.Server[0].Method)
	}
	if !file.Server[0].MethodExplicit {
		t.Fatal("expected MethodExplicit to be true")
	}
}

func TestParseServerSectionDefaultMethod(t *testing.T) {
	tokens, err := lexer.Lex(`<server>x := 1</server><body><p>{{ x }}</p></body>`)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if file.Server[0].Method != "GET" {
		t.Fatalf("expected default method GET, got %q", file.Server[0].Method)
	}
	if file.Server[0].MethodExplicit {
		t.Fatal("expected MethodExplicit to be false for default method")
	}
}
